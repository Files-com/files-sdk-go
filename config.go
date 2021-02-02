package files_sdk

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

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

type Config struct {
	APIKey                   string `header:"X-FilesAPI-Key"`
	SessionId                string `header:"X-FilesAPI-Auth"`
	Endpoint                 string
	Subdomain                string
	HttpClient               HttpClient
	AdditionalHeaders        map[string]string
	Logger                   retryablehttp.Logger
	Debug                    *bool
	maxConcurrentConnections int
}

func (s *Config) GetHttpClient() HttpClient {
	if s.HttpClient == nil || reflect.ValueOf(s.HttpClient).IsNil() {
		retryClient := retryablehttp.NewClient()
		retryClient.Logger = s.GetLogger()
		retryClient.RetryMax = 3
		s.HttpClient = retryClient.StandardClient()
	}
	return s.HttpClient
}

type NullLogger struct{}

func (n NullLogger) Printf(_ string, _ ...interface{}) {
}

func (s *Config) GetLogger() retryablehttp.Logger {
	var debugLevel string
	if s.Debug == nil {
		debugLevel = os.Getenv("FILES_SDK_DEBUG")
	} else {
		if *s.Debug {
			debugLevel = "debug"
		}
	}

	log.New(os.Stderr, "", log.LstdFlags)
	if debugLevel == "" {
		s.Logger = NullLogger{}
	} else {
		s.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return s.Logger
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
	headers.Set("User-Agent", UserAgent)
	if s.GetAPIKey() != "" {
		headers.Set("X-FilesAPI-Key", s.GetAPIKey())
	} else if s.SessionId != "" {
		headers.Set("X-FilesAPI-Auth", s.SessionId)
	}

	for key, value := range s.AdditionalHeaders {
		headers.Set(key, value)
	}
}

func (s *Config) SetMaxConcurrentConnections(value int) {
	s.maxConcurrentConnections = value
}

func (s *Config) MaxConcurrentConnections() int {
	if s.maxConcurrentConnections == 0 {
		s.SetMaxConcurrentConnections(10)
	}
	return s.maxConcurrentConnections
}
