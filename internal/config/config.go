package config

import (
	"fmt"
	"os"
	"strings"
)

type Client interface {
	ValidateEnvCred(field string, prefix string) string
}

func ValidateEnvCred(field string, prefix string) string {
	envKey := fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(field))

	value, _ := os.LookupEnv(envKey)

	if len(value) == 0 {
		fmt.Println("You have not provided the value for ", strings.ToUpper(field), ". We cannot continue.")
		os.Exit(1)
	}

	return value
}