package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateCardNumber() string {
	rand.Seed(time.Now().UnixNano())
	base := make([]int, 15)

	base[0] = 4
	for i := 1; i < 15; i++ {
		base[i] = rand.Intn(10)
	}

	sum := 0
	for i := 0; i < 15; i++ {
		n := base[i]
		if i%2 == 0 {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
	}
	checkDigit := (10 - (sum % 10)) % 10
	cardNumber := ""
	for _, d := range base {
		cardNumber += fmt.Sprintf("%d", d)
	}
	cardNumber += fmt.Sprintf("%d", checkDigit)
	return cardNumber
}

func GenerateCardExpiry() string {
	year := time.Now().Year() + 3
	month := rand.Intn(12) + 1
	return fmt.Sprintf("%02d/%d", month, year)
}
