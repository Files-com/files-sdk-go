package lib

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func CheckRequired(iStruct interface{}, values *url.Values) error {
	var errors = make([]string, 0)
	val := reflect.ValueOf(iStruct)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if iStruct == nil {
		return nil
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("CheckRequired expects struct input. Got %v", val.Kind())
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous { // unexported
			continue
		}
		tag := sf.Tag.Get("required")
		urlTag := sf.Tag.Get("url")
		key := strings.Split(urlTag, ",")[0]

		if tag == "true" && values.Get(key) == "" {
			errors = append(
				errors,
				CreateError(iStruct, sf.Name).Error(),
			)
		}
	}
	if len(errors) != 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

func CreateError(i interface{}, name string) error {
	structName := reflect.TypeOf(i).Name()
	return fmt.Errorf("missing required field: %v{}.%v", structName, name)
}
