package auth

import (
	"strconv"
	"strings"
	"time"

	"math/rand"
)

func GeneratePasscode() string {
	// get current ms
	curMs := time.Now().Nanosecond() / 1000

	// convert ms to str and get first 4 char
	msStr := strconv.Itoa(curMs)[:4]

	// generate random char between A and Z
	var alphb []int
	for i := 0; i < 4; i++ {
		alphb = append(alphb, rand.Intn(26)+65)
	}

	// Convert ascii values to character and join them
	var alphChar []string
	for _, a := range alphb {
		alphChar = append(alphChar, string(rune(a)))
	}
	alphStr := strings.Join(alphChar, "")

	// combine alphabet string and ms string
	return alphStr + msStr
}
