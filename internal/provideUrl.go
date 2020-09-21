package internal

import "net/url"

func CreateUrl(s string) func() (*url.URL, error) {
	return func() (*url.URL, error) {
		return url.Parse(s)
	}
}

