package bmber

import (
	"errors"
	"net/http"
	"net/url"
)

var (
	ErrFake503 = errors.New("Fake redirect to 503")
)

type HeaderTransport struct {
	Upstream http.RoundTripper
	Header   http.Header
}

func (t HeaderTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	for k, v := range t.Header {
		r.Header.Set(k, v)
	}

RT:
	res, err := t.RoundTrip(r)
	if ue, ok := err.(*url.Error); ok && ue.Err == ErrFake503 {
		goto TR
	}

	return res, err
}
