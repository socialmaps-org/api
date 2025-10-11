package env

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
}
