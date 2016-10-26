package watney_test

import (
	"net/http"
	"testing"

	"github.com/fortytw2/watney"
)

func TestWatney(t *testing.T) {
	tr := watney.Configure(http.DefaultTransport, t)

	c := &http.Client{
		Transport: tr,
	}

	resp, err := c.Get("https://www.google.com")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	watney.Save(c)
	println("saved file")
}
