-- +migrate Up

-- Garante que o schema 'master' exista antes de tentar criar a tabela nele.
-- 'IF NOT EXISTS' torna o script seguro para ser executado várias vezes.
CREATE SCHEMA IF NOT EXISTS master;

-- Cria a tabela dentro do schema 'master'.
CREATE TABLE IF NOT EXISTS master.company_globals (
    id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    cgc VARCHAR(14) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT pk_company_globals PRIMARY KEY (id)
);

-- Os índices também devem ser criados no schema correto.
CREATE UNIQUE INDEX uk_company_globals_cgc ON master.company_globals (cgc);
CREATE INDEX idx_company_globals_deleted_at ON master.company_globals (deleted_at);