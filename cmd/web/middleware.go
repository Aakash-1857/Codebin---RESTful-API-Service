// cmd/web/middleware.go
package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// recoverPanic is a middleware that recovers from panics and logs them.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic).
		defer func() {
			// Use the builtin recover function to check if there has been a panic.
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the serverErrorResponse helper to send a 500 error.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the value of the Authorization header.
		authorizationHeader := r.Header.Get("Authorization")

		// If the header is missing, send a 401 Unauthorized response.
		if authorizationHeader == "" {
			app.authenticationRequiredResponse(w, r)
			return
		}

		// The header value is expected to be in the format "Bearer <token>".
		// We split the string to isolate the token itself.
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		tokenString := headerParts[1]

		// Parse the JWT string.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// The signing method must be HMAC.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return our secret key.
			return []byte(app.config.jwt.secret), nil
		})

		// If there was an error parsing or the token is not valid, send a 401.
		if err != nil || !token.Valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// If the token is valid, call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}