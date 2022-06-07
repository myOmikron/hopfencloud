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
	CLISockPath             string
	ListenAddress           string
	AllowedHosts            []*AllowedHost
	UseForwardedProtoHeader bool
}

type Files struct {
	DataPath string
}

type Config struct {
	Database Database
	Files    Files
	Server   Server
}

func (c *Config) CheckConfig() error {
	return nil
}
