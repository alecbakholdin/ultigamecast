package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv(filename string) {
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		vals := strings.Split(scanner.Text(), "=")
		if len(vals) < 2 {
			continue
		}
		os.Setenv(strings.TrimSpace(vals[0]), strings.Join(vals[1:], "="))
	}
}

func MustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Errorf("env variable %s cannot be empty", key))
	}
	return val
}