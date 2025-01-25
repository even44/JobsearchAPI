package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/even44/JobsearchAPI/pkg/initializers"
	"github.com/golang-jwt/jwt/v4"
)

var authLogger *log.Logger = log.New(os.Stdout, "[AUTH] ", log.Ldate+log.Ltime+log.Lmsgprefix)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get cookie
		tokenString, err := r.Cookie("Authorization")

		if err != nil {
			authLogger.Println("Could not extract auth cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(initializers.ApiSecret), nil
		})
		if err != nil {
			authLogger.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {

			// Check if token is expired
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				authLogger.Println("Token expired")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Check if user with provided id exists
			user, err := initializers.Store.GetById(int(claims["sub"].(float64)))
			if err != nil {
				authLogger.Printf("No user with id %d", claims["sub"])
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if user != nil {
				next.ServeHTTP(w, r)
			}
			w.WriteHeader(http.StatusUnauthorized)
			return

		} else {
			authLogger.Printf("Could not map claims of token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}
