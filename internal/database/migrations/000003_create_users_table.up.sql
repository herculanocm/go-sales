-- Cria a tabela dentro do schema 'master'.
CREATE TABLE IF NOT EXISTS master.users (
    id UUID NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,
    company_global_id UUID NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT pk_users PRIMARY KEY (id),

    -- A referência da chave estrangeira também deve incluir o schema.
    CONSTRAINT fk_users_company_globals
        FOREIGN KEY (company_global_id)
        REFERENCES master.company_globals(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

-- Os índices também devem ser criados no schema correto.
CREATE UNIQUE INDEX uk_users_email_company_global ON master.users (company_global_id, email_address);
CREATE INDEX idx_users_deleted_at ON master.users (deleted_at);