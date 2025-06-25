package feedbin

import "net/http"

type Authenticator interface {
	Authenticate(req *http.Request)
}

type BasicAuth struct {
	Username string
	Password string
}

func (b *BasicAuth) Authenticate(req *http.Request) {
	req.SetBasicAuth(b.Username, b.Password)
}