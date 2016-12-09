package watney

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

var gold bool

func init() {
	flag.BoolVar(&gold, "watney", false, "record new watney fixtures")
	flag.Parse()
}

type TestingT interface {
	Fatal(args ...interface{})
}

// Configure reconfigures a given http.Client to use a fixture provided transport
func Configure(tr http.RoundTripper, t TestingT) http.RoundTripper {
	if gold {
		return newRecorder(tr)
	}
	return newReplayer(t)
}

// Save writes the interactions to a file, if `-watney` is set for the test run
func Save(c *http.Client) {
	if !gold {
		return
	}

	r, ok := c.Transport.(*recorder)
	if !ok {
		panic(c.Transport)
	}

	_, file, num, ok := runtime.Caller(1)
	if !ok {
		fmt.Printf("called from %s#%d\n", file, num)
	}

	f, err := os.Create(file + ".har")
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(f).Encode(r.export())
	if err != nil {
		panic(err)
	}
}
