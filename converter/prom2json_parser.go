package converter

import (
	"crypto/tls"
	"fmt"
	dto "github.com/prometheus/client_model/go"
	"net/http"
	"os"
	"time"
)

func ExtractMetrics(promUrl string, cert string, key string) ([]*Family, error) {
	var err error
	result := []*Family{}
	mfChan := make(chan *dto.MetricFamily, 1024)
	transport, err := makeTransport(cert, key, true)
	if err != nil {
		fmt.Println(os.Stderr, err)
	}
	go func() {
		err := FetchMetricFamilies(promUrl, mfChan, transport)
		if err != nil {
			fmt.Println(err)
		}
	}()

	for mf := range mfChan {
		result = append(result, NewFamily(mf))
	}
	return result, err
}

func makeTransport(
	certificate string, key string,
	skipServerCertCheck bool,
) (*http.Transport, error) {
	// Start with the DefaultTransport for sane defaults.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// Conservatively disable HTTP keep-alives as this program will only
	// ever need a single HTTP request.
	transport.DisableKeepAlives = true
	// Timeout early if the server doesn't even return the headers.
	transport.ResponseHeaderTimeout = time.Minute
	tlsConfig := &tls.Config{InsecureSkipVerify: skipServerCertCheck}
	if certificate != "" && key != "" {
		cert, err := tls.LoadX509KeyPair(certificate, key)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	transport.TLSClientConfig = tlsConfig
	return transport, nil
}
