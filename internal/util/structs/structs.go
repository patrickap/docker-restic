package structs

import "github.com/mitchellh/mapstructure"

func ToMap(s interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := mapstructure.Decode(s, &m)

	return m, err
}
