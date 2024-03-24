package utils

import "encoding/json"

func StructToMap(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Marshal struct to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON to map
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
