package conf

type Database struct {
	Driver   string
	Name     string
	Port     uint
	Host     string
	User     string
	Password string
}

type AllowedHost struct {
	Host  string
	Https bool
}

type Server struct {
	ListenAddress           string
	AllowedHosts            []*AllowedHost
	UseForwardedProtoHeader bool
}

type General struct {
	SiteName string
}

type Config struct {
	Database Database
	Server   Server
	General  General
}

func (c *Config) CheckConfig() error {
	return nil
}
