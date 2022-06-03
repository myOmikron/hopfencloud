package crypt

import (
	"crypto/rand"
	"fmt"
	"github.com/myOmikron/hopfencloud/modules/logger"
)

func GetToken() (string, error) {
	var rbytes = make([]byte, 32)

	if _, err := rand.Read(rbytes); err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return fmt.Sprintf("%x", rbytes), nil
}
