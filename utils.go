package main

import (
	"net/http"
	"os"
)

func GetEnv(env, fallback string) string {
	value := os.Getenv(env)

	if value == "" {
		return fallback
	}

	return value
}

func authenticated(r *http.Request) string {
	session, _ := store.Get(r, sessionName)
	if u, ok := session.Values["hydra-token"]; !ok {
		return ""
	} else if token, ok := u.(string); !ok {
		return ""
	} else {
		return token
	}
}
