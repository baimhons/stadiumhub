package internal

var ENV struct {
	Server      `mapstructure:"server"`
	Database    `mapstructure:"database"`
	Redis       `mapstructure:"redis"`
	APIFootball `mapstructure:"apifootball"`
	Stripe      `mapstructure:"stripe"`
	AdminData   `mapstructure:"adminData"`
	EmailKey    `mapstructure:"emailKey"`
	SecretKey   `mapstructure:"secretKey"`
}

type Server struct {
	Port int `mapstructure:"port" defaultValue:"8080"`
}

type Database struct {
	Host     string `mapstructure:"host" defaultValue:"dpg-d3ugk38dl3ps73f50fh0-a.singapore-postgres.render.com"`
	Port     int    `mapstructure:"port" defaultValue:"5432"`
	User     string `mapstructure:"user" defaultValue:"stadiumuser"`
	Password string `mapstructure:"password" defaultValue:"YpwWnlSfTBdnsJV9EtToStTrst37QI5M"`
	Name     string `mapstructure:"name" defaultValue:"stadiumhub"`
	Driver   string `mapstructure:"driver" defaultValue:"postgres"`
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

type AdminData struct {
	Username    string `mapstructure:"username" defaultValue:"adminStadiumHub"`
	Email       string `mapstructure:"email" defaultValue:"adminStadiumHub05@gmail.com"`
	Password    string `mapstructure:"password" defaultValue:"4Dm1n3_7-0"`
	PhoneNumber string `mapstructure:"phoneNumber" defaultValue:"012345678"`
	Role        string `mapstructure:"role" defaultValue:"admin"`
}

type SecretKey struct {
	SecretKey string `mapstructure:"secretKey" defaultValue:"s3cr3t_k3y_f0r_stad1um_hub7-0"`
}

type EmailKey struct {
	EmailKey string `mapstructure:"emailKey" defaultValue:"scfv mowy chgo tmjy"`
}
