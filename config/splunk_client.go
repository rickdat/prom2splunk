package config

import (
	"crypto/tls"
	"fmt"
	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk/v2"
	"net/http"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

var singleInstance *splunk.Client

func GetSplunkInstance() *splunk.Client {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			httpClient := &http.Client{Timeout: time.Second * 20, Transport: tr}
			singleInstance := splunk.NewClient(
				httpClient,
				"https://prd-p-k9fe9.splunkcloud.com:8088/services/collector/event",
				"a0b2a086-44d4-4a36-bd22-116a92952444",
				"prometheus",
				"",
				"main",
			)
			return singleInstance
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstance
}
