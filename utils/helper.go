package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func AbortWithErrJson(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, EJ(err))
}

func AbortWithStrJson(c *gin.Context, code int, str string) {
	c.AbortWithStatusJSON(code, SJ(str))
}

func EJ(e error) gin.H {
	return gin.H{
		"error": e.Error(),
	}
}

func SJ(str string) gin.H {
	return gin.H{
		"error": str,
	}
}

func StringToTime(c *gin.Context, t string) time.Time {
	nt, e := time.Parse(time.RFC3339, t)

	if e != nil {
		AbortWithStrJson(c, 500, "Time typecasting Error")
	}

	return nt
}

func GetBearerToken(h *string) (*string, error) {
	parts := strings.SplitN(*h, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("incorrect authorize token")
	}

	token := parts[1]

	return &token, nil
}

func InterfaceToStruct(src interface{}, dst interface{}) error {
	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonData, &dst); err != nil {
		return err
	}
	return nil
}
