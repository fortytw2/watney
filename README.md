Watney [![Build Status](https://travis-ci.org/fortytw2/watney.svg?branch=master)](https://travis-ci.org/fortytw2/watney) [![codecov](https://codecov.io/gh/fortytw2/watney/branch/master/graph/badge.svg)](https://codecov.io/gh/fortytw2/watney)
------

A port of the Ruby `vcr` library to Go, backed by `github.com/google/martian/har`

Run tests that depend on the network, reliably, without ever touching the
network itself.

Consider this package WIP - Pull Requests / Issues are greatly appreciated
while the exact API is worked out :)

### Usage

watney provides a `http.RoundTripper` that is automatically configured to use
fixtures unless `-watney` is set in the test flag.

First record fixtures using `go test -watney -v ./...`, then replay them by
omitting `-watney` in your future test runs.

To record a single test, use `go test -watney -v -run CaseName` and a new `.har`
will be written.

Bonus! You can inspect any `.har` individually, with any number of tools,
such as [HAR Viewer](http://www.softwareishard.com/har/viewer/), or [Har Viewer](https://ericduran.github.io/chromeHAR/)


```go
func TestGoogle(t *testing.T) {
	// you can use your own transport here, it is fully preserved during
	// recording
	tr := watney.Configure(http.DefaultTransport, t)

	c := &http.Client{
		Transport: tr,
	}

	// writes to a file like watney_test.go.har
	defer watney.Save(c)

	resp, err := c.Get("https://www.google.com")
	if err != nil {
		t.Fatal(err)
	}

	// test some things about resp

	defer resp.Body.Close()
}
```

LICENSE
------

Originally derived from a helper I found in `github.com/cardigann/cardigann`,
which is MIT licensed. The same license applies here
