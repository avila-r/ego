package httpx

import (
	"net/url"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/list"
	"github.com/avila-r/ego/pair"
)

type Request struct {
	Target RequestTarget
	Method Method
	Body   Body
}

type RequestTarget struct {
	url   string
	query collection.List[pair.EntryPair[string, string]]
}

func Uri(uri string) *RequestTarget {
	return &RequestTarget{
		url:   uri,
		query: list.Empty[pair.EntryPair[string, string]](),
	}
}

func Url(url string) *RequestTarget {
	return &RequestTarget{
		url:   url,
		query: list.Empty[pair.EntryPair[string, string]](),
	}
}

func (rt *RequestTarget) WithQueryParam(key, value string) *RequestTarget {
	rt.query.Add(pair.EntryPair[string, string]{Key: key, Value: value})
	return rt
}

func (rt *RequestTarget) Build() (string, error) {
	u, err := url.Parse(rt.url)
	if err != nil {
		return "", err
	}

	q := u.Query()
	rt.query.ForEach(func(entry pair.EntryPair[string, string]) {
		q.Add(entry.Key, entry.Value)
	})

	u.RawQuery = q.Encode()

	return u.String(), nil
}
