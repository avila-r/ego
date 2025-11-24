package httpx

import "net/http"

func Get(url string) Response {
	var err error

	base, err := http.Get(url)
	if err != nil {
		return Response{err: err}
	}

	status := GetStatus(base.StatusCode)

	header := HeaderOf(base.Header)

	body, err := BodyOf(base.Body)

	response := Response{
		status: *status,
		header: header,
		body:   body,
		res:    base,
		err:    err,
	}

	return response
}
