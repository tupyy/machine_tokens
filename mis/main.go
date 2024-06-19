package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	keycloakUrl  string
	clientId     string
	clientSecret string
	username     string
	password     string
)

type MisResponse struct {
	Role string `json:"role"`
	Jwt  string `json:"jwt"`
}

type OidcResponse struct {
	AccessToken string `json:"access_token"`
}

func auth(serverUrl string, formData map[string]string) (OidcResponse, error) {
	data := url.Values{}
	for k, v := range formData {
		data.Add(k, v)
	}
	resp, err := http.PostForm(serverUrl, data)
	if err != nil {
		return OidcResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return OidcResponse{}, fmt.Errorf("oidc response status code %s", resp.Status)
	}

	oidcData := OidcResponse{}

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return OidcResponse{}, fmt.Errorf("error reading response body %s", err)
	}

	if err := json.Unmarshal(bodyData, &oidcData); err != nil {
		return OidcResponse{}, fmt.Errorf("error unmarshaling the response body %s", err)
	}

	return oidcData, nil
}

func main() {
	flag.StringVar(&keycloakUrl, "server_url", "http://127.0.0.1:3000", "Keycloak url")
	flag.StringVar(&clientId, "client_id", "vault", "Client id")
	flag.StringVar(&clientSecret, "client_secret", "vault", "Client secret")
	flag.StringVar(&username, "username", "cosmin", "Username")
	flag.StringVar(&password, "password", "cosmin", "Password")

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
}
