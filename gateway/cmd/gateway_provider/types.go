package main

import (
	"encoding/json"
	"fmt"
)

type jsonMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func setenv(k, v string) error {
	return sendToCompose(jsonMessage{
		Type:    "setenv",
		Message: fmt.Sprintf("%v=%v", k, v),
	})
}

func infoMessage(message string) error {
	return sendToCompose(jsonMessage{
		Type:    "info",
		Message: message,
	})
}

func errorMessage(message string, err error) error {
	return sendToCompose(jsonMessage{
		Type:    "error",
		Message: fmt.Sprintf("%s: %v", message, err),
	})
}

func sendToCompose(message jsonMessage) error {
	marshal, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = fmt.Println(string(marshal))
	return err
}
