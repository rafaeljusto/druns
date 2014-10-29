package config

var (
	DrunsConfig Config
)

type Config struct {
	Server struct {
		IP   string
		Port string
	}

	Database struct {
		Name string
		URI  string
	}
}
