package extraction

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type JsonExtractor struct {
}

var unmarshalJsonCapture = func(result gjson.Result) (interface{}, error) {
	if result.IsObject() {
		jObject := map[string]interface{}{}
		err := json.Unmarshal([]byte(result.Raw), &jObject)
		if err == nil {
			return jObject, err
		}
	}

	if result.IsArray() {
		jStrSlice := []string{}
		err := json.Unmarshal([]byte(result.Raw), &jStrSlice)
		if err == nil {
			return jStrSlice, err
		}

		jFloatSlice := []float64{}
		err = json.Unmarshal([]byte(result.Raw), &jFloatSlice)
		if err == nil {
			return jFloatSlice, err
		}

		jIntSlice := []int{}
		err = json.Unmarshal([]byte(result.Raw), &jIntSlice)
		if err == nil {
			return jIntSlice, err
		}

		jBoolSlice := []bool{}
		err = json.Unmarshal([]byte(result.Raw), &jBoolSlice)
		if err == nil {
			return jBoolSlice, err
		}

	}

	if result.IsBool() {
		jBool := false
		err := json.Unmarshal([]byte(result.Raw), &jBool)
		if err == nil {
			return jBool, err
		}
	}

	return nil, fmt.Errorf("json could not be unmarshaled")
}

func (je JsonExtractor) extractFromString(source string, jsonPath string) (interface{}, error) {
	result := gjson.Get(source, jsonPath)

	// path not found
	if result.Raw == "" && result.Type == gjson.Null {
		return "", fmt.Errorf("json path not found")
	}

	switch result.Type {
	case gjson.String:
		return result.String(), nil
	case gjson.Null:
		return nil, nil
	case gjson.False:
		return false, nil
	case gjson.Number:
		number := result.String()
		if strings.Contains(number, ".") { // float
			return result.Float(), nil
		}
		return result.Int(), nil
	case gjson.True:
		return true, nil
	case gjson.JSON:
		return unmarshalJsonCapture(result)
	default:
		return "", nil
	}
}

func (je JsonExtractor) extractFromByteSlice(source []byte, jsonPath string) (interface{}, error) {
	result := gjson.GetBytes(source, jsonPath)

	// path not found
	if result.Raw == "" && result.Type == gjson.Null {
		return "", fmt.Errorf("json path not found")
	}

	switch result.Type {
	case gjson.String:
		return result.String(), nil
	case gjson.Null:
		return nil, nil
	case gjson.False:
		return false, nil
	case gjson.Number:
		number := result.String()
		if strings.Contains(number, ".") { // float
			return result.Float(), nil
		}
		return result.Int(), nil
	case gjson.True:
		return true, nil
	case gjson.JSON:
		return unmarshalJsonCapture(result)
	default:
		return "", nil
	}
}
