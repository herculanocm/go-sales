INSERT INTO master.roles (id, name) VALUES
    (gen_random_uuid(), 'ROLE_ADMIN'),
    (gen_random_uuid(), 'ROLE_USER'),
    (gen_random_uuid(), 'ROLE_GUEST')
ON CONFLICT (name) DO NOTHING;