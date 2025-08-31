CREATE TABLE IF NOT EXISTS master.permissions (
    id UUID NOT NULL,
    company_global_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    PRIMARY KEY (id)
);

-- Index by company_global_id
CREATE INDEX IF NOT EXISTS idx_permissions_company_global_id ON master.permissions (company_global_id);
-- Index by deleted_at
CREATE INDEX IF NOT EXISTS idx_permissions_deleted_at ON master.permissions (deleted_at);



CREATE TABLE IF NOT EXISTS master.roles (
    id UUID NOT NULL,
    company_global_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    can_edit BOOLEAN NOT NULL DEFAULT TRUE,
    can_delete BOOLEAN NOT NULL DEFAULT TRUE,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (id),
    CONSTRAINT unique_name UNIQUE (company_global_id, name)
);



-- Index by company_global_id
CREATE INDEX IF NOT EXISTS idx_roles_company_global_id ON master.roles (company_global_id);
-- Index by deleted_at
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON master.roles (deleted_at);

CREATE TABLE IF NOT EXISTS master.roles_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    -- A chave primária é composta para garantir que uma role não possa ter a mesma permissão duas vezes.
    CONSTRAINT pk_role_permissions PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role_permissions_role
        FOREIGN KEY (role_id)
        REFERENCES master.roles(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY (permission_id)
        REFERENCES master.permissions(id)
        ON DELETE CASCADE
);