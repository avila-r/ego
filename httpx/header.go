package httpx

import (
	"net/http"

	"github.com/avila-r/ego/collection"
	"github.com/avila-r/ego/list"
	"github.com/avila-r/ego/optional"
)

type Header struct {
	headers map[string][]string
}

func (h *Header) First(key string) optional.Optional[string] {
	return optional.Of("")
}

func (h *Header) All(key string) collection.List[string] {
	return list.Empty[string]()
}

func HeaderOf(header http.Header) Header {
	return Header{} // Placeholder
}
