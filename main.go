package main

import (
	"github.com/gorilla/sessions"
	lc "github.com/piensa/logico/logconsent"
	"net/http"
)

var (
	d     *lc.Database
	h     *lc.Hydra
	store = sessions.NewCookieStore([]byte("something-very-secret-keep-it-safe"))
)

const sessionName = "authentication"

func main() {
	d = &lc.Database{
		Dbuser: GetEnv("DB_USER", "dbuser"),
		Dbpw:   GetEnv("DB_PW", "secret"),
		Dbname: GetEnv("DB_NAME", "testusers"),
		Dbhost: GetEnv("DB_HOST", "localhost"),
		Dbport: GetEnv("DB_PORT", "5433"),
	}

	d.Connect()
	defer d.Db.Close()

	port := GetEnv("PORT", "3000")

	// Create hydra client.
	hydraConfig := map[string]string{
		"public_url":    GetEnv("HYDRA_PUBLIC_URL", "http://api.logi.co"),
		"admin_url":     GetEnv("HYDRA_ADMIN_URL", "http://admin.logi.co"),
		"client_id":     GetEnv("HYDRA_CLIENT_ID", "piensa"),
		"client_secret": GetEnv("HYDRA_CLIENT_SECRET", "piensa"),
		"scopes":        GetEnv("HYDRA_SCOPES", "openid,offline,eat,sleep,rave,repeat"),
	}
	h = lc.CreateHydraClient(hydraConfig)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/consent", consentHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.ListenAndServe("0.0.0.0:"+port, nil)

}
