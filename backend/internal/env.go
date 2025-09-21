package internal

var ENV struct {
	Server   `mapstructure:"server"`
	Database `mapstructure:"database"`
}

type Server struct {
	Port int `mapstructure:"port" defaultValue:"8000"`
}

type Database struct {
	Host     string `mapstructure:"host" defaultValue:"localhost"`
	Port     int    `mapstructure:"port" defaultValue:"9000"`
	User     string `mapstructure:"user" defaultValue:"postgres"`
	Password string `mapstructure:"password" defaultValue:"root"`
	Name     string `mapstructure:"name" defaultValue:"postgres"`
	Driver   string `mapstructure:"driver" defaultValue:"postgres"`
}
