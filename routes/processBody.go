package routes

import (
	"encoding/json"
	"io"
)

func ProcessBody(closer io.ReadCloser, responseStruct interface{}) error {
	body, err := io.ReadAll(closer)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, responseStruct); err != nil {
		return err
	}

	return nil
}