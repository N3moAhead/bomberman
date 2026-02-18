package router

import (
	"context"
	"net/http"

	"github.com/N3moAhead/bomberman/website/internal/db"
)

func userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, appSessionName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userIDValue, ok := session.Values["user_id"]
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		user, err := db.GetUserByID(userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
