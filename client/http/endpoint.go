package http

import (
	"fmt"
	"net/url"
)

type endpoint struct {
	path string
}

func newBaseEndpoint(base string) func(string) *endpoint {
	return func(path string) *endpoint {
		return &endpoint{
			path: fmt.Sprintf("%s%s", base, path),
		}
	}
}

func (e endpoint) buildURL(values ...string) (*url.URL, error) {
	u, err := url.ParseRequestURI(fmt.Sprintf(e.path, values))
	if err != nil {
		return nil, err
	}

	return u, nil
}
