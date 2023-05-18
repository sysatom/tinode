package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
	"unicode"
)

func HasHan(txt string) bool {
	for _, runeValue := range txt {
		if unicode.Is(unicode.Han, runeValue) {
			return true
		}
	}
	return false
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

const (
	UrlRegex = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
)

func IsUrl(text string) bool {
	re := regexp.MustCompile("^" + UrlRegex + "$")
	return re.MatchString(text)
}

func Masker(input string, start int) string {
	if len(input) <= start {
		return input
	}
	lenStart := len(input[start:])
	switch {
	case lenStart <= 3:
		return input[:start] + strings.Repeat("*", lenStart)
	case 3 < lenStart && lenStart <= 5:
		return input[:start+1] + strings.Repeat("*", lenStart-2) + input[lenStart+start-1:]
	case 5 < lenStart && lenStart <= 10:
		return input[:start+2] + strings.Repeat("*", lenStart-4) + input[lenStart+start-2:]
	case lenStart > 10:
		return input[:start+4] + strings.Repeat("*", lenStart-8) + input[lenStart+start-4:]
	default:
		return ""
	}
}

func Fn(public interface{}) string {
	if v, ok := public.(map[string]interface{}); ok {
		if s, ok := v["fn"]; ok {
			if ss, ok := s.(string); ok {
				return ss
			}
		}
	}
	return ""
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
