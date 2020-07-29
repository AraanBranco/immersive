package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Configuration wrapper
type Configuration struct {
	v   *viper.Viper
	env string
}

var conf *Configuration

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	fmt.Println("[Env]: ", env)
	v := viper.New()
	v.SetConfigName("config")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(env)
	v.AddConfigPath(os.Getenv("CONFIG_DIR"))
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Println(context.TODO(), " Ocorreu um problema ao ler o arquivo de configuração:", err)
	}

	conf = &Configuration{v, env}
}

// Get returns configuration wrapper
func Get() *Configuration {
	return conf
}

// GetEnvConfString returns a value and key on configuration
func (config *Configuration) GetEnvConfString(key string) string {
	return config.v.GetString(key)
}

// GetEnvConfBool returns a value of and key on configuration as boolean
func (config *Configuration) GetEnvConfBool(key string) bool {
	return config.v.GetBool(key)
}

// GetString return a value of a key as String
func (config *Configuration) GetString(key string) string {
	return config.v.GetString(key)
}

func (config *Configuration) GetEnvConfStringSlice(key string) []string {
	return config.v.GetStringSlice(key)
}

func (config *Configuration) GetEnvConfInteger(key string) int {
	return config.v.GetInt(key)
}

func (config *Configuration) GetInteger(key string) int {
	return config.v.GetInt(key)
}
