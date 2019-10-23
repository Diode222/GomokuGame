package json

import "encoding/json"

func JsonResponse(code int, data interface{}) string {
	response := make(map[string]interface{})
	response["code"] = code
	response["data"] = data

	js, err := json.Marshal(response)
	if err != nil {
		return err.Error()
	}
	return string(js)
}
