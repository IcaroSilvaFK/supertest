// This package is a library for testing http requests
// Inspired by supertest on JavaScript ecosystem
package supertest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
)

type Tester struct {
	httpUrl      string `validate:"required,url"`
	httpHeaders  map[string]string
	httpBody     io.Reader
	httpMethod   string `validate:"required"`
	httpStatus   int
	response     interface{}
	resp         *http.Response
	validation   *validator.Validate
	validateBody bool
}

var errors = make(map[string]string, 0)

type TesterInterface interface {
	Method(string) TesterInterface
	Url(string) TesterInterface
	Json(interface{}) TesterInterface
	Headers(map[string]string) TesterInterface
	Query(map[string]string) TesterInterface
	Body([]byte) TesterInterface
	Status(int) TesterInterface
	ValidateBody() TesterInterface
	GetUrl() string
	Build(*testing.T) *Tester
}

// This method return an instance of Tester
//
// instance := supertest.NewHttpTester()
//
// The snapshot is used to create a test builder
func New() TesterInterface {
	return &Tester{
		httpHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		validation: validator.New(),
	}
}

// Using the function from set method from http request
//
// instance.Method(http.MethodGet)
//
// This return an instance of Tester
func (tt *Tester) Method(m string) TesterInterface {
	if m == "" {
		errors["method"] = "Method is required"
		return nil
	}
	tt.httpMethod = m

	return tt
}

// Using the function from set url from http request
//
// instance.Url("http://httpbin.org/get")
//
// This return an instance of Tester
func (tt *Tester) Url(url string) TesterInterface {
	if url == "" {
		errors["url"] = "Url is required"
		return tt
	}

	tt.httpUrl = url

	return tt
}

// Using the function from set query from http request
//
// instance.Query(map[string]string{"foo": "bar", "baz": "qux", "key": "value"})
//
// This return an instance of Tester
func (tt *Tester) Query(q map[string]string) TesterInterface {
	if q == nil {
		errors["query"] = "Query is required"
		return tt
	}

	var r string
	for k, v := range q {
		if v == "" {
			errors["query"] = "Query value is required"
			return tt
		}

		if r != "" {
			r += "&"
		} else {
			r = "?"
		}

		r += fmt.Sprintf("%s=%s", k, v)
	}

	tt.httpUrl += r

	return tt
}

// Using the function from set body return from http request
//
// instance.Json(&body)
//
// This return an instance of Tester
func (tt *Tester) Json(i interface{}) TesterInterface {
	if i == nil {
		errors["json"] = "Json is required"
		return tt
	}

	tt.response = i

	return tt
}

// Using the function from set status expected from http request
//
// instance.Status(http.StatusOK)
//
// This return an instance of Tester
func (tt *Tester) Status(s int) TesterInterface {
	if s == 0 {
		errors["status"] = "Status is required"
		return tt
	}
	if s < 100 && s > 599 {
		errors["status"] = "Status must be between 100 and 599"
		return tt
	}

	tt.httpStatus = s

	return tt
}

// Using the function from set headers from http request
// the default header as ben set is Content-Type: application/json
// instance.Headers(map[string]string{"Content-Type": "application/json"})
// This return an instance of Tester
func (tt *Tester) Headers(h map[string]string) TesterInterface {
	if h == nil {
		errors["headers"] = "Headers is required"
		return tt
	}

	for k, v := range h {
		if v == "" {
			errors["headers"] = "Header value is required"
			return tt
		}

		tt.httpHeaders[k] = v
	}

	return tt
}

// Using the function from set body from http request
//
// instance.Body([]byte(`{"title": "foo", "body": "bar", "userId": 1}`))
//
// This return an instance of Tester
func (tt *Tester) Body(bt []byte) TesterInterface {
	if bt == nil {
		errors["body"] = "Body is required"
		return tt
	}

	buff := bytes.NewBuffer(bt)

	tt.httpBody = buff

	return tt
}

func (tt *Tester) ValidateBody() TesterInterface {
	tt.validateBody = true

	return tt
}

// This return an instance of Tester expected testing
//
// instance.Method("GET").Url("http://httpbin.org/status/404").Status(404).Build(t)
//
// This return an instance of Tester
func (tt *Tester) Build(t *testing.T) *Tester {
	if !tt.validateStructRequestIsValid() {
		t.Errorf("Struct validation error: %s", errors["makeValidateStructRequest"])
	}

	tt.makeRequest()
	tt.makeResponse()
	tt.validateBodyReturned()

	if len(errors) > 0 {
		message := ""

		for k, v := range errors {
			message += k + ": " + v + "\n"
		}

		t.Error(message)
		return nil
	}

	return tt
}

func (tt *Tester) makeRequest() {
	r, err := http.NewRequest(tt.httpMethod, tt.httpUrl, tt.httpBody)
	if err != nil {
		errors["Error on create request"] = err.Error()
	}

	tt.makeHeaders(r)

	c := http.DefaultClient

	res, err := c.Do(r)
	if err != nil {
		errors["Error on execute request"] = err.Error()
		return
	}

	tt.checkWithStatusIsEqualExpected(res.StatusCode)

	tt.resp = res
}

func (tt *Tester) makeResponse() {
	if tt.response == nil {
		return
	}

	bt, err := io.ReadAll(tt.resp.Body)

	defer tt.resp.Body.Close()

	if err != nil {
		errors["Error on read response"] = err.Error()
		return
	}

	if err := json.Unmarshal(bt, tt.response); err != nil {
		errors["Error on unmarshal response"] = err.Error()
		return
	}
}

func (tt *Tester) makeHeaders(r *http.Request) {
	if len(tt.httpHeaders) == 0 {
		return
	}

	for k, v := range tt.httpHeaders {
		r.Header.Add(k, v)
	}
}

func (tt *Tester) validateStructRequestIsValid() bool {
	if err := tt.validation.Struct(tt); err != nil {
		errors["Make validate build invalid"] = err.Error()
		return false
	}

	return true
}

func (tt *Tester) validateBodyReturned() {
	if tt.validateBody && tt.response != nil {
		if err := tt.validation.Struct(tt.response); err != nil {
			errors["Body struct invalid"] = err.Error()
		}
	}
}

func (tt *Tester) checkWithStatusIsEqualExpected(status int) {
	if tt.httpStatus != status {
		str := fmt.Sprintf("Expected status: %d but got: %d", tt.httpStatus, status)
		errors["Error expected status"] = str
		return
	}
}

// GetUrl return the url of http request
//
// instance.GetUrl() == "http://httpbin.org/get"
func (tt *Tester) GetUrl() string {
	return tt.httpUrl
}

func (tt *Tester) GetStatus() int {
	return tt.httpStatus
}

func (tt *Tester) GetResponse() *http.Response {
	return tt.resp
}

func (tt *Tester) GetBody() interface{} {
	return tt.response
}

func (tt *Tester) GetHeaders() map[string]string {
	return tt.httpHeaders
}
