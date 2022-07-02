package httpclient

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
	"github.com/mcuadros/go-defaults"
)

//go:generate mockgen -source=client.go -destination=../../example/generated_mock.go -package=example HTTPClient
type HTTPClient interface {
	Get(uri string, headers map[string]string) (*http.Response, error)
	Post(uri string, body io.Reader, headers map[string]string) (*http.Response, error)
	PostForm(uri string, formVal map[string]string, headers map[string]string) (*http.Response, error)
	Put(uri string, body io.Reader, headers map[string]string) (*http.Response, error)
	Delete(uri string, headers map[string]string) (*http.Response, error)
	AddPlugin(plg heimdall.Plugin)
}

type httpClient struct {
	client *hystrix.Client
}

type customHttpClient struct {
	client http.Client
}

func (c *customHttpClient) Do(request *http.Request) (*http.Response, error) {
	return c.client.Do(request)
}

type HttpClientConfig struct {
	Timeout                time.Duration `default:"3000ms"`
	MaxConcurrentRequests  int           `default:"100"`
	ErrorPercentThreshold  int           `default:"10"`
	SleepWindow            int           `default:"20"`
	RequestVolumeThreshold int           `default:"100"`
	RetryCount             int           `default:"3"`

	BackoffInitialTimeout time.Duration `default:"1000ms"`
	BackoffMaxTimeout     time.Duration `default:"5000ms"`
	BackoffJitterInterval time.Duration `default:"200ms"`
	BackoffExponentFactor float64       `default:"2"`

	ProxyUrl string
}

func NewClient(cfg HttpClientConfig, name string) *httpClient {
	var customHc = func(*hystrix.Client) {}

	backoff := heimdall.NewExponentialBackoff(
		cfg.BackoffInitialTimeout,
		cfg.BackoffMaxTimeout,
		cfg.BackoffExponentFactor,
		cfg.BackoffJitterInterval)

	retrier := heimdall.NewRetrier(backoff)

	// Use proxy request
	if cfg.ProxyUrl != "" {
		proxyUrl, err := url.Parse(cfg.ProxyUrl)
		if err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
			customHc = hystrix.WithHTTPClient(&customHttpClient{client: http.Client{Transport: transport}})
		}
	}

	client := hystrix.NewClient(
		customHc,
		hystrix.WithHTTPTimeout(cfg.Timeout),
		hystrix.WithCommandName(name),
		hystrix.WithHystrixTimeout(cfg.Timeout),
		hystrix.WithMaxConcurrentRequests(cfg.MaxConcurrentRequests),
		hystrix.WithErrorPercentThreshold(cfg.ErrorPercentThreshold),
		hystrix.WithSleepWindow(cfg.SleepWindow),
		hystrix.WithRequestVolumeThreshold(cfg.RequestVolumeThreshold),
		hystrix.WithRetryCount(cfg.RetryCount),
		hystrix.WithRetrier(retrier),
		hystrix.WithFallbackFunc(FallbackFunction),
	)

	return &httpClient{
		client: client,
	}
}

func (h *httpClient) Get(uri string, headers map[string]string) (*http.Response, error) {
	var (
		httpHeader http.Header = http.Header{}
	)

	for k, v := range headers {
		httpHeader.Add(k, v)
	}

	return h.client.Get(uri, httpHeader)
}

func (h *httpClient) Post(uri string, body io.Reader, headers map[string]string) (res *http.Response, err error) {
	var (
		httpHeader http.Header = http.Header{}
	)

	for k, v := range headers {
		httpHeader.Add(k, v)
	}

	return h.client.Post(uri, body, httpHeader)
}

func (h *httpClient) PostForm(uri string, formVal map[string]string, headers map[string]string) (*http.Response, error) {
	var (
		req  *http.Request
		err  error
		form = url.Values{}
	)

	for k, v := range formVal {
		form.Set(k, v)
	}

	req, err = http.NewRequest(http.MethodPost, uri, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.PostForm = form

	return h.client.Do(req)
}

func (h *httpClient) Put(uri string, body io.Reader, headers map[string]string) (*http.Response, error) {
	var (
		httpHeader http.Header = http.Header{}
	)

	for k, v := range headers {
		httpHeader.Add(k, v)
	}

	return h.client.Put(uri, body, httpHeader)
}

func (h *httpClient) Delete(uri string, headers map[string]string) (*http.Response, error) {
	var (
		httpHeader http.Header = http.Header{}
	)

	for k, v := range headers {
		httpHeader.Add(k, v)
	}

	return h.client.Delete(uri, httpHeader)
}

func (h *httpClient) AddPlugin(plg heimdall.Plugin) {
	h.client.AddPlugin(plg)
}

//FallbackFunction for hystrix fallback func
func FallbackFunction(err error) error {
	return err
}

func DefaultConfig() HttpClientConfig {
	cfg := HttpClientConfig{}
	defaults.SetDefaults(&cfg)
	return cfg
}
