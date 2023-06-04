package configs

type Config struct {
	Port int
	Env  string
	DB   struct {
		Dsn string
	}
	Stripe struct {
		Secret string
		Key    string
	}
	SMTP struct {
		Host      string
		Port      int
		Username  string
		Password  string
		EmailFrom string
	}
	SecretKey string
	Api       string
	FrontEnd  string
}
