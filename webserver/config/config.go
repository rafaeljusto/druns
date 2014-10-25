package config

var (
	DrunsConfig Config
)

type Config struct {
	Database struct {
		Name string
		URI  string
	}
}
