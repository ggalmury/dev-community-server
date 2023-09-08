package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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

func ErrHandledMarshal(c *gin.Context, i interface{}) []byte {
	v, e := json.Marshal(i)

	if e != nil {
		AbortWithStrJson(c, 500, "Marshal Error")
	}

	return v
}

func ErrHandledUnmarshal[T any](c *gin.Context, b []byte) T {
	var v T

	e := json.Unmarshal(b, &v)

	if e != nil {
		AbortWithStrJson(c, 500, "Marshal Error")
	}

	return v
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
