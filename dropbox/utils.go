package dropbox

import (
	"encoding/json"
	"fmt"
)

func parseEntries(respBody []byte) ([]*Entry, error) {
	msg := make(map[string]interface{})
	if err := json.Unmarshal(respBody, &msg); err != nil {
		return nil, err
	}

	results := []*Entry{}

	if _, exist := msg["entries"]; !exist {
		return nil, fmt.Errorf("invalid response body %s", string(respBody))
	}
	entries, ok := msg["entries"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response body %s", string(respBody))
	}

	for _, entry := range entries {
		e, ok := entry.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid response body %s", string(respBody))
		}
		tag, ok := e[".tag"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid response body %s", string(respBody))
		}
		name, ok := e["name"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid response body %s", string(respBody))
		}
		id, ok := e["id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid response body %s", string(respBody))
		}

		ret := &Entry{
			Id:   id,
			Name: name,
			Tag:  tag,
		}

		results = append(results, ret)
	}

	return results, nil
}
