package proxy

import (
	"net/http"
	"strings"

	"domain-proxy/src/base/define"
)

func getCookie(r *http.Request, name string) *http.Cookie {
	for _, c := range r.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	domain := strings.Split(strings.Replace(r.Host, "login.", "", 1), ":")[0]
	if r.URL.Path == "/logout" {
		http.SetCookie(w, &http.Cookie{Name: define.ServiceName, Value: "", MaxAge: 3600 * 24, Domain: domain})
		_, _ = w.Write([]byte("logout ok"))
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	user := r.URL.Query().Get("user")
	pass := r.URL.Query().Get("pass")
	if !verifyUserNameAndPassword(user, pass) {
		return
	}
	token := signToken(user)
	http.SetCookie(w, &http.Cookie{Name: define.ServiceName, Value: token, MaxAge: 3600 * 24, Domain: domain})
	log.Tracef("Domain: %s token: %s", domain, token)
	_, _ = w.Write([]byte("ok"))
}

func verifyLoginToken(w http.ResponseWriter, r *http.Request) bool {
	cookie := getCookie(r, define.ServiceName)
	if cookie == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	verify := verifyToken(cookie.Value)
	if !verify {
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}
