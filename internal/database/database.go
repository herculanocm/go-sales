package database

import (
	"fmt"
	"go-sales/internal/config"
	"go-sales/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// schemaAwareNamer é uma estratégia de nomenclatura que adiciona um schema
// e respeita o método TableName() do modelo.
type schemaAwareNamer struct {
	schema.NamingStrategy
	schemaPrefix string
}

// TableName é a nossa implementação personalizada.
func (n schemaAwareNamer) TableName(table string) string {
	// Esta implementação padrão não é mais o que queremos.
	// O GORM chama um método diferente durante a migração (Migrator().CreateTable())
	// que nos permite acessar o modelo real.
	// Deixamos isso como um fallback.
	return fmt.Sprintf("%s.%s", n.schemaPrefix, n.NamingStrategy.TableName(table))
}

// Migrator é o método chave. Ele nos dá acesso ao modelo (`dst`).
func (n schemaAwareNamer) Migrator(db *gorm.DB) gorm.Migrator {
	// Criamos uma função para ser chamada antes de criar a tabela.
	db.Callback().Create().Before("gorm:create_table").
		Register("custom:table_name", func(d *gorm.DB) {
			if d.Statement.Schema != nil {
				// Verifica se o modelo (d.Statement.Model) implementa nossa interface model.Tabler
				if tabler, ok := d.Statement.Model.(model.Tabler); ok {
					// Se sim, usa o nome da tabela do modelo.
					d.Statement.Table = tabler.TableName()
				}
				// Adiciona o prefixo do schema ao nome da tabela (seja o padrão ou o personalizado).
				d.Statement.Table = fmt.Sprintf("%s.%s", n.schemaPrefix, d.Statement.Table)
			}
		})
	return db.Migrator()
}

func Connect(cfg *config.Config) error {
	log.Println("Connecting to database...")
	var err error

	namingStrategy := schemaAwareNamer{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		schemaPrefix: cfg.DBSchema,
	}

	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		NamingStrategy: namingStrategy,
	})

	if err != nil {
		log.Fatal("Fail to connect to database:", err)
		return err
	}

	log.Println("Connected to database successfully!")

	return nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Migrating database...")

	err := db.AutoMigrate(model.RegisteredModels...)
	if err != nil {
		log.Fatal("Fail to migrate database:", err)
		return err
	}

	log.Println("Database migrated successfully!")
	return nil
}
