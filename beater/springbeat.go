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

type Springbeat struct {
    done   chan struct{}
    config config.Config
    client publisher.Client

    period          time.Duration
    urls            []*url.URL

    metricsStats    bool
    healthStats     bool
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
    config := config.DefaultConfig
    if err := cfg.Unpack(&config); err != nil {
        return nil, fmt.Errorf("Error reading config file: %v", err)
    }

    bt := &Springbeat{
        done: make(chan struct{}),
        config: config,
    }

    //define default URL if none provided
    var urlConfig []string
    if config.URLs != nil {
        urlConfig = config.URLs
    } else {
        urlConfig = []string{"http://127.0.0.1"}
    }

    bt.urls = make([]*url.URL, len(urlConfig))
    for i := 0; i < len(urlConfig); i++ {
        u, err := url.Parse(urlConfig[i])
        if err != nil {
            logp.Err("Invalid Spring Boot URL: %v", err)
            return nil, err
        }
        bt.urls[i] = u
    }

    if config.Stats.Metrics != nil {
        bt.metricsStats = *config.Stats.Metrics
    } else {
        bt.metricsStats = true
    }

    if config.Stats.Health != nil {
        bt.healthStats = *config.Stats.Health
    } else {
        bt.healthStats = true
    }

    if !bt.metricsStats && !bt.metricsStats {
        return nil, errors.New("Invalid statistics configuration")
    }

    logp.Debug("springbeat", "Init springbeat")
    logp.Debug("springbeat", "Period %v\n", bt.period)
    logp.Debug("springbeat", "Watch %v", bt.urls)
    logp.Debug("springbeat", "Metrics statistics %t\n", bt.metricsStats)
    logp.Debug("springbeat", "Health statistics %t\n", bt.healthStats)

    return bt, nil
}

func (bt *Springbeat) Run(b *beat.Beat) error {
    logp.Info("springbeat is running! Hit CTRL-C to stop it.")

    for _, u := range bt.urls {
        go func(u *url.URL) {

            ticker := time.NewTicker(bt.config.Period)
            counter := 1
            for {
                select {
                case <-bt.done:
                    goto GotoFinish
                case <-ticker.C:
                }

                timerStart := time.Now()

                if bt.metricsStats {
                    logp.Debug("springbeat", "Metrics stats for url: %v", u)
                    metrics_stats, err := bt.GetMetricsStats(*u)

                    if err != nil {
                        logp.Err("Error reading Metrics stats: %v", err)
                    } else {
                        logp.Debug("springbeat", "Metrics stats detail: %+v", metrics_stats)

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
                    logp.Debug("springbeat", "Health stats for url: %v", u)
                    health_stats, err := bt.GetHealthStats(*u)

                    if err != nil {
                        logp.Err("Error reading Health stats: %v", err)
                    } else {
                        logp.Debug("springbeat", "Health stats detail: %+v", health_stats)

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

func (bt *Springbeat) Stop() {
    logp.Debug("springbeat", "Stop springbeat")
    bt.client.Close()
    close(bt.done)
}
