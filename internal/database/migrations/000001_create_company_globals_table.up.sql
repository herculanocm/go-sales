-- Cria a tabela dentro do schema 'master'.
CREATE TABLE IF NOT EXISTS master.company_globals (
    id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    social_name VARCHAR(255),
    description TEXT,
    cgc VARCHAR(14) NOT NULL,
    email VARCHAR(150),
    enabled BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT pk_company_globals PRIMARY KEY (id),
    CONSTRAINT uk_company_globals_cgc UNIQUE (cgc)
);

CREATE INDEX idx_company_globals_deleted_at ON master.company_globals (deleted_at);


-- Tabela de contatos globais da empresa
CREATE TABLE IF NOT EXISTS master.company_global_contacts (
    id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(150),
    phone VARCHAR(20),
    cgc VARCHAR(40),

    CONSTRAINT pk_company_global_contacts PRIMARY KEY (id),
    CONSTRAINT fk_company_global_contacts_company_id FOREIGN KEY (company_id)
        REFERENCES master.company_globals(id) ON DELETE CASCADE
);


-- Tabela de endereços globais da empresa
CREATE TABLE IF NOT EXISTS master.company_global_addresses (
    id BIGINT NOT NULL,
    company_id BIGINT NOT NULL,
    street VARCHAR(255) NOT NULL,
    street_number VARCHAR(50),
    street_complement VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,

    CONSTRAINT pk_company_global_addresses PRIMARY KEY (id),
    CONSTRAINT fk_company_global_addresses_company_id FOREIGN KEY (company_id)
        REFERENCES master.company_globals(id) ON DELETE CASCADE
);