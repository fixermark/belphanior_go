package servant

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

type matchHandler struct {
	name string
	method string
	pathRegexp *regexp.Regexp
	dataRegexp *regexp.Regexp
	handler reflect.Value
}

// Handler for incoming messages. Servant can be accessed via CallHandler, which
// will either match a request to a relevant registered Handler or yield a 404.
type Servant struct {
	role RoleImplementation
	handlers []matchHandler
}

func (s *Servant) SetRoleUrl(url string) {
	s.role.RoleUrl = url
}

// Adds the handler to the set of registered handlers and serves it
func (s *Servant) RegisterHandler(description Handler) {
	s.role.Handlers = append(s.role.Handlers, description)
	mh := matchHandler{
		name: description.Name,
		method: description.Method,
		handler: reflect.ValueOf(description.Handler),
	}
	mh.pathRegexp = substitutionToRegexp(description.Path)
	if description.Data != "" {
		mh.dataRegexp = substitutionToRegexp(description.Data)
	}

	s.handlers = append(s.handlers, mh)
}

func (s *Servant) ReportRoles(w http.ResponseWriter, r *http.Request) {
	var protocol struct {
		Roles []RoleImplementation `json:"roles"`
	}
	protocol.Roles = append(protocol.Roles, s.role)

	encoded, err := json.Marshal(protocol)
	if err != nil {
		reportError(w, err)
		return
	}

	headers := w.Header()
	headers.Add("Content-Type", "application/json")

	w.Write(encoded)
}

// Attempts to call the handler that matches to the given request
func (s *Servant) CallHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var data string
	if (r.Body != nil) {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			reportError(w, err)
			return
		}
		data = string(bs)
	}
	method := r.Method

	for _, handler := range s.handlers {
		if method != handler.method {
			continue
		}
		var pathMatches []string
		var dataMatches []string

		pathMatches = handler.pathRegexp.FindStringSubmatch(
			path)
		if pathMatches == nil {
			continue
		}
		if handler.dataRegexp != nil {
			dataMatches = handler.dataRegexp.FindStringSubmatch(
				data)
			if dataMatches == nil {
				continue
			}
		}
		var pathArgs, dataArgs []string
		if pathMatches != nil {
			pathArgs = pathMatches[1:]
		}
		if dataMatches != nil {
			dataArgs = dataMatches[1:]
		}
		// If we have not continued, this handler matched.
		runHandler(w, pathArgs, dataArgs, &handler)
		return
	}
	http.Error(w, "No handler found", 404)
}

// Evaluates the specified handler with the given arguments
func runHandler(
	w http.ResponseWriter,
	pathMatches, dataMatches []string,
	handler *matchHandler) {
	inputValues := make(
		[]reflect.Value,
		0,
		len(pathMatches) + len(dataMatches))
	for _, value := range pathMatches {
		inputValues = append(inputValues, reflect.ValueOf(value))
	}
	for _, value := range dataMatches {
		inputValues = append(inputValues, reflect.ValueOf(value))
	}
	result := handler.handler.Call(inputValues)
	switch len(result) {
	case 0:
		io.WriteString(w,"")
	case 1:
		io.WriteString(w, result[0].String())
	default:
		if result[1].IsNil() {
			io.WriteString(w, result[0].String())
		} else {
			// convert result[1] to error type and call String method

			errorMethod := result[1].MethodByName("Error")
			errorString := errorMethod.Call(nil)[0].String()
			http.Error(w, errorString, 500)
		}
	}
}

// Reports an error to the response
func reportError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 500)
}

// Convert a Belphanior role substitution to a regexp that will match strings
// conforming to the role substitution.
func substitutionToRegexp(sub string) *regexp.Regexp {
	argMatcher := regexp.MustCompile("\\$\\(.*?\\)")
	regexpString := argMatcher.ReplaceAllString(sub, "(.*)")
	return regexp.MustCompile(regexpString)
}

// Runs the servant as an HTTP server using net/http.  Defines the 'port' flag,
// which allows specification of the port the server runs on.

func (s *Servant) Run() {
	portFlag := flag.Int(
		"port",
		8080,
		"The port the servant should listen on.")
	flag.Parse()

	// TODO(mtomczak): Role reporting.
	http.HandleFunc("/protocol", s.ReportRoles)
	http.HandleFunc("/", s.CallHandler)

	fmt.Printf("Servant is listening on port %d\n", *portFlag)
	http.ListenAndServe(":" + strconv.Itoa(*portFlag), nil)
}
