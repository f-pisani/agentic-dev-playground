package feedbin

import (
	"net/url"
	"strconv"
)

type RequestOption func(v url.Values)

func WithPage(page int) RequestOption {
	return func(v url.Values) {
		v.Set("page", strconv.Itoa(page))
	}
}

func WithPerPage(perPage int) RequestOption {
	return func(v url.Values) {
		v.Set("per_page", strconv.Itoa(perPage))
	}
}