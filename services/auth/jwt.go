package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"mohaafaridd.dev/ecom/config"
	"mohaafaridd.dev/ecom/types"
	"mohaafaridd.dev/ecom/utils"
)

type contextKey string

const UserKey contextKey = "userId"

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from authorization header
		tokenString := getTokenFromRequest(r)
		// Validate JWT
		token, err := validateToken(tokenString)

		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			PermissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			PermissionDenied(w)
			return
		}

		// Fetch user
		claims := token.Claims.(jwt.MapClaims)
		userIdStringified := claims["userId"].(string)
		userId, _ := strconv.Atoi(userIdStringified)
		u, err := store.GetUserById(userId)

		if err != nil {
			log.Println("invalid token")
			PermissionDenied(w)
			return
		}

		// set context "userId"
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if strings.HasPrefix(token, "Bearer") {
		token = strings.Replace(token, "Bearer ", "", 1)
		return token
	}

	return token
}

func validateToken(token string) (*jwt.Token, error) {
	log.Println(token)
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func PermissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int)

	if !ok {
		return -1
	}

	return userId
}
