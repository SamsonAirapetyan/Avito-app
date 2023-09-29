package utils

import (
	logger2 "SergeyProject/pkg/logger"
	"encoding/json"
	"io"
	"net/http"
)

func StructDecode(r *http.Request, req interface{}) error {
	logger := logger2.GetLogger()
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		logger.Error("Unable to decode the request data", "error", err)
		return err
	}

	return nil
}

func ToJSON(req interface{}, w io.Writer) error {
	i := json.NewEncoder(w)
	return i.Encode(req)
}
