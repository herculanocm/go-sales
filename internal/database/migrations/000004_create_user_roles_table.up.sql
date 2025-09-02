CREATE TABLE master.user_roles (
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    -- A chave primária é composta para garantir que um usuário não possa ter a mesma role duas vezes.
    CONSTRAINT pk_user_roles PRIMARY KEY (user_id, role_id),

    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
        REFERENCES master.users(id)
        ON DELETE CASCADE, -- Se um usuário for deletado, suas associações de role também são.

    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
        REFERENCES master.roles(id)
        ON DELETE CASCADE -- Se uma role for deletada, suas associações também são.
);