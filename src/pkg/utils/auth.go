package utils

import (
	"athena/src/models/consts"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	OwnerId         uint
	OwnerType       string
	TokenUUID       string
	IsAuthenticated bool
	RequestUuid     string
}

func GetAuthData(c *gin.Context) *Auth {
	auth := &Auth{IsAuthenticated: false}
	id, exists := c.Get("authenticated-user-id")
	var authenticatedUserID uint
	if exists {
		authenticatedUserID = id.(uint)
	}
	authenticatedUserType := c.GetString("authenticated-user-type")
	accessTokenUUID := c.GetString("access-token-uuid")

	if authenticatedUserID != 0 && authenticatedUserType != "" && accessTokenUUID != "" {
		auth.IsAuthenticated = true
		auth.TokenUUID = accessTokenUUID
		auth.OwnerType = authenticatedUserType
		auth.OwnerId = authenticatedUserID
		auth.RequestUuid = c.GetString(string(consts.RequestUuid))
	}

	return auth
}

func GenerateRandomCodes(number int) ([]string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const codeLength = 5
	var randomCodes []string

	for i := 0; i < number; i++ {
		var part1, part2 string
		for j := 0; j < codeLength; j++ {
			randomIndex1, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			randomIndex2, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			part1 += string(charset[randomIndex1.Int64()])
			part2 += string(charset[randomIndex2.Int64()])
		}
		randomCode := fmt.Sprintf("%s-%s", part1, part2)
		randomCodes = append(randomCodes, randomCode)
	}

	return randomCodes, nil
}
