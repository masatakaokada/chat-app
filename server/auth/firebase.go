package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

// JWTを検証する
func FirebaseMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Firebase SDK のセットアップ
			opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
			app, err := firebase.NewApp(context.Background(), nil, opt)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}

			client, err := app.Auth(context.Background())
			if err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}

			// クライアントから送られてきた JWT 取得
			auth := c.Request().Header.Get("Authorization")
			idToken := strings.Replace(auth, "Bearer ", "", 1)

			if idToken == "" {
				idToken = c.QueryParam("token")
			}

			// JWT の検証
			token, err := client.VerifyIDToken(context.Background(), idToken)
			if err != nil {
				fmt.Printf("error verifying ID token: %v\n", err)
				return c.String(http.StatusUnauthorized, "トークンの有効期限切れです")
			}

			c.Set("token", token)
			return next(c)
		}
	}
}
