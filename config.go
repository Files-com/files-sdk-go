package files_sdk

import (
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	UserAgent          = "Files.com Go SDK"
)

var ProductionEndpoint string = "https://app.files.com"
var APIKey string

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
	Get(string) (*http.Response, error)
}

type Config struct {
	APIKey                   string `header:"X-FilesAPI-Key"`
	Endpoint                 string
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
		_, err := os.Stat("../log")
		if os.IsNotExist(err) {
			os.Mkdir("../log", 0700)
		}

		f, err := os.OpenFile("../log/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		s.Logger = log.New(f, "", log.LstdFlags)
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

func (s *Config) SetHeaders(headers *http.Header) {
	headers.Set("User-Agent", UserAgent)
	headers.Set("X-FilesAPI-Key", s.GetAPIKey())
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
