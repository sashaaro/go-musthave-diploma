package auth

import (
	"context"
	"github.com/sashaaro/go-musthave-diploma/internal/config"
	privateRouter "github.com/sashaaro/go-musthave-diploma/internal/http/middlware/private_router"
	"net/http"
	"time"
)

func SetAuthCookie(w http.ResponseWriter, authToken string, tokenExp time.Duration) {
	//tokenExp = tokenExp
	cookie := http.Cookie{
		Name:  config.AuthCookie, // accessToken
		Value: authToken,
		Path:  "/",
		//MaxAge:   int(tokenExp),
		//HttpOnly: true,
		//Secure:   true,
		//SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func GetUserIDFromContext(ctx context.Context) int {
	userIDAny, _ := ctx.Value(privateRouter.KeyUserID).(int)
	return userIDAny
}
