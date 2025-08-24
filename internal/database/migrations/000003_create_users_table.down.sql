-- +migrate Down
-- Os índices não precisam ser removidos explicitamente, pois eles são removidos com a tabela.
DROP TABLE IF EXISTS master.users;
