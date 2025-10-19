package internal

var ENV struct {
	Server      `mapstructure:"server"`
	Database    `mapstructure:"database"`
	Redis       `mapstructure:"redis"`
	APIFootball `mapstructure:"apifootball"`
	Stripe      `mapstructure:"stripe"`
}

type Server struct {
	Port int `mapstructure:"port" defaultValue:"8080"`
}

type Database struct {
	Host     string `mapstructure:"host" defaultValue:"mysql"`
	Port     int    `mapstructure:"port" defaultValue:"3306"`
	User     string `mapstructure:"user" defaultValue:"stadiumhubuser"`
	Password string `mapstructure:"password" defaultValue:"root"`
	Name     string `mapstructure:"name" defaultValue:"stadiumhub"`
	Driver   string `mapstructure:"driver" defaultValue:"mysql"`
}

type Redis struct {
	Host     string `mapstructure:"host" defaultValue:"redis"`
	Port     int    `mapstructure:"port" defaultValue:"6379"`
	Password string `mapstructure:"password" defaultValue:""`
}

type APIFootball struct {
	APIKey string `mapstructure:"apikey" defaultValue:"e1009c5181884003ac55cd624a81502e"`
}

type Stripe struct {
	StripeKey string `mapstructure:"stripeKey" defaultValue:"sk_test_51SJdxYBVJneBwH8Vw62Ms5memSg9TXmtovK2ZhGCwkIa8VYRDA3sHyEQRH5g7SNgSMwOkuDNdBjWUeiWAR93XyA000LmowvRrO"`
}
