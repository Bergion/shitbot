package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

type ContextKey string

const key ContextKey = "user"

type TelegramUser struct {
	ID              int64  `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Username        string `json:"username"`
	LanguageCode    string `json:"language_code"`
	IsPremium       bool   `json:"is_premium"`
	AllowsWriteToPm bool   `json:"allows_write_to_pm"`
}

func TelegramAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		telegramBotToken := os.Getenv("TELEGRAM_TOKEN")
		authCheckString, err := getAuthCheckString(r.URL.Query())
		if err != nil {
			log.Println(err)
			authFailed(w)
			return
		}

		secretKey := getHmac256Signature([]byte("WebAppData"), []byte(telegramBotToken))
		expectedHashString := hex.EncodeToString(getHmac256Signature(secretKey, []byte(authCheckString)))
		fmt.Println(expectedHashString)
		if expectedHashString != r.URL.Query().Get("hash") {
			log.Println("Invalid data. Hashes are not equal")
			authFailed(w)
			return
		}

		tgUser := &TelegramUser{}
		err = json.Unmarshal([]byte(r.URL.Query().Get("user")), tgUser)
		if err != nil {
			log.Println(err)
			authFailed(w)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), key, tgUser))
		next.ServeHTTP(w, r)
	})
}

func FromContext(ctx context.Context) (*TelegramUser, bool) {
	user, ok := ctx.Value(key).(*TelegramUser)
	return user, ok
}

func getAuthCheckString(values url.Values) (string, error) {
	paramKeys := make([]string, 0)
	for key, v := range values {
		if key == "hash" {
			continue
		}
		if len(v) != 1 {
			return "", errors.New("is not a valid auth query")
		}
		paramKeys = append(paramKeys, key)
	}

	// sort keys
	sort.Strings(paramKeys)

	dataCheckArr := make([]string, len(paramKeys))
	for i, key := range paramKeys {
		dataCheckArr[i] = key + "=" + values.Get(key)
	}

	return strings.Join(dataCheckArr, "\n"), nil
}

func getHmac256Signature(secretKey []byte, data []byte) []byte {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(data)
	sum := mac.Sum(nil)
	return sum
}

func authFailed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}
