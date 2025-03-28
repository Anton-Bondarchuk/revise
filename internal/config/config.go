package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StorageConfig  StorageConfig     `yaml:"storage" env-required:"true"`
	OauthConfig 
}

type StorageConfig struct {
	Host     string `yml:"host"`
	Port     int    `yml:"port"`
	Username string `yml:"username"`
	Password string `yml:"password"`
	Database string `yml:"database"`
}

type OauthConfig struct {
	RedirectURI  string `yaml:"redirect_uri" env-default:"/auth/google/callback"`
	ClientID     string `yaml:"client_id" env-required:"true"`
	ClientSecret string `yaml:"client_secter" env-required:"true"`
	Scopes       []string `yaml:"scopes" env-default:"https://www.googleapis.com/auth/userinfo.email"`
	GoogleAuthURL string `yaml:"google_auth_url" env-default:"https://www.googleapis.com/oauth2/v2/userinfo?access_token="`
}


func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	// TODO: check the current path
	fmt.Println("config path: ", res)
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
