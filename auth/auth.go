package auth

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetToken Получает токены из файла по заданному названию токена
func GetToken(tag string) (TgBotToken string, err error) {
	defer func() {
		if p := recover(); p != nil {
			log.Fatalf("Panic! can't read %s from secret file", tag)
		}
		if err != nil {
			err = errors.Wrapf(err, "can't read %s from secret file", tag)
		}
	}()

	path := filepath.Join("secrets", "tokens.inf")

	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	configLines := strings.Split(string(file), "\n")

	for i, str := range configLines {
		if strings.Contains(str, tag) {
			TgBotToken = strings.TrimSpace(configLines[i+1])
			break
		}
	}

	return
}
