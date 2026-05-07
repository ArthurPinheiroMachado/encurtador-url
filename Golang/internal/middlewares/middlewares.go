package middlewares

import (
	"encoding/base64"
	"fmt"
	"golang/internal/util"
	"net/http"
	"strings"
)

func Auth(
	user,
	pass string,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trace := util.CreateErrorContext("authMiddleware")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			util.SendHttpError(w, http.StatusUnauthorized, trace.Apply(fmt.Errorf("Unauthorized")))
			return
		}

		token := strings.TrimPrefix(authHeader, "Basic ")

		creds, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(creds), ":", 2)

		if len(parts) != 2 || parts[0] != user || parts[1] != pass {
			util.SendHttpError(w, http.StatusUnauthorized, fmt.Errorf("Token incorrect. Unauthorized"))
			return
		}

		next(w, r)
	}
}
