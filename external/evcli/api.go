package evcli

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/exograd/eventline/pkg/eventline"
	"github.com/exograd/go-daemon/check"
)

type APIError struct {
	Message string          `json:"error"`
	Code    string          `json:"code,omitempty"`
	RawData json.RawMessage `json:"data,omitempty"`
	Data    interface{}     `json:"-"`
}

type InvalidRequestBodyError struct {
	ValidationErrors check.ValidationErrors `json:"validation_errors"`
}

func (err APIError) Error() string {
	return err.Message
}

func (err *APIError) UnmarshalJSON(data []byte) error {
	type APIError2 APIError

	err2 := APIError2(*err)
	if jsonErr := json.Unmarshal(data, &err2); jsonErr != nil {
		return jsonErr
	}

	switch err2.Code {
	case "invalid_request_body":
		var errData InvalidRequestBodyError

		if err2.RawData != nil {
			if err := json.Unmarshal(err2.RawData, &errData); err != nil {
				return fmt.Errorf("invalid jsv errors: %w", err)
			}

			err2.Data = &errData
		}
	}

	*err = APIError(err2)
	return nil
}

func IsInvalidRequestBodyError(err error) (bool, check.ValidationErrors) {
	var apiError *APIError

	if !errors.As(err, &apiError) {
		return false, nil
	}

	requestBodyErr, ok := apiError.Data.(*InvalidRequestBodyError)
	if !ok {
		return false, nil
	}

	return true, requestBodyErr.ValidationErrors
}

type ProjectPage struct {
	Elements eventline.Projects `json:"elements"`
	Previous *eventline.Cursor  `json:"previous,omitempty"`
	Next     *eventline.Cursor  `json:"next,omitempty"`
}

type Parameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
}

type Parameters []*Parameter

type JobPage struct {
	Elements eventline.Jobs    `json:"elements"`
	Previous *eventline.Cursor `json:"previous,omitempty"`
	Next     *eventline.Cursor `json:"next,omitempty"`
}
