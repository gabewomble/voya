package data

import "encoding/json"

func CreateMetadata(data map[string]any) ([]byte, error) {
	metadata, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return metadata, nil
}
