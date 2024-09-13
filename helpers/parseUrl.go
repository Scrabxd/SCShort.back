package helpers

import "net/url"

func IsValidUrl(u string) bool {

	parsedUrl, err := url.Parse(u)
	if err != nil {
		return false
	}

	return parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
}
