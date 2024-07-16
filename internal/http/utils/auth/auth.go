package auth

import (
	"context"
	"github.com/sashaaro/go-musthave-diploma-tpl/internal/config"
	"net/http"
	"time"
)

func SetAuthCookie(w http.ResponseWriter, authToken string, tokenExp time.Duration) {
	//tokenExp = tokenExp
	cookie := http.Cookie{
		Name:  config.AUTH_COOKIE, // accessToken
		Value: authToken,
		Path:  "/",
		//MaxAge:   int(tokenExp),
		//HttpOnly: true,
		//Secure:   true,
		//SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func GetUserIdFromContext(ctx context.Context) *int {
	userIdAny, ok := ctx.Value("userId").(int)
	if ok {
		return &userIdAny
	}

	return nil
}
