package httpx

import (
	"net/http"
)

type Response struct {
	status Status
	header Header
	body   Body

	res *http.Response
	err error
}

func (r *Response) Status() Status {
	return r.status
}

func (r *Response) Header() Header {
	return r.header
}

func (r *Response) Body() Body {
	return r.body
}

func (r *Response) Error() error {
	return r.err
}

func (r *Response) HasFailed() bool {
	return r.err != nil
}

func (r *Response) Raw() *http.Response {
	return r.res
}
