package param

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseIntParam(c *gin.Context, key string) (int64, error) {
	value := c.Param(key)
	if value == "" {
		return 0, errors.New("missing param: " + key)
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.New("invalid param: " + key)
	}
	return id, nil
}
