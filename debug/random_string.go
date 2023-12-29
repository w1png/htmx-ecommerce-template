package debug

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func GenerateRandomString() string {
	hash_input := strconv.Itoa(rand.Intn(9999)) + strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(os.Getppid())
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(hash_input)))

	return hash
}
