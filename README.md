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
  "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
  assert := assert.New(t)

  // assert equality
  assert.Equal(123, 123, "they should be equal")

  // assert inequality
  assert.NotEqual(123, 456, "they should not be equal")

  // assert for nil (good for errors)
  assert.Nil(object)

  // assert for not nil (good when you expect something)
  if assert.NotNil(object) {

    // now we know that object isn't nil, we are safe to make
    // further assertions without causing any errors
    assert.Equal("Something", object.Value)
  }
}
```

---

# Staying up to date

To update Testify to the latest version, use `go get -u github.com/IcaroSilvaFK/supertest`.

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
