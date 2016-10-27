package watney

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/martian/har"
)

// A recorder is an `http.Roundtripper` capable of recording and replaying
// network interactions
type recorder struct {
	http.RoundTripper
	hl *har.Logger
}

// newRecorder returns a new recorder object that fulfills the http.RoundTripper interface
func newRecorder(t http.RoundTripper) http.RoundTripper {
	hl := har.NewLogger()
	hl.SetOption(har.PostDataLogging(true))
	hl.SetOption(har.BodyLogging(true))

	return &recorder{
		RoundTripper: t,
		hl:           hl,
	}
}

// RoundTrip implements http.RoundTripper
func (r *recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())

	if err := r.hl.RecordRequest(id, req); err != nil {
		return nil, err
	}

	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if err = r.hl.RecordResponse(id, resp); err != nil {
		return resp, err
	}

	return resp, err
}

// Export returns the in-memory log.
func (r *recorder) export() *har.HAR {
	return r.hl.Export()
}
