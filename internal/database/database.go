package database

import (
	"go-sales/internal/config"
	"go-sales/internal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error

	// Com o search_path na DSN, não precisamos mais de NamingStrategy ou callbacks para o schema.
	// O GORM usará a configuração padrão, que é o que queremos.
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})

	if err != nil {
		log.Error().Err(err).Msg("Fail to connect to database")
		return err
	}
	return nil
}

func Migrate(db *gorm.DB, cfg *config.Config) error {
	// Apenas executa a operação de apagar tabelas se a flag estiver explicitamente ativada.
	// É uma camada extra de segurança para evitar acidentes em produção.
	if cfg.DBRecreate && cfg.AppEnv != "production" {
		log.Warn().Msg("DB_RECREATE is true. Dropping all tables...")
		// Dropa as tabelas na ordem inversa para respeitar as chaves estrangeiras.
		err := db.Migrator().DropTable(model.RegisteredModels...)
		if err != nil {
			log.Error().Err(err).Msg("Fail to drop tables")
			return err
		}
		log.Info().Msg("Tables dropped successfully.")
	}

	// O GORM agora criará as tabelas no schema definido pelo search_path na DSN.
	log.Info().Msg("Running database migrations...")
	err := db.AutoMigrate(model.RegisteredModels...)
	if err != nil {
		log.Error().Err(err).Msg("Fail to migrate database")
		return err
	}
	return nil
}
