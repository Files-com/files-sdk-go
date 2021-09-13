package files_sdk

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

//go:embed _VERSION
var VERSION string

const (
	ProductionEndpoint = "https://{SUBDOMAIN}.files.com"
	UserAgent          = "Files.com Go SDK"
	DefaultDomain      = "app"
	APIPath            = "/api/rest/v1"
)

var APIKey string
var GlobalConfig Config

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
	Get(string) (*http.Response, error)
}

type Logger interface {
	Printf(string, ...interface{})
}

type Config struct {
	APIKey            string `header:"X-FilesAPI-Key"`
	SessionId         string `header:"X-FilesAPI-Auth"`
	Endpoint          string
	Subdomain         string
	standardClient    HttpClient
	rawClient         *retryablehttp.Client
	AdditionalHeaders map[string]string
	logger            Logger
	Debug             *bool
}

func (s *Config) SetHttpClient(client *http.Client) {
	s.GetRawClient().HTTPClient = client
}

func (s *Config) GetHttpClient() HttpClient {
	if s.standardClient == nil {
		s.standardClient = s.GetRawClient().StandardClient()
	}
	return s.standardClient
}

func (s *Config) GetRawClient() *retryablehttp.Client {
	if s.rawClient == nil {
		s.rawClient = retryablehttp.NewClient()
		s.rawClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
		s.rawClient.Logger = s.Logger()
		s.rawClient.RetryMax = 3
	}

	return s.rawClient
}

type NullLogger struct{}

func (n NullLogger) Printf(_ string, _ ...interface{}) {
}

func (s *Config) InDebug() bool {
	return s.Debug != nil && *s.Debug || (os.Getenv("FILES_SDK_DEBUG") != "")
}

func (s *Config) Logger() retryablehttp.Logger {
	if s.InDebug() {
		s.SetLogger(log.New(os.Stderr, "", log.LstdFlags))
	} else {
		s.SetLogger(NullLogger{})
	}
	return s.logger
}

func (s *Config) SetLogger(l Logger) {
	s.logger = l
}

func (s *Config) RootPath() string {
	if s.Subdomain == "" {
		s.Subdomain = DefaultDomain
	}
	if s.Endpoint == "" {
		s.Endpoint = strings.Replace(ProductionEndpoint, "{SUBDOMAIN}", s.Subdomain, 1)
	}
	return s.Endpoint + APIPath
}

func (s *Config) GetAPIKey() string {
	if s.APIKey != "" {
		return s.APIKey
	}
	if APIKey != "" {
		return APIKey
	}
	if os.Getenv("FILES_API_KEY") != "" {
		return os.Getenv("FILES_API_KEY")
	}
	return ""
}

func (s *Config) SetHeaders(headers *http.Header) {
	headers.Set("User-Agent", fmt.Sprintf("%v %v", UserAgent, strings.TrimSpace(VERSION)))
	if s.GetAPIKey() != "" {
		headers.Set("X-FilesAPI-Key", s.GetAPIKey())
	} else if s.SessionId != "" {
		headers.Set("X-FilesAPI-Auth", s.SessionId)
	}

	for key, value := range s.AdditionalHeaders {
		headers.Set(key, value)
	}
}
