package middleware

import (
	"context"
	"net/http"

	"chord-drawer/app/internal/session"
)

const (
	cookieSessionName = "session_id"
)

func SessionHandler(next http.HandlerFunc, mgr *session.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieSessionName)
		var sessionID string
		if err != nil || cookie == nil {
			sessionID, err = mgr.Create("")
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			cookie = &http.Cookie{
				Name:   cookieSessionName,
				Value:  sessionID,
				MaxAge: 0,
			}
			http.SetCookie(w, cookie)
		} else if mgr.Check(cookie.Value) {
			mgr.UpdateTime(cookie.Value)
		} else {
			sessionID, err = mgr.Create(cookie.Value)
		}
		if r.Method == http.MethodPost {
			ctx := context.WithValue(r.Context(), cookieSessionName, mgr.Get(cookie.Value).FilePath)
			next(w, r.WithContext(ctx))
		} else {
			next(w, r)
		}
	}
}
