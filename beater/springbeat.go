package beater

import (
    "errors"
    "fmt"
    "net/url"
    "time"

    "github.com/elastic/beats/libbeat/beat"
    "github.com/elastic/beats/libbeat/common"
    "github.com/elastic/beats/libbeat/logp"
    "github.com/elastic/beats/libbeat/publisher"

    "github.com/consulthys/springbeat/config"
)

const selector = "springbeat"

type Springbeat struct {
    period          time.Duration
    urls            []*url.URL

    beatConfig      *config.Config

    done            chan struct{}
    client          publisher.Client

    metricsStats    bool
    healthStats     bool
}

// Creates beater
func New() *Springbeat {
    return &Springbeat{
        done: make(chan struct{}),
    }
}

/// *** Beater interface methods ***///

func (bt *Springbeat) Config(b *beat.Beat) error {

    // Load beater beatConfig
    err := b.RawConfig.Unpack(&bt.beatConfig)
    if err != nil {
        return fmt.Errorf("Error reading config file: %v", err)
    }

    //define default URL if none provided
    var urlConfig []string
    if bt.beatConfig.Springbeat.URLs != nil {
        urlConfig = bt.beatConfig.Springbeat.URLs
    } else {
        urlConfig = []string{"http://127.0.0.1"}
    }

    bt.urls = make([]*url.URL, len(urlConfig))
    for i := 0; i < len(urlConfig); i++ {
        u, err := url.Parse(urlConfig[i])
        if err != nil {
            logp.Err("Invalid Spring Boot URL: %v", err)
            return err
        }
        bt.urls[i] = u
    }

    if bt.beatConfig.Springbeat.Stats.Metrics != nil {
        bt.metricsStats = *bt.beatConfig.Springbeat.Stats.Metrics
    } else {
        bt.metricsStats = true
    }

    if bt.beatConfig.Springbeat.Stats.Health != nil {
        bt.healthStats = *bt.beatConfig.Springbeat.Stats.Health
    } else {
        bt.healthStats = true
    }

    if !bt.metricsStats && !bt.metricsStats {
        return errors.New("Invalid statistics configuration")
    }

    return nil
}

func (bt *Springbeat) Setup(b *beat.Beat) error {

    // Setting default period if not set
    if bt.beatConfig.Springbeat.Period == "" {
        bt.beatConfig.Springbeat.Period = "10s"
    }

    bt.client = b.Publisher.Connect()

    var err error
    bt.period, err = time.ParseDuration(bt.beatConfig.Springbeat.Period)
    if err != nil {
        return err
    }

    logp.Debug(selector, "Init springbeat")
    logp.Debug(selector, "Period %v\n", bt.period)
    logp.Debug(selector, "Watch %v", bt.urls)
    logp.Debug(selector, "Metrics statistics %t\n", bt.metricsStats)
    logp.Debug(selector, "Health statistics %t\n", bt.healthStats)

    return nil
}

func (bt *Springbeat) Run(b *beat.Beat) error {
    logp.Info("springbeat is running! Hit CTRL-C to stop it.")

    for _, u := range bt.urls {
        go func(u *url.URL) {

            ticker := time.NewTicker(bt.period)
            counter := 1
            for {
                select {
                case <-bt.done:
                    goto GotoFinish
                case <-ticker.C:
                }

                timerStart := time.Now()

                if bt.metricsStats {
                    logp.Debug(selector, "Metrics stats for url: %v", u)
                    metrics_stats, err := bt.GetMetricsStats(*u)

                    if err != nil {
                        logp.Err("Error reading Metrics stats: %v", err)
                    } else {
                        logp.Debug(selector, "Metrics stats detail: %+v", metrics_stats)

                        event := common.MapStr{
                            "@timestamp":   common.Time(time.Now()),
                            "type":         "metrics",
                            "counter":      counter,
                            "metrics":      metrics_stats,
                        }

                        bt.client.PublishEvent(event)
                        logp.Info("Spring Boot /metrics stats sent")
                        counter++
                    }
                }

                if bt.healthStats {
                    logp.Debug(selector, "Health stats for url: %v", u)
                    health_stats, err := bt.GetHealthStats(*u)

                    if err != nil {
                        logp.Err("Error reading Health stats: %v", err)
                    } else {
                        logp.Debug(selector, "Health stats detail: %+v", health_stats)

                        event := common.MapStr{
                            "@timestamp":   common.Time(time.Now()),
                            "type":         "health",
                            "counter":      counter,
                            "health":       health_stats,
                        }

                        bt.client.PublishEvent(event)
                        logp.Info("Spring Boot /health stats sent")
                        counter++
                    }
                }

                timerEnd := time.Now()
                duration := timerEnd.Sub(timerStart)
                if duration.Nanoseconds() > bt.period.Nanoseconds() {
                    logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
                }
            }

        GotoFinish:
        }(u)
    }

    <-bt.done
    return nil
}

func (bt *Springbeat) Cleanup(b *beat.Beat) error {
    return nil
}

func (bt *Springbeat) Stop() {
    logp.Debug(selector, "Stop springbeat")
    close(bt.done)
}
