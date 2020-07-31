package files_sdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Call(method string, config Config, resource string, params url.Values) (*[]byte, *http.Response, error) {
	res, err := CallRaw(method, config, resource, params)
	defaultValue := make([]byte, 0)
	if err != nil {
		return &defaultValue, res, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &defaultValue, res, err
	}
	config.GetLogger().Printf("Response Body: %v", string(data))
	re := ResponseError{}
	err = re.UnmarshalJSON(data)
	if err != nil {
		return &data, res, err
	}
	return &data, res, err
}

func CallRaw(method string, config Config, path string, params url.Values) (*http.Response, error) {
	httpClient := config.GetHttpClient()
	var body []byte
	var urlWithParams string
	bodyParams := make(map[string]string)
	switch method {
	case "GET", "HEAD", "DELETE":
		urlWithParams = config.RootPath() + path + "?" + params.Encode()
	default:
		for key, value := range params {
			bodyParams[key] = value[0]
		}
		body, _ = json.Marshal(bodyParams)
		config.GetLogger().Printf("Body: %v", string(body))
		params = url.Values{}
		urlWithParams = config.RootPath() + path
	}
	req, err := http.NewRequest(method, urlWithParams, bytes.NewBuffer(body))
	if err != nil {
		return &http.Response{}, err
	}
	config.SetHeaders(req.Header)
	if method == "post" {
		req.Header.Add("Content-Type", "application/json")
	}
	config.GetLogger().Printf("Headers: %v", req.Header)
	return httpClient.Do(req)
}
