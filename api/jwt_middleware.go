package api

import (
	"github.com/auth0-community/go-auth0"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	jose "gopkg.in/square/go-jose.v2"
	"net/http"
)

var authenticatedEndpoints = []string{"/admin"}

// AuthConfig holds the auth configuration
type AuthConfig struct {
	ClientAudience string `envconfig:"AUTH_CLIENT_AUDIENCE"`
	ClientDomain   string `envconfig:"AUTH_CLIENT_DOMAIN"`
	ClientSecret   string `envconfig:"AUTH_CLIENT_SECRET"`
}

// GetConfig loads the config object from env vars and returns it
func GetConfig() (*AuthConfig, error) {
	var config AuthConfig

	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func authMiddleware(next http.Handler) http.Handler {
	config, err := GetConfig()
	if err != nil {
		log.Error("Unable to configure auth0")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte(config.ClientSecret)
		secretProvider := auth0.NewKeyProvider(secret)
		audience := []string{config.ClientAudience}
		configuration := auth0.NewConfiguration(secretProvider, audience, config.ClientDomain, jose.HS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			log.WithFields(log.Fields{
				"err":   err,
				"token": token,
			}).Error("Token is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
