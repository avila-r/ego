package httpx

import "io"

type Body struct {
	bytes []byte
	empty bool
}

func EmptyBody() Body {
	return Body{
		bytes: []byte{},
		empty: true,
	}
}

func BodyOf(data io.ReadCloser) (Body, error) {
	defer data.Close()

	bytes, err := io.ReadAll(data)
	if err != nil {
		return EmptyBody(), err
	}

	return Body{bytes: bytes}, nil
}

func (b Body) Bytes() []byte {
	return b.bytes
}

func (b Body) IsEmpty() bool {
	return b.empty
}
