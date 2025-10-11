package env

import "time"

type AuthEnv struct {
	// Host is the host + port
	//
	// Examples:
	//
	// * 127.0.0.1:8080
	// * boramalper.org
	Host                 string             `env:"HOST"`
	OSMClientID          string             `env:"OSM_CLIENT_ID"`
	OSMClientSecret      string             `env:"OSM_CLIENT_SECRET"`
	CookieSecret         Base64EncodedBytes `env:"COOKIE_SECRET"`
	OAuth2Secret         Base64EncodedBytes `env:"OAUTH2_SECRET"`
	Oauth2PrivateKeyFile string             `env:"OAUTH2_PRIVATE_KEY_FILE"`

	// OAuth2 Token Lifespans
	AccessTokenLifespan  time.Duration `env:"OAUTH2_ACCESS_TOKEN_LIFESPAN" envDefault:"30m"`
	RefreshTokenLifespan time.Duration `env:"OAUTH2_REFRESH_TOKEN_LIFESPAN" envDefault:"720h"` // 30 days
	AuthCodeLifespan     time.Duration `env:"OAUTH2_AUTH_CODE_LIFESPAN" envDefault:"10m"`
	IDTokenLifespan      time.Duration `env:"OAUTH2_ID_TOKEN_LIFESPAN" envDefault:"6h"`
}
