package watney

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/google/martian/har"
)

// replayer is an http.RoundTripper that replays network interactions, allowing
// for perfectly reproducible network interactions
type replayer struct {
	entries []*har.Entry
}

func newReplayer(t TestingT) http.RoundTripper {
	_, file, num, ok := runtime.Caller(2)
	if !ok {
		fmt.Printf("called from %s#%d\n", file, num)
	}

	rt, err := newReplayerFromFile(file + ".har")
	if err != nil {
		t.Fatal("you must first record .har files with -watney")
	}

	return rt
}

// newReplayerFromFile crafts a new replayer from the given file
func newReplayerFromFile(path string) (*replayer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var h har.HAR
	if err = json.Unmarshal(b, &h); err != nil {
		return nil, err
	}

	return &replayer{entries: h.Log.Entries}, nil
}

func (r *replayer) matchEntry(req *http.Request) (*har.Entry, error) {
	if len(r.entries) == 0 {
		return nil, errors.New("No matching entry found")
	}

	entry := r.entries[0]
	r.entries = r.entries[1:]

	return entry, nil

}

func (r *replayer) RoundTrip(req *http.Request) (*http.Response, error) {
	entry, err := r.matchEntry(req)
	if err != nil {
		panic(err)
	}
	resp, err := createResponse(entry.Response)
	if err != nil {
		panic(err)
	}
	resp.Request = req
	return resp, nil
}

func readHTTPVersion(v string) (string, int, int) {
	if v == "HTTP/1.0" {
		return "HTTP/1.0", 1, 0
	}
	return "HTTP/1.1", 1, 1
}

func createResponse(hresp *har.Response) (*http.Response, error) {
	h := http.Header{}

	for _, hrow := range hresp.Headers {
		h.Add(hrow.Name, hrow.Value)
	}

	v, major, minor := readHTTPVersion(hresp.HTTPVersion)

	return &http.Response{
		Status:        hresp.StatusText,
		StatusCode:    hresp.Status,
		Header:        h,
		ContentLength: hresp.Content.Size,
		Body:          ioutil.NopCloser(bytes.NewReader(hresp.Content.Text)),
		Proto:         v,
		ProtoMajor:    major,
		ProtoMinor:    minor,
	}, nil
}
