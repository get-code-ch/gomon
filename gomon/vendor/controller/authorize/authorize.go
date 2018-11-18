package authorize

import (
	"github.com/gorilla/sessions"
)

const UserContext = "user-context"

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("1FE3B7IIB05CGBK2F6D17KJ61H36OLJJ")
	Store = sessions.NewCookieStore(key)
)
