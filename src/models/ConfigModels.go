package models

type AppConfiguration struct {
	Database  DatabaseConfiguration
	SecretKey SecretKeys
}

type DatabaseConfiguration struct {
	Name     string `mapstructure:"DB_NAME"`
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}

type SecretKeys struct {
	HashKey      string `mapstructure:"HASH_KEY"`
	SharedSecret string `mapstructure:"SHARED_SECRET"`
	SignatureMsg string `mapstructure:"SIGNATURE_MSG"`
	RSAKeyFile   string `mapstructure:"RSA_KEY_FILE"`
}
