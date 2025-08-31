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

    CONSTRAINT pk_company_globals PRIMARY KEY (id),
    CONSTRAINT uk_company_globals_cgc UNIQUE (cgc)
);


CREATE INDEX idx_company_globals_deleted_at ON master.company_globals (deleted_at);