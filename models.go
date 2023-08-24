package main

type EnvironmentConfig struct {
	FileName      string `env:"LOG_FILE" envDefault:"reterm.log"`
	MaxSize       int    `env:"LOG_MAX_SIZE" envDefault:"10"`
	MaxBackups    int    `env:"LOG_MAX_BACKUPS" envDefault:"5"`
	MaxAge        int    `env:"LOG_MAX_AGE" envDefault:"30"`
	RedisAddress  string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}
