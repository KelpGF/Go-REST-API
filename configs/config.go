package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type ConfigType struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	WebServerHost string `mapstructure:"WEB_SERVER_HOST"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn  int    `mapstructure:"JWT_EXPIRES_IN"`
	TokenAuth     *jwtauth.JWTAuth
}

func LoadConfig(path string) *ConfigType {
	var configs *ConfigType

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // environment variables are prioritized instead of .env file

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&configs)
	if err != nil {
		panic(err)
	}

	if configs.WebServerPort == "" {
		configs.WebServerPort = "3000"
	}

	configs.TokenAuth = jwtauth.New("HS256", []byte(configs.JWTSecret), nil)

	return configs
}
