package beater

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const METRICS_STATS = "/metrics"
const HEALTH_STATS = "/health"

type HealthStats struct {
    Status string `json:"status"`
    DiskSpace struct {
        Status string `json:"status"`
        Total uint64 `json:"total"`
        Free uint64 `json:"free"`
        Threshold uint64 `json:"threshold"`
    } `json:"diskSpace"`
    DB struct {
        Status string `json:"status"`
        Database string `json:"database"`
        Hello uint64 `json:"hello"`
    } `json:"db"`
}

type MetricsStats struct {
    Mem struct {
        Total uint64 `json:"total"`
        Free uint64 `json:"free"`
    } `json:"mem"`
    Processors uint64 `json:"processors"`
    LoadAverage float64 `json:"load_average"`
    Uptime struct {
        Total uint64 `json:"total"`
        Instance uint64 `json:"instance"`
    } `json:"uptime"`
    Heap struct {
        Total uint64 `json:"total"`
        Committed uint64 `json:"committed"`
        Init uint64 `json:"init"`
        Used uint64 `json:"used"`
    } `json:"heap"`
    NonHeap struct {
        Total uint64 `json:"total"`
        Committed uint64 `json:"committed"`
        Init uint64 `json:"init"`
        Used uint64 `json:"used"`
    } `json:"non_heap"`
    Threads struct {
        Total uint64 `json:"total"`
        TotalStarted uint64 `json:"started"`
        Peak uint64 `json:"peak"`
        Daemon uint64 `json:"daemon"`
    } `json:"non_heap"`
    Classes struct {
        Total uint64 `json:"total"`
        Loaded uint64 `json:"loaded"`
        Unloaded uint64 `json:"unloaded"`
    } `json:"classes"`
    GC struct {
        Scavenge struct {
            Count uint64 `json:"count"`
            Time uint64 `json:"time"`
        } `json:"scavenge"`
        Marksweep struct {
            Count uint64 `json:"count"`
            Time uint64 `json:"time"`
        } `json:"marksweep"`
    } `json:"gc"`
    Http struct {
        SessionsMax int64 `json:"max_sessions"`
        SessionsActive uint64 `json:"active_sessions"`
    } `json:"http"`
    DataSource struct {
        PrimaryActive uint64 `json:"primary_active"`
        PrimaryUsage float64 `json:"primary_usage"`
    } `json:"data_source"`
    GaugeResponse struct {
        Actuator float64 `json:"actuator,omitempty"`
        Autoconfig float64 `json:"autoconfig,omitempty"`
        Beans float64 `json:"beans,omitempty"`
        Configprops float64 `json:"configprops,omitempty"`
        Dump float64 `json:"dump,omitempty"`
        Env float64 `json:"env,omitempty"`
        Health float64 `json:"health,omitempty"`
        Info float64 `json:"info,omitempty"`
        Root float64 `json:"root,omitempty"`
        Trace float64 `json:"trace,omitempty"`
        Unmapped float64 `json:"unmapped,omitempty"`
    } `json:"gauge_response"`
    Status struct {
        TWO00 struct {
            Actuator uint64 `json:"actuator,omitempty"`
            Autoconfig uint64 `json:"autoconfig,omitempty"`
            Beans uint64 `json:"beans,omitempty"`
            Configprops uint64 `json:"configprops,omitempty"`
            Dump uint64 `json:"dump,omitempty"`
            Env uint64 `json:"env,omitempty"`
            Health uint64 `json:"health,omitempty"`
            Info uint64 `json:"info,omitempty"`
            Root uint64 `json:"root,omitempty"`
            Trace uint64 `json:"trace,omitempty"`
        } `json:"200"`
    } `json:"status"`
}

type RawMetricsStats struct {
    Mem uint64 `json:"mem"`
    MemFree uint64 `json:"mem.free"`
    Processors uint64 `json:"processors"`
    InstanceUptime uint64 `json:"instance.uptime"`
    Uptime uint64 `json:"uptime"`
    SystemloadAverage float64 `json:"systemload.average"`
    HeapCommitted uint64 `json:"heap.committed"`
    HeapInit uint64 `json:"heap.init"`
    HeapUsed uint64 `json:"heap.used"`
    Heap uint64 `json:"heap"`
    NonheapCommitted uint64 `json:"nonheap.committed"`
    NonheapInit uint64 `json:"nonheap.init"`
    NonheapUsed uint64 `json:"nonheap.used"`
    Nonheap uint64 `json:"nonheap"`
    ThreadsPeak uint64 `json:"threads.peak"`
    ThreadsDaemon uint64 `json:"threads.daemon"`
    ThreadsTotalStarted uint64 `json:"threads.totalStarted"`
    Threads uint64 `json:"threads"`
    Classes uint64 `json:"classes"`
    ClassesLoaded uint64 `json:"classes.loaded"`
    ClassesUnloaded uint64 `json:"classes.unloaded"`
    GCPsScavengeCount uint64 `json:"gc.ps_scavenge.count"`
    GCPsScavengeTime uint64 `json:"gc.ps_scavenge.time"`
    GCPsMarksweepCount uint64 `json:"gc.ps_marksweep.count"`
    GCPsMarksweepTime uint64 `json:"gc.ps_marksweep.time"`
    HttpSessionsMax int64 `json:"httpsessions.max"`
    HttpSessionsActive uint64 `json:"httpsessions.active"`
    DateSourcePrimaryActive uint64 `json:"datasource.primary.active"`
    DateSourcePrimaryUsage float64 `json:"datasource.primary.usage"`
    GaugeResponseActuator float64 `json:"gauge.response.actuator"`
    GaugeResponseBeans float64 `json:"gauge.response.beans"`
    GaugeResponseTrace float64 `json:"gauge.response.trace"`
    GaugeResponseAutoconfig float64 `json:"gauge.response.autoconfig"`
    GaugeResponseDump float64 `json:"gauge.response.dump"`
    GaugeResponseHealth float64 `json:"gauge.response.health"`
    GaugeResponseRoot float64 `json:"gauge.response.root"`
    GaugeResponseUnmapped float64 `json:"gauge.response.unmapped"`
    GaugeResponseInfo float64 `json:"gauge.response.info"`
    GaugeResponseEnv float64 `json:"gauge.response.env"`
    GaugeResponseConfigprops float64 `json:"gauge.response.configprops"`
    CounterStatus200Actuator uint64 `json:"counter.status.200.actuator"`
    CounterStatus200Autoconfig uint64 `json:"counter.status.200.autoconfig"`
    CounterStatus200Beans uint64 `json:"counter.status.200.beans"`
    CounterStatus200Configprops uint64 `json:"counter.status.200.configprops"`
    CounterStatus200Dump uint64 `json:"counter.status.200.dump"`
    CounterStatus200Env uint64 `json:"counter.status.200.env"`
    CounterStatus200Health uint64 `json:"counter.status.200.health"`
    CounterStatus200Info uint64 `json:"counter.status.200.info"`
    CounterStatus200Root uint64 `json:"counter.status.200.root"`
    CounterStatus200Trace uint64 `json:"counter.status.200.trace"`
}

func (bt *Springbeat) GetHealthStats(u url.URL) (*HealthStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + HEALTH_STATS)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP%s", res.Status)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    stats := &HealthStats{}
    err = json.Unmarshal([]byte(body), &stats)
    if err != nil {
        return nil, err
    }
    return stats, nil
}

func (bt *Springbeat) GetMetricsStats(u url.URL) (*MetricsStats, error) {
    res, err := http.Get(strings.TrimSuffix(u.String(), "/") + METRICS_STATS)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP%s", res.Status)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    raw_stats := &RawMetricsStats{}
    err = json.Unmarshal([]byte(body), &raw_stats)
    if err != nil {
        return nil, err
    }

    // Transform into usable JSON format
    stats := &MetricsStats{}
    stats.Mem.Free = raw_stats.MemFree
    stats.Mem.Total = raw_stats.Mem
    stats.Processors = raw_stats.Processors
    stats.LoadAverage = raw_stats.SystemloadAverage
    stats.Uptime.Total = raw_stats.Uptime
    stats.Uptime.Instance = raw_stats.InstanceUptime
    stats.Heap.Total = raw_stats.Heap
    stats.Heap.Init = raw_stats.HeapInit
    stats.Heap.Committed = raw_stats.HeapCommitted
    stats.Heap.Used = raw_stats.HeapUsed
    stats.NonHeap.Total = raw_stats.Nonheap
    stats.NonHeap.Init = raw_stats.NonheapInit
    stats.NonHeap.Committed = raw_stats.NonheapCommitted
    stats.NonHeap.Used = raw_stats.NonheapUsed
    stats.Threads.Total = raw_stats.Threads
    stats.Threads.TotalStarted = raw_stats.ThreadsTotalStarted
    stats.Threads.Peak = raw_stats.ThreadsPeak
    stats.Threads.Daemon = raw_stats.ThreadsDaemon
    stats.Classes.Total = raw_stats.Classes
    stats.Classes.Loaded = raw_stats.ClassesLoaded
    stats.Classes.Unloaded = raw_stats.ClassesUnloaded
    stats.GC.Scavenge.Count = raw_stats.GCPsScavengeCount
    stats.GC.Scavenge.Time = raw_stats.GCPsScavengeTime
    stats.GC.Marksweep.Count = raw_stats.GCPsMarksweepCount
    stats.GC.Marksweep.Time = raw_stats.GCPsMarksweepTime
    stats.Http.SessionsActive = raw_stats.HttpSessionsActive
    stats.Http.SessionsMax = raw_stats.HttpSessionsMax
    stats.DataSource.PrimaryActive = raw_stats.DateSourcePrimaryActive
    stats.DataSource.PrimaryUsage = raw_stats.DateSourcePrimaryUsage
    stats.GaugeResponse.Actuator = raw_stats.GaugeResponseActuator
    stats.GaugeResponse.Autoconfig = raw_stats.GaugeResponseAutoconfig
    stats.GaugeResponse.Beans = raw_stats.GaugeResponseBeans
    stats.GaugeResponse.Configprops = raw_stats.GaugeResponseConfigprops
    stats.GaugeResponse.Dump = raw_stats.GaugeResponseDump
    stats.GaugeResponse.Env = raw_stats.GaugeResponseEnv
    stats.GaugeResponse.Health = raw_stats.GaugeResponseHealth
    stats.GaugeResponse.Info = raw_stats.GaugeResponseInfo
    stats.GaugeResponse.Root = raw_stats.GaugeResponseRoot
    stats.GaugeResponse.Trace = raw_stats.GaugeResponseTrace
    stats.GaugeResponse.Unmapped = raw_stats.GaugeResponseUnmapped
    stats.Status.TWO00.Actuator = raw_stats.CounterStatus200Actuator
    stats.Status.TWO00.Autoconfig = raw_stats.CounterStatus200Autoconfig
    stats.Status.TWO00.Beans = raw_stats.CounterStatus200Beans
    stats.Status.TWO00.Configprops = raw_stats.CounterStatus200Configprops
    stats.Status.TWO00.Dump = raw_stats.CounterStatus200Dump
    stats.Status.TWO00.Env = raw_stats.CounterStatus200Env
    stats.Status.TWO00.Health = raw_stats.CounterStatus200Health
    stats.Status.TWO00.Info = raw_stats.CounterStatus200Info
    stats.Status.TWO00.Root = raw_stats.CounterStatus200Root
    stats.Status.TWO00.Trace = raw_stats.CounterStatus200Trace

    return stats, nil
}
