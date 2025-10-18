package env

type Env struct {
	// Host is the host + port
	//
	// Examples:
	//
	// * 127.0.0.1:8080
	// * boramalper.org
	Host string `env:"HOST"`
	Post int    `env:"PORT"`
}
