package lib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

var QueryEscape = url.QueryEscape

func PathEscape(path string) string {
	if path == "nil" {
		return ""
	}
	pathParts := strings.Split(path, "/")
	newParts := make([]string, len(pathParts))

	for i, part := range pathParts {
		newParts[i] = url.PathEscape(part)
	}

	return strings.Join(newParts, "/")
}

type Path struct {
	Path string
}

func BuildPath(resourcePath string, values interface{}) (string, error) {
	r := regexp.MustCompile(`\{(.*)\}`)
	matches := r.FindSubmatch([]byte(resourcePath))
	if len(matches) > 0 {
		j, err := json.Marshal(&values)
		if err != nil {
			return "", err
		}
		var inInterface map[string]interface{}
		err = json.Unmarshal(j, &inInterface)
		if err != nil {
			return "", err
		}
		value := inInterface[string(matches[1])]
		valueInt, OkInt := value.(float64)
		if OkInt && valueInt == 0 {
			return "", CreateError(values, string(matches[1]))
		}

		stringValue := fmt.Sprintf("%v", value)
		if value == nil {

			stringValue = ""
			t := reflect.TypeOf(values)
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)

				// Get the field tag value
				tag := field.Tag.Get("path")
				if tag == string(matches[1]) {
					stringValue = fmt.Sprintf("%v", reflect.ValueOf(values).FieldByName(field.Name))
				}
			}
		}
		if string(matches[1]) == "path" {
			stringValue = PathEscape(stringValue)
		}
		return strings.ReplaceAll(resourcePath, string(matches[0]), stringValue), nil
	}

	return resourcePath, nil
}
