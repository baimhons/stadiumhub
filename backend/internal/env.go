package internal

var ENV struct {
	Server      `mapstructure:"server"`
	Database    `mapstructure:"database"`
	Redis       `mapstructure:"redis"`
	JWTSecret   `mapstructure:"jwtsecret"`
	APIFootball `mapstructure:"apifootball"`
}

type Server struct {
	Port int `mapstructure:"port" defaultValue:"8080"`
}

type Database struct {
	Host     string `mapstructure:"host" defaultValue:"localhost"`
	Port     int    `mapstructure:"port" defaultValue:"3306"`
	User     string `mapstructure:"user" defaultValue:"stadiumhubuser"`
	Password string `mapstructure:"password" defaultValue:"root"`
	Name     string `mapstructure:"name" defaultValue:"stadiumhub"`
	Driver   string `mapstructure:"driver" defaultValue:"mysql"`
}

type Redis struct {
	Host     string `mapstructure:"host" defaultValue:"localhost"`
	Port     int    `mapstructure:"port" defaultValue:"6379"`
	Password string `mapstructure:"password" defaultValue:""`
}

type JWTSecret struct {
	Secret string `mapstructure:"jwtsecret" defaultValue:"stadium_hub_secret"`
}

type APIFootball struct {
	APIKey string `mapstructure:"apikey" defaultValue:"e1009c5181884003ac55cd624a81502e"`
}
