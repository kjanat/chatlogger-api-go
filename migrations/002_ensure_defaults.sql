-- Add a superadmin user for initial setup (password: admin1234, and user1234)
-- This is just a placeholder; in production, you would create users via API
INSERT INTO organizations
    (name, slug)
VALUES
    ('System Admin', 'system-admin'),
    ('Not assigned to any org', 'unassigned');

INSERT INTO users
    (email, password_hash, role, organization_id, first_name, last_name)
VALUES
    ('admin@example.com', '$2a$10$w04EwYhTl/aFFubCHTWGDu94gxydNHpYmGr/IdSAtFDExIop2Zwfm', 'superadmin', 1, 'System', 'Admin'),
    ('user@example.com', '$2a$10$ACqSWyOSs51loUHoCCFSbeMAOErHwjpaXnv6NDNsvYJb19ZBLZQQm', 'user', 2, 'Default', 'User');
