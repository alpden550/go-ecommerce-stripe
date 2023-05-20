package configs

type Config struct {
	Port int
	Env  string
	Api  string
	DB   struct {
		Dsn string
	}
	Stripe struct {
		Secret string
		Key    string
	}
}
