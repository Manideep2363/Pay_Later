package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {

        currentRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "role not found",
            })
            c.Abort()
            return
        }

        allowed := false

        for _, role := range roles {
            if currentRole.(string) == role {
                allowed = true
                break
            }
        }

        if !allowed {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "access denied",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}