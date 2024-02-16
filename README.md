# Supertest - Library from write http tests

Go code (golang) set of packages that provide many tools for testifying that your code will behave as you intend.

Get started:

- Install supertest with [one line of code](#installation), or [update it with another](#staying-up-to-date)
- A little about [Test-Driven Development (TDD)](https://en.wikipedia.org/wiki/Test-driven_development)

## [`supertest`]

The `supertest` package provides some helpful methods that allow you to write better http test code in Go.

- Prints friendly, easy to read failure descriptions
- Allows for very readable code
- Optionally annotate each assertion with a message

See it in action:

```go
package yours

import (
  "testing"
  "github.com/IcaroSilvaFK/supertest"
)


func TestShouldRequestWithParams(t *testing.T) {

	tt := supertest.New()

	tt.Method("GET").Url("http://httpbin.org/get").Status(200).Build(t)

}
```

- Every assert func takes the `testing.T` object as the first argument. This is how it writes the errors out through the normal `go test` capabilities.

```go
package yours

import (
  "testing"
	"github.com/IcaroSilvaFK/supertest"
)

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

```

---

# Supported go versions

We currently support the most recent major Go versions from 1.22 onward.

---

# Contributing

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue. Extra credit for those using Testify to write the test code that demonstrates it.

---

# License

This project is licensed under the terms of the MIT license.
