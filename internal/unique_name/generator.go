package uniquename

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateUniqueName() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	if err != nil {
		panic(err)
	}
	randomAdjective := adjectives[n.Int64()]
	n, err = rand.Int(rand.Reader, big.NewInt(int64(len(colors))))
	if err != nil {
		panic(err)
	}
	randomColor := colors[n.Int64()]
	n, err = rand.Int(rand.Reader, big.NewInt(int64(len(animals))))
	if err != nil {
		panic(err)
	}
	randomAnimal := animals[n.Int64()]

	return fmt.Sprintf("%s-%s-%s", randomAdjective, randomColor, randomAnimal)
}
