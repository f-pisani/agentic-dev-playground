package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"n
	"net/url"
	"fmt"
)

// GenerateExtractContentURL builds the signed URL for Feedbin's extract-full-content service.
func GenerateExtractContentURL(username, secret, targetURL string) (string, error) {
	digest := hmac.New(sha1.New, []byte(secret))
	_, err := digest.Write([]byte(targetURL))
	if err != nil {
		return "", err
	}
	signature := fmt.Sprintf("%x", digest.Sum(nil))
	base64URL := base64.URLEncoding.EncodeToString([]byte(targetURL))
	base64URL = strings.TrimRight(base64URL, "=") // urlsafe, no padding
	return (&url.URL{
		Scheme:   "https",
		Host:     "extract.feedbin.com",
		Path:     fmt.Sprintf("/parser/%s/%s", username, signature),
		RawQuery: "base64_url=" + base64URL,
	}).String(), nil
}
