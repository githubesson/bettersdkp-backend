package middleware

import (
	"net/http"
	"strings"
)

func ExtractSessionID(r *http.Request) string {
	var sessionID string

	sessionID = r.Header.Get("X-Session-ID")

	if sessionID == "" {
		cookie, err := r.Cookie("JSESSIONID")
		if err == nil {
			sessionID = cookie.Value
		} else {

			cookieHeader := r.Header.Get("Cookie")
			if cookieHeader != "" {
				cookies := strings.Split(cookieHeader, ";")
				for _, c := range cookies {
					parts := strings.Split(strings.TrimSpace(c), "=")
					if len(parts) == 2 && parts[0] == "JSESSIONID" {
						sessionID = parts[1]
						break
					}
				}
			}
		}
	}

	return sessionID
}
