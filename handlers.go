package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("error")) > 0 {
		fmt.Printf("Got error: %s\n", r.URL.Query().Get("error_description"))

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<html><body><h1>An error occurred</h1><h2>%s</h2><p>%s</p><p>%s</p><p>%s</p></body></html>", r.URL.Query().Get("error"), r.URL.Query().Get("error_description"), r.URL.Query().Get("error_hint"), r.URL.Query().Get("error_debug"))
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.Conf.Exchange(context.Background(), code)
	if err != nil {
		log.Print(err)
		http.Error(w, "token exchange error", http.StatusBadRequest)
		return
	}

	// Add token to user session.
	session, _ := store.Get(r, sessionName)
	session.Values["hydra-token"] = token.AccessToken
	session.Values["refresh-token"] = token.RefreshToken

	// Store the session in the cookie
	if err := store.Save(r, w, session); err != nil {
		http.Error(w, "Could not persist cookie", http.StatusBadRequest)
		return
	}

	fmt.Println(token.AccessToken)
	http.Redirect(w, r, "/", 301)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)
	session.Options.MaxAge = -1
	http.Redirect(w, r, "/", 301)
}

func consentHandler(w http.ResponseWriter, r *http.Request) {
	consentRequestId := r.URL.Query().Get("consent_challenge")
	if consentRequestId == "" {
		http.Error(w, "Consent endpoint was called without a consent request id", http.StatusBadRequest)
		return
	}

	consentRequest, response, err := h.Client.GetConsentRequest(consentRequestId)
	if err != nil {
		http.Error(w, "The consent request endpoint does not respond", http.StatusBadRequest)
		return
	} else if response.StatusCode != http.StatusOK {
		http.Error(w, "Consent request endpoint", http.StatusBadRequest)
		return
	}

	// TODO: Check if user is authenticated.
	if r.Method == "POST" {
		// Parse the form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		var grantedScopes = []string{}
		for key := range r.PostForm {
			grantedScopes = append(grantedScopes, key)
		}

		redirect := h.AcceptConsent(consentRequestId, grantedScopes)
		http.Redirect(w, r, redirect, http.StatusFound)
	}

	// TODO: Check POST method.

	values := map[string]interface{}{
		"Scopes":   consentRequest.RequestedScope,
		"ClientID": consentRequest.Client.ClientId,
	}

	t, _ := template.New("consent.html").ParseFiles("templates/consent.html")
	t.Execute(w, &values)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	access_token := authenticated(r)
	if access_token == "" {
		// Perform hydra auth workflow.
		authURL := h.AuthURL()
		w.Write([]byte(fmt.Sprintf(`
		<html>
			<head></head>
			<body>
				<h2>You are not logged in</h2>
				<p><a href="%s">Authorize application</a></p>
			</body>
		</html>`, authURL)))
		return
	}

	// Validate token with Hydra server.
	values := map[string]string{
		"proxy_url":    GetEnv("PROXY_URL", "http://okproxy.logi.co"),
		"access_token": access_token,
	}

	t, err := template.New("index.html").ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, &values)

	// w.Write([]byte(fmt.Sprintf(`
	// <html>
	// 	<head></head>
	// 	<body>
	// 		<p>You are now logged in and token %s is Active!</p>
	// 	</body>
	// </html>`, access_token)))

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("login_challenge")

	values := map[string]string{
		"Challenge": challenge,
		"Error":     "",
	}
	t, _ := template.New("login.html").ParseFiles("templates/login.html")
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Could not parse form", 500)
			return
		}

		// Checking user credentials
		email := r.Form.Get("email")
		pw := r.Form.Get("password")
		if !d.Validate(email, pw) {
			values["Error"] = "Incorrect password"
			t.Execute(w, &values)
			return
		}

		redirect := h.AcceptLogin(email, challenge)
		http.Redirect(w, r, redirect, http.StatusFound)
	}

	t.Execute(w, &values)
}
