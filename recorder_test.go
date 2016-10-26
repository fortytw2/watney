package watney

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHarRecorder(t *testing.T) {
	var responseStr = "Hello, client\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, responseStr)
	}))
	defer ts.Close()

	tr := newRecorder(http.DefaultTransport)
	c := &http.Client{
		Transport: tr,
	}

	res, err := c.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if string(greeting) != responseStr {
		t.Fatalf(`Expected %q, got %q`, responseStr, greeting)
	}

}
