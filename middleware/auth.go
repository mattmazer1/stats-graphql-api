package middleware

import (
	"context"
	"net/http"

	db "github.com/mattmazer1/graphql-api/database"
	"github.com/mattmazer1/graphql-api/graph/model"
	"github.com/mattmazer1/graphql-api/utils"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate jwt token
			tokenStr := header
			username, err := utils.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// Create user and check if user exists in db
			user := model.User{Username: username}
			id, err := db.GetUserId(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			user.ID = id
			// Put it in context
			ctx := context.WithValue(r.Context(), UserCtxKey, &user)

			// Call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(UserCtxKey).(*model.User)
	return raw
}
