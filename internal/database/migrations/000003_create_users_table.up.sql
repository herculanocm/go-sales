-- Cria a tabela dentro do schema 'master'.

CREATE TABLE IF NOT EXISTS master.users (
    id BIGINT NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    email_recovery VARCHAR(255),
    email_verified BOOLEAN NOT NULL DEFAULT false,
    email_verified_at TIMESTAMPTZ,
    phone_number VARCHAR(20),
    phone_verified BOOLEAN NOT NULL DEFAULT false,
    phone_verified_at TIMESTAMPTZ,
    password_hash VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,
    actived BOOLEAN NOT NULL DEFAULT false,
    activation_key VARCHAR(255),
    activated_at TIMESTAMPTZ,
    reset_key VARCHAR(255),
    reset_requested TIMESTAMPTZ,
    reset_at TIMESTAMPTZ,
    company_global_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT uk_users_email_company_global UNIQUE (company_global_id, email_address),

    CONSTRAINT fk_users_company_globals
        FOREIGN KEY (company_global_id)
        REFERENCES master.company_globals(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

-- Índices para buscas frequentes
CREATE INDEX idx_users_company_global_id ON master.users (company_global_id);
CREATE INDEX idx_users_email_address ON master.users (email_address);
CREATE INDEX idx_users_enabled ON master.users (enabled);
CREATE INDEX idx_users_company_global_id_enabled_email ON master.users (company_global_id, enabled, email_address);
CREATE INDEX idx_users_deleted_at ON master.users (deleted_at);