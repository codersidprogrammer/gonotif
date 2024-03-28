package transport

import (
	"sync"
	"time"

	"github.com/gojek/heimdall"
	"github.com/gojek/heimdall/v7/httpclient"
)

var HttpClient *httpclient.Client

type httpClient struct {
	*Transport
	name string
}

func NewHttpClientTransport(name string) TransportService {
	return &httpClient{
		name: name,
		Transport: &Transport{
			lock: &sync.Mutex{},
		},
	}
}

// DoConnect implements TransportService.
func (*httpClient) DoConnect() error {
	// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
	backoffInterval := 2 * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := 5 * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	timeout := 1000 * time.Millisecond

	if HttpClient == nil {
		HttpClient = httpclient.NewClient(
			httpclient.WithHTTPTimeout(timeout),
			httpclient.WithRetrier(retrier),
			httpclient.WithRetryCount(5),
		)
	}

	return nil
}

// Close implements TransportService.
func (*httpClient) Close() error {
	panic("unimplemented")
}
