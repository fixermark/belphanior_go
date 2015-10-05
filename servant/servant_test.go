package servant

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerRegistrationAndCall(t *testing.T) {
	s := Servant{}
	handlerInput := "unmodified"

	err := s.RegisterHandler(Handler{
		Name: "test 1",
		Method: "GET",
		Path: "/test1/$(arg 1)",
		Handler: func(arg1 string) string {
			handlerInput = arg1
			return "ok"
		},
	})
	if err != nil {
		t.Error("Could not register handler: ", err)
	}

	req, err := http.NewRequest("GET", "http://example.com/test1/argumentReceived", nil)
	if err != nil {
		t.Error("Could not create request: ", err)
	}
	resp := httptest.NewRecorder()
	s.CallHandler(resp, req)
	if handlerInput != "argumentReceived" {
		t.Error("Expected arg1 to be argumentReceived, was", handlerInput)
	}
	if resp.Code != 200 {
		t.Error("Expected 200 response, was ", resp.Code)
	}
	body := resp.Body.String()
	if body != "ok" {
		t.Error("Expected 'ok' body, was ", body)
	}

}

func TestNoHandlerForMessage(t *testing.T) {
	s := Servant{}
	handlerInput := "unmodified"

	err := s.RegisterHandler(Handler{
		Name: "test 1",
		Method: "GET",
		Path: "/test1/$(arg 1)",
		Handler: func(arg1 string) string {
			handlerInput = arg1
			return "ok"
		},
	})

	if err != nil {
		t.Error("Could not register handler: ", err)
	}


	req, err := http.NewRequest("GET", "http://example.com/notatest", nil)
	if err != nil {
		t.Error("Could not create request: ", err)
	}
	resp := httptest.NewRecorder()
	s.CallHandler(resp, req)
	if handlerInput != "unmodified" {
		t.Error("Expected unmodified handler input, was ", handlerInput)
	}
	if resp.Code != 404 {
		t.Error("Expected 404 response, was ", resp.Code)
	}
}
