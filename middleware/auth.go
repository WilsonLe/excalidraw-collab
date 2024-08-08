package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/wilsonle/excalidraw-collab/constants"
	"github.com/wilsonle/excalidraw-collab/models"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(parts[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username, password := pair[0], pair[1]
		authenticated, err := models.AuthenticateUser(username, password)
		if err != nil || !authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), constants.USERNAME_CONTEXT_KEY, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
