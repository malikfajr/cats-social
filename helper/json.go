package helper

import (
	"encoding/json"
	"net/http"
)

func ParsingBody(request *http.Request, result interface{}) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	return err
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}, statusCode int) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Add("Content-Type", "Application/json")
	writer.WriteHeader(statusCode)

	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)

	PanicIfError(err)
}

type WebResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
