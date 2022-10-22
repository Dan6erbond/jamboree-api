package randomstring

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomString(length int) string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	charactersLength := len(characters)
	result := ""
	for i := 0; i <= length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(charactersLength)))
		if err != nil {
			panic(err)
		}
		result += string(characters[n.Int64()])
	}
	return result
}
