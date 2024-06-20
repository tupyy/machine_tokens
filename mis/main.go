package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	port         int
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
	bodyData, err := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return OidcResponse{}, fmt.Errorf("oidc response status code %s: %s", resp.Status, string(bodyData))
	}

	oidcData := OidcResponse{}

	if err != nil {
		return OidcResponse{}, fmt.Errorf("error reading response body %s", err)
	}

	if err := json.Unmarshal(bodyData, &oidcData); err != nil {
		return OidcResponse{}, fmt.Errorf("error unmarshaling the response body %s", err)
	}

	log.Printf("User logged in: %+v", oidcData)
	return oidcData, nil
}

func main() {
	flag.IntVar(&port, "port", 8083, "MIS port")
	flag.StringVar(&keycloakUrl, "server_url", "http://127.0.0.1:3000", "Keycloak url")
	flag.StringVar(&clientId, "client_id", "vault", "Client id")
	flag.StringVar(&clientSecret, "client_secret", "vault", "Client secret")
	flag.StringVar(&username, "username", "cosmin", "Username")
	flag.StringVar(&password, "password", "cosmin", "Password")

	flag.Parse()

	log.Printf("Start MIS app. Listen on port %d", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		formData := map[string]string{
			"client_id":     clientId,
			"client_secret": clientSecret,
			"scope":         "openid",
			"grant_type":    "password",
			"username":      username,
			"password":      password,
		}
		oidcUrl := fmt.Sprintf("%s/realms/vault/protocol/openid-connect/token", strings.TrimRight(keycloakUrl, "/"))
		oidcResponse, err := auth(oidcUrl, formData)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, oidcResponse.AccessToken)
		return
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
