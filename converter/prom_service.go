package converter

import (
	"fmt"
	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk/v2"
	"os"
	"promToSplunk/config"
	"sync"
	"time"
)

func UploadMetricJson() (string, error) {
	resp, err := ExtractMetrics("http://localhost:9090/metrics", "", "")
	if err != nil {
		return "", err
	}
	var wg sync.WaitGroup
	wg.Add(len(resp))
	var client = config.GetSplunkInstance()
	var host, _ = os.Hostname()
	var event = splunk.Event{
		Time:       splunk.EventTime{time.Now()},
		Host:       host,
		Source:     "prometheus",
		SourceType: "prometheus",
		Index:      "main",
		Event:      nil,
	}
	for _, s := range resp {
		go func(s *Family) {
			event.Event = s
			err = client.LogEvent(&event)
			if err != nil {
				fmt.Println(err)
			}
		}(s)

	}
	wg.Wait()
	if err != nil {
		return "", err
	}
	return "success sent", err
}

type SplunkMetric struct {
}
