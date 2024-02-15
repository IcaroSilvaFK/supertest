package supertest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type Tester struct {
	httpUrl     string
	httpHeaders map[string]string
	httpBody    io.Reader
	httpMethod  string
	httpStatus  int
	response    interface{}
	resp        *http.Response
}

var errors = make(map[string]string, 0)

type TesterInterface interface {
	Method(string) TesterInterface
	Url(string) TesterInterface
	Json(interface{}) TesterInterface
	Headers(map[string]string) TesterInterface
	Body([]byte) TesterInterface
	Status(int) TesterInterface
	Build(*testing.T) *Tester
}

func NewHttpTester() TesterInterface {
	return &Tester{
		httpHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (tt *Tester) Method(m string) TesterInterface {
	if m == "" {
		errors["method"] = "Method is required"
		return nil
	}
	tt.httpMethod = m

	return tt
}

func (tt *Tester) Url(url string) TesterInterface {
	if url == "" {
		errors["url"] = "Url is required"
		return tt
	}

	tt.httpUrl = url

	return tt
}

func (tt *Tester) Json(i interface{}) TesterInterface {
	if i == nil {
		errors["json"] = "Json is required"
		return tt
	}

	tt.response = i

	return tt
}

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

func (tt *Tester) Body(bt []byte) TesterInterface {
	if bt == nil {
		errors["body"] = "Body is required"
		return tt
	}

	buff := bytes.NewBuffer(bt)

	tt.httpBody = buff

	return tt
}

func (tt *Tester) Build(t *testing.T) *Tester {

	tt.makeRequest()
	tt.makeResponse()

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
		errors["makeRequest"] = err.Error()
	}

	tt.makeHeaders(r)

	c := http.DefaultClient

	res, err := c.Do(r)

	if err != nil {
		errors["makeRequest"] = err.Error()
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
		errors["makeResponse"] = err.Error()
		return
	}

	if err := json.Unmarshal(bt, tt.response); err != nil {
		errors["makeResponse"] = err.Error()
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

func (tt *Tester) checkWithStatusIsEqualExpected(status int) {

	if tt.httpStatus != status {
		str := fmt.Sprintf("Expected status: %d but got: %d", tt.httpStatus, status)
		errors["checkWithStatusIsEqualExpected"] = str
		return
	}
}
