package env

import libenv "github.com/caarlos0/env/v11"

var Var struct {
	OAuth2IntrospectURL string `env:"OAUTH2_INTROSPECT_URL"`
	OAuth2ClientID      string `env:"OAUTH2_CLIENT_ID" envDefault:"org.socialmaps.api"`
	OAuth2ClientSecret  string `env:"OAUTH2_CLIENT_SECRET"`
	MistralSecret       string `env:"MISTRAL_SECRET"`
	DatabaseDSN         string `env:"DATABASE_DSN" envDefault:"postgres://socialmaps_api@localhost:5432/socialmaps"`
	ListenAddr          string `env:"LISTEN_ADDR" envDefault:"127.0.0.1:8080"`
	NominatimEndpoint   string `env:"NOMINATIM_ENDPOINT" envDefault:"https://nominatim.openstreetmap.org"`
}

func init() {
	err := libenv.Parse(&Var)
	if err != nil {
		panic(err)
	}
}
