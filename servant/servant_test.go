package servant

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoleReport(t *testing.T) {
	s := Servant{}
	s.SetRoleUrl("https://example.com/foo")
	s.RegisterHandler(Handler{
		Name:    "test 1",
		Method:  "GET",
		Path:    "/test1/$(arg 1)",
		Handler: func() {},
	})
	s.RegisterHandler(Handler{
		Name:    "test 2",
		Method:  "POST",
		Path:    "test2/$(arg 1)",
		Data:    "$(arg 2)",
		Handler: func() {},
	})
	req, err := http.NewRequest("GET", "http://example.com/protocol", nil)
	if err != nil {
		t.Error("Could not create request: ", err)
	}
	resp := httptest.NewRecorder()
	s.ReportRoles(resp, req)
	if resp.Code != 200 {
		t.Error("Expected 200 response, was ", resp.Code)
	}
	body := resp.Body.String()
	const match = "{\"roles\":[{\"role_url\":\"https://example.com/foo\"," +
		"\"handlers\":[{\"name\":\"test 1\",\"method\":\"GET\"," +
		"\"path\":\"/test1/$(arg 1)\"},{\"name\":\"test 2\"," +
		"\"method\":\"POST\",\"path\":\"test2/$(arg 1)\"," +
		"\"data\":\"$(arg 2)\"}]}]}"
	if body != match {
		t.Errorf(
			"Protocol did not match, needed\n%s\nwas\n%s\n",
			match, body)
	}
}

func TestHandlerRegistrationAndCall(t *testing.T) {
	s := Servant{}
	handlerInput := "unmodified"

	s.RegisterHandler(Handler{
		Name:   "test 1",
		Method: "GET",
		Path:   "/test1/$(arg 1)",
		Handler: func(arg1 string) string {
			handlerInput = arg1
			return "ok"
		},
	})

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

	s.RegisterHandler(Handler{
		Name:   "test 1",
		Method: "GET",
		Path:   "/test1/$(arg 1)",
		Handler: func(arg1 string) string {
			handlerInput = arg1
			return "ok"
		},
	})

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

func TestMultipleHandlerRegistration(t *testing.T) {
	s := Servant{}
	handlerInput := "unmodified"

	s.RegisterHandler(Handler{
		Name:   "test 1",
		Method: "GET",
		Path:   "/test1/$(arg 1)",
		Handler: func(arg1 string) string {
			handlerInput = arg1
			return "ok"
		},
	})
	s.RegisterHandler(Handler{
		Name:   "test 2",
		Method: "GET",
		Path:   "/test2/$(arg 1)/$(arg 2)",
		Handler: func(arg1, arg2 string) string {
			handlerInput = arg1 + arg2
			return "ok"
		},
	})

	req, err := http.NewRequest("GET", "http://example.com/test2/cat/dog", nil)
	if err != nil {
		t.Error("Could not create request: ", err)
	}
	resp := httptest.NewRecorder()
	s.CallHandler(resp, req)
	if handlerInput != "catdog" {
		t.Error("Expected arg1 to be catdog, was", handlerInput)
	}
	if resp.Code != 200 {
		t.Error("Expected 200 response, was ", resp.Code)
	}
	body := resp.Body.String()
	if body != "ok" {
		t.Error("Expected 'ok' body, was ", body)
	}

}

func TestNoHandlerResult(t *testing.T) {
	s := Servant{}
	handlerInput := "unmodified"

	s.RegisterHandler(Handler{
		Name:   "test 1",
		Method: "GET",
		Path:   "/test1/$(arg 1)",
		Handler: func(arg1 string) {
			handlerInput = arg1
		},
	})

	req, err := http.NewRequest("GET", "http://example.com/test1/catdog", nil)
	if err != nil {
		t.Error("Could not create request: ", err)
	}
	resp := httptest.NewRecorder()
	s.CallHandler(resp, req)
	if handlerInput != "catdog" {
		t.Error("Expected arg1 to be catdog, was", handlerInput)
	}
	if resp.Code != 200 {
		t.Error("Expected 200 response, was ", resp.Code)
	}
	body := resp.Body.String()
	if body != "" {
		t.Error("Expected '' body, was ", body)
	}

}

func TestErrorHandlerResult(t *testing.T) {
	s := Servant{}

	s.RegisterHandler(Handler{
		Name:   "test 1",
		Method: "GET",
		Path:   "/test1/$(arg 1)",
		Handler: func(arg1 string) (result string, err error) {
			if arg1 == "bomb" {
				err = errors.New("kaboom")
			} else {
				result = arg1
			}
			return
		},
	})

	req, err := http.NewRequest("GET", "http://example.com/test1/catdog", nil)
	if err != nil {
		t.Error("Could not create request: ", err)
	}
	resp := httptest.NewRecorder()
	s.CallHandler(resp, req)
	if resp.Code != 200 {
		t.Error("Expected 200 response, was ", resp.Code)
	}
	body := resp.Body.String()
	if body != "catdog" {
		t.Error("Expected 'catdog' body, was ", body)
	}

	req, err = http.NewRequest("GET", "http://example.com/test1/bomb", nil)
	if err != nil {
		t.Error("Could not create request: ", err)
	}
	resp = httptest.NewRecorder()
	s.CallHandler(resp, req)
	if resp.Code != 500 {
		t.Error("Expected 500 response, was ", resp.Code)
	}
	body = resp.Body.String()
	if body != "kaboom\n" {
		t.Error("Expected 'kaboom' body, was", body)
	}

}
