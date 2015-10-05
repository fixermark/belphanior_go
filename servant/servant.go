package servant

import (
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
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
func (s *Servant) RegisterHandler(description Handler) error {
	var err error
	s.role.Handlers = append(s.role.Handlers, description)
	mh := matchHandler{
		name: description.Name,
		method: description.Method,
		handler: reflect.ValueOf(description.Handler),
	}
	mh.pathRegexp, err =
		substitutionToRegexp(description.Path)
	if err != nil {
		return err
	}
	if description.Data != "" {
		mh.dataRegexp, err =
			substitutionToRegexp(description.Data)
		if err != nil {
			return err
		}
	}

	s.handlers = append(s.handlers, mh)
	return nil
}

// Attempts to call the handler that matches to the given request
func(s *Servant) CallHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println("Path matches: ", pathMatches)
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
	fmt.Println("Argument count: ", len(inputValues))
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
func substitutionToRegexp(sub string) (r *regexp.Regexp, err error) {
	fmt.Println("substitutionToRegexp: ", sub)
	argMatcher := regexp.MustCompile("\\$\\(.*?\\)")
	regexpString := argMatcher.ReplaceAllString(sub, "(.*)")
	r, err = regexp.Compile(regexpString)
	fmt.Println("substituted regexp: ", r)
	return
}

// General notes on implementation:
// - When handlers are registered, dissect them into two regexes (for path and body) and an arg count (for path and body)
// - Match: method, path and body regex must match method, path and body of request
// - On match, bundle args into Value slice (path followed by body)
// - Call handler with args
// - Marshall handler return value (if any) to string and respond with it and 200
// - if no matches, 404.


