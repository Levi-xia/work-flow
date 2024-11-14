package utils

import "encoding/json"

func CompactJSON(content string) (string, error) {
	if content == "" {
		return "", nil
	}
	compactContent := ""
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(content), &jsonObj); err != nil {
		return "", err
	}
	compactBytes, err := json.Marshal(jsonObj)
	if err != nil {
		return "", err
	}
	compactContent = string(compactBytes)
	return compactContent, nil
}
