package simple_gin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthProcessor struct {
}

func (ap *AuthProcessor) Process(c *gin.Context) (*gin.Context, error) {
	token := c.GetHeader("token")
	if token != "token" {
		err := errors.New("invalid token")
		c.String(http.StatusUnauthorized, "invalid token")
		return nil, err
	}

	return c, nil
}
