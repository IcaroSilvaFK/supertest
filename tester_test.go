package supertest_test

import (
	"testing"

	"github.com/IcaroSilvaFK/supertest"
	"github.com/stretchr/testify/assert"
)

func TestShouldRequestWithParams(t *testing.T) {

	tt := supertest.New()

	tt.Method("GET").Url("http://httpbin.org/get").Status(200).Build(t)

}

func TestShouldExpectHttpErrorNotFoundOnApiNotExists(t *testing.T) {

	tt := supertest.New()

	tt.Method("GET").Url("http://httpbin.org/status/404").Status(404).Build(t)
}

func TestShouldRequestBodyNotEmpty(t *testing.T) {

	var body struct {
		UserId    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	tt := supertest.New()

	tt.Method("GET").Url("https://jsonplaceholder.typicode.com/todos/1").Json(&body).Status(200).Build(t)

	assert.NotNil(t, body)
	assert.Equal(t, 1, body.ID)
	assert.Equal(t, "delectus aut autem", body.Title)
	assert.Equal(t, false, body.Completed)
	assert.Equal(t, 1, body.UserId)
}

func TestShouldRequestPostOnApi(t *testing.T) {

	tt := supertest.New()

	var res struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		UserId int    `json:"userId"`
	}

	r := []byte(`{"title": "foo", "body": "bar", "userId": 1}`)

	tt.Method("POST").Url("https://jsonplaceholder.typicode.com/todos").Body(r).Json(&res).Status(201).Build(t)

	assert.Equal(t, "foo", res.Title)
	assert.Equal(t, "bar", res.Body)
	assert.Equal(t, 1, res.UserId)
}

func TestShouldRequestNotPostOnRouteNotFound(t *testing.T) {

	tt := supertest.New()

	r := []byte(`{"title": "foo", "body": "bar", "userId": 1}`)

	tt.Method("POST").Url("https://jsonplaceholder.typicode.com/tod").Body(r).Status(404).Build(t)
}

func TestShouldUrlMatchOnAddQueryParams(t *testing.T) {

	tt := supertest.New()

	tt.Url("http://httpbin.org/get").Query(map[string]string{"foo": "bar", "baz": "qux", "key": "value"})

	url := tt.GetUrl()

	assert.Equal(t, "http://httpbin.org/get?foo=bar&baz=qux&key=value", url)
}

func TestShouldValidateBody(t *testing.T) {

	var body struct {
		UserId    int    `json:"userId" validate:"required"`
		ID        int    `json:"id" validate:"required"`
		Title     string `json:"title" validate:"required"`
		Completed bool   `json:"completed" validate:"required"`
	}

	tt := supertest.New()

	tt.Method("GET").Url("https://jsonplaceholder.typicode.com/todos/1").Json(&body).Status(200).ValidateBody().Build(t)
}
