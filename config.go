package files_sdk

import (
	"github.com/hashicorp/go-retryablehttp"
	"log"
	"net/http"
	"os"
	"reflect"
)

const (
	ProductionEndpoint = "https://app.files.com"
	UserAgent          = "Files.com Go SDK"
)

var APIKey string

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Config struct {
	APIKey            string `header:"X-FilesAPI-Key"`
	Endpoint          string
	HttpClient        HttpClient
	AdditionalHeaders map[string]string
	Logger            retryablehttp.Logger
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
	debugLevel := os.Getenv("FILES_SDK_DEBUG")
	log.New(os.Stderr, "", log.LstdFlags)
	if debugLevel == "" {
		s.Logger = NullLogger{}
	} else {
		s.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return s.Logger
}

func (s *Config) RootPath() string {
	if s.Endpoint == "" {
		s.Endpoint = ProductionEndpoint
	}
	return s.Endpoint + "/api/rest/v1"
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

func (s *Config) SetHeaders(headers http.Header) {
	headers.Set("User-Agent", UserAgent)
	headers.Set("X-FilesAPI-Key", s.GetAPIKey())
	for key, value := range s.AdditionalHeaders {
		headers.Set(key, value)
	}
}
