package logconsent

import (
	"github.com/ory/hydra/sdk/go/hydra"
	"github.com/ory/hydra/sdk/go/hydra/swagger"
	"github.com/ory/x/randx"
	"golang.org/x/oauth2"
	"strings"
)

type Hydra struct {
	Client hydra.SDK
	Conf   oauth2.Config
}

func CreateHydraClient(params map[string]string) *Hydra {
	client, err := hydra.NewSDK(&hydra.Configuration{
		AdminURL: params["admin_url"],
	})
	Handle(err)

	h := &Hydra{Client: client}
	h.updateConfig(params)

	return h
}

func (h *Hydra) updateConfig(params map[string]string) {
	scopes := strings.Split(params["scopes"], ",")

	h.Conf = oauth2.Config{
		ClientID:     params["client_id"],
		ClientSecret: params["client_secret"],
		Endpoint: oauth2.Endpoint{
			TokenURL: params["public_url"],
			AuthURL:  params["browser_url"],
		},
		RedirectURL: "http://localhost:3000/callback",
		Scopes:      scopes,
	}
}

func (h *Hydra) AcceptLogin(subject, challenge string) string {
	login_request := swagger.AcceptLoginRequest{
		Acr:         "oauth2",
		Remember:    false,
		RememberFor: 3600,
		Subject:     subject,
	}
	compl, _, err := h.Client.AcceptLoginRequest(challenge, login_request)

	Handle(err)
	return compl.RedirectTo
}

func (h *Hydra) AcceptConsent(id string, grantedScopes []string) string {
	compl, _, err := h.Client.AcceptConsentRequest(id, swagger.AcceptConsentRequest{
		GrantScope:  grantedScopes,
		Remember:    false,
		RememberFor: 3600,
	})
	Handle(err)

	return compl.RedirectTo
}

func (h *Hydra) AuthURL() string {
	state, _ := randx.RuneSequence(24, randx.AlphaLower)
	nonce, _ := randx.RuneSequence(24, randx.AlphaLower)
	authCodeURL := h.Conf.AuthCodeURL(
		string(state),
		oauth2.SetAuthURLParam("audience", ""),
		oauth2.SetAuthURLParam("nonce", string(nonce)),
		oauth2.SetAuthURLParam("prompt", ""),
		oauth2.SetAuthURLParam("max_age", "0"),
	)

	return authCodeURL
}
