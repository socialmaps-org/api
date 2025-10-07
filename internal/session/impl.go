package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type session struct {
	ID string `json:"id"`
}

type cookieValue struct {
	Session session `json:"session"`
}

// > Cookies with names starting with `__Host-Http-` must be set with
// > the `Secure` flag by a secure page (HTTPS) and must have the
// > `HttpOnly` attribute set to prove that they were set via the
// > `Set-Cookie` header. In addition, they also have the same
// > restrictions as `__Host-`-prefixed cookies [0]. This combination
// > yields a cookie that is as close as can be to treating the origin
// > as a security boundary while at the same time ensuring developers
// > and server operators know that its scope is limited to HTTP
// > requests.
//
// [0]
//
// > [Cookies with names starting with `__Host-`] must not have a
// > `Domain` attribute specified, and the `Path` attribute must be set
// > to `/`. This guarantees that such cookies are only sent to the host
// > that set them, and not to any other host on the domain. It also
// > guarantees that they are set host-wide and cannot be overridden on
// > any path on that host.
//
// Source:
// https://developer.mozilla.org/docs/Web/HTTP/Guides/Cookies#cookie_prefixes
const COOKIE_NAME = "__Host-Http-SocialMaps-Auth-Session"

func ToCookie(key []byte, sessionID string) *http.Cookie {
	val := cookieValue{Session: session{ID: sessionID}}

	msg, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}

	hash := hmac.New(sha256.New, key)
	_, err = hash.Write(msg)
	if err != nil {
		panic(err)
	}

	tag := hash.Sum(nil)

	var tagmsg []byte
	tagmsg = append(tagmsg, tag...)
	tagmsg = append(tagmsg, msg...)

	return &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    base64.URLEncoding.EncodeToString(tagmsg),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
}

func FromCookie(key []byte, cookie *http.Cookie) string {
	tagmsg, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		panic(err)
	}

	tag := tagmsg[:sha256.Size]
	msg := tagmsg[sha256.Size:]

	hash := hmac.New(sha256.New, key)
	_, err = hash.Write(msg)
	if err != nil {
		panic(err)
	}
	act := hash.Sum(nil)

	ok := hmac.Equal(act, tag)
	if !ok {
		panic("hmac invalid")
	}

	var val cookieValue
	err = json.Unmarshal(msg, &val)
	if err != nil {
		panic(err)
	}

	return val.Session.ID
}

func EmptyCookie() *http.Cookie {
	return &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	}
}
