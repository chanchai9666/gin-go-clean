package configs

type Config struct {
	Port       int    `env:"PORT" envDefault:"8080"`
	AppEnv     string `env:"APP_ENV" envDefault:"local"`
	DbHost     string `env:"DB_HOST" envDefault:"localhost"`
	DbPort     string `env:"DB_PORT" envDefault:"5432"`
	DbDatabase string `env:"DB_DATABASE" envDefault:"arczedxxxx"`
	DbUsername string `env:"DB_USERNAME" envDefault:"admin"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"1234"`
}
