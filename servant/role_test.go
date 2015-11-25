package servant

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRoleImplementationJson(t *testing.T) {
	exampleRole := RoleImplementation{
		RoleUrl: "http://example.com/role",
		Handlers: []Handler{
			Handler{
				Name:   "test 1",
				Method: "GET",
				Path:   "/test1/$(arg 1)",
			},
			Handler{
				Name:   "test 2",
				Method: "POST",
				Path:   "test2/$(arg 1)",
				Data:   "$(arg 2)",
			},
		},
	}
	result, err := json.Marshal(exampleRole)

	if err != nil {
		t.Error("Error while marshaling role: ", err)
	}

	expectedResult := []byte(`{"role_url":"http://example.com/role","handlers":[{"name":"test 1","method":"GET","path":"/test1/$(arg 1)"},{"name":"test 2","method":"POST","path":"test2/$(arg 1)","data":"$(arg 2)"}]}`)
	if !bytes.Equal(result, expectedResult) {
		t.Errorf("Marshaled role was not equal to expected role. Marshaled role: %s", result)
	}

}
