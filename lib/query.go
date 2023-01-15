package lib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

func PathEscape(path string) (string, error) {
	if path == "nil" {
		return "", nil
	}
	pathParts := strings.Split(path, "/")
	newParts := make([]string, len(pathParts))

	for i, part := range pathParts {
		newParts[i] = url.PathEscape(part)
	}
	var err error
	if len(newParts) > 1 {
		path, err = url.JoinPath(newParts[0], newParts[1:]...)
		if err != nil {
			return path, err
		}
	} else {
		path = newParts[0]
	}

	return Path{Path: path}.PruneStartingSlash().String(), nil
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
			stringValue, err = PathEscape(stringValue)
			if err != nil {
				return "", err
			}
		}
		return strings.ReplaceAll(resourcePath, string(matches[0]), stringValue), nil
	}

	return resourcePath, nil
}
