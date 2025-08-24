# Inclui as variáveis do arquivo .env e as exporta para os comandos do shell.
include .env
export

# Define o caminho para as migrações para não repetir.
MIGRATE_PATH=internal/database/migrations

# .PHONY garante que o make execute o comando mesmo que exista um arquivo com o mesmo nome.
.PHONY: migrate-create migrate-up migrate-down migrate-reset

# Comando para criar um novo arquivo de migração.
# Uso: make migrate-create name=add_description_to_users
migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

# Comando para aplicar todas as migrações pendentes.
# Uso: make migrate-up
migrate-up:
	@migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) up

# Comando para reverter a última migração aplicada.
# Uso: make migrate-down
migrate-down:
	@migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) down 1

# Comando para dropar tudo e recriar do zero. CUIDADO: APAGA TODOS OS DADOS.
# Uso: make migrate-reset
migrate-reset:
	@echo "⚠️  ATENÇÃO: Este comando irá apagar TODOS os dados do banco de dados '$(DB_NAME)'."
	@echo "Pressione ENTER para continuar ou CTRL+C para cancelar."
	@read dummy
	@echo "-> Dropando o banco de dados..."
	@migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) drop -f
	@echo "-> Recriando o banco de dados a partir das migrações..."
	@migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) up
	@echo "✅ Banco de dados resetado com sucesso."