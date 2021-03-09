package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/jwt"
	"net/http"
	"time"
)

var oauth2Provider fosite.OAuth2Provider

func clientCredentialsFlowHandler(rw http.ResponseWriter, r *http.Request) {

	context := r.Context()

	jwtClaims := &jwt.JWTClaims{}
	jwtClaims.Add("myCustomClaim", "myCustomClaimValue")
	var jwtSessionData = &oauth2.JWTSession{
		JWTClaims: jwtClaims,
	}
	accessRequest, err := oauth2Provider.NewAccessRequest(context, r, jwtSessionData)
	if err != nil {
		oauth2Provider.WriteAccessError(rw, accessRequest, err)
		return
	}

	response, err := oauth2Provider.NewAccessResponse(context, accessRequest)
	if err != nil {
		oauth2Provider.WriteAccessError(rw, accessRequest, err)
		return
	}

	oauth2Provider.WriteAccessResponse(rw, accessRequest, response)
}

func main() {

	var store = storage.NewExampleStore()
	var secret = []byte("Atlas Tiguan Jetta Passat Arteon")
	var config = &compose.Config{
		AccessTokenLifespan: time.Minute * 30,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("unable to create private key")
	}

	strategy := &compose.CommonStrategy{
		CoreStrategy:               compose.NewOAuth2JWTStrategy(privateKey, compose.NewOAuth2HMACStrategy(config, secret, nil)),
		OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(config, privateKey),
	}

	oauth2Provider = compose.Compose(config, store, strategy, nil, compose.OAuth2ClientCredentialsGrantFactory)

	http.HandleFunc("/token", clientCredentialsFlowHandler)
	http.HandleFunc("/k8s", k8sProbeHandler)
	fmt.Println("Listener : Started : Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func k8sProbeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("ok"))
}
