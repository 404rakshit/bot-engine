package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// VerifyTelegramAuth validates the payload authenticity and checks for replay attacks.
func VerifyTelegramAuth(data map[string]string, botToken string) error {

	// utils.Print("format: %+v", data)

	hash, hasHash := data["hash"]
	if !hasHash {
		return fmt.Errorf("missing hash parameter")
	}

	// 1. Replay Attack Prevention (Check timestamp)
	authDateStr, hasDate := data["auth_date"]
	if !hasDate {
		return fmt.Errorf("missing auth_date parameter")
	}

	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid auth_date format")
	}

	// Reject requests older than 24 hours (86400 seconds)
	if time.Now().Unix()-authDate > 86400 {
		return fmt.Errorf("authentication payload expired")
	}

	// 2. Build the data_check_string
	var keys []string
	for k := range data {
		if k != "hash" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys) // MUST be alphabetical

	var dataCheckArr []string
	for _, k := range keys {
		// Format: key=value
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, data[k]))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// 3. Create the secret key (SHA256 of the bot token)
	secretHash := sha256.New()
	secretHash.Write([]byte(botToken))
	secretKey := secretHash.Sum(nil)

	// 4. Calculate HMAC-SHA256 signature
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	// 5. Securely compare hashes
	if !hmac.Equal([]byte(expectedHash), []byte(hash)) {
		return fmt.Errorf("data integrity check failed: invalid signature")
	}

	return nil
}
