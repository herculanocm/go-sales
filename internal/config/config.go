package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config armazena todas as configurações da aplicação.
// As tags 'mapstructure' são usadas para o Viper mapear os valores para a struct.
type Config struct {
	AppEnv       string `mapstructure:"APP_ENV"`
	AppAPIPrefix string `mapstructure:"APP_API_PREFIX"`
	AppAPIPort   string `mapstructure:"APP_API_PORT"`

	DBSchema string `mapstructure:"DEFAULT_SCHEMA"`
	DBHost   string `mapstructure:"DB_HOST"`
	DBUser   string `mapstructure:"DB_USER"`
	DBPass   string `mapstructure:"DB_PASS"`
	DBName   string `mapstructure:"DB_NAME"`
	DBPort   string `mapstructure:"DB_PORT"`
}

// LoadConfig lê a configuração e retorna uma instância de Config.
// Não há mais variável global.
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Ignora o erro se o arquivo não for encontrado, pois as variáveis de ambiente podem ser usadas.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("erro ao ler arquivo de config: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da config: %w", err)
	}

	return &cfg, nil
}

// DSN retorna a string de conexão com o banco de dados (Data Source Name).
func (c *Config) DSN() string {
	// Adicionar "search_path" é a forma mais confiável de definir o schema padrão.
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.DBSchema,
	)
}

func (c *Config) String() string {
	return fmt.Sprintf("AppEnv: %s, AppAPIPrefix: %s, AppAPIPort: %s, DBSchema: %s, DBHost: %s, DBUser: %s, DBName: %s, DBPort: %s",
		c.AppEnv, c.AppAPIPrefix, c.AppAPIPort, c.DBSchema, c.DBHost, c.DBUser, c.DBName, c.DBPort,
	)
}
