-- Seed roles and permissions
INSERT INTO roles (id, name, description) VALUES
    ('11111111-1111-1111-1111-111111111111', 'admin', 'Administrator with full access'),
    ('22222222-2222-2222-2222-222222222222', 'advocate', 'Advocate with matter access'),
    ('33333333-3333-3333-3333-333333333333', 'staff', 'Staff with limited access');

INSERT INTO permissions (id, name, description) VALUES
    ('11111111-1111-1111-1111-111111111111', 'manage_users', 'Manage users and roles'),
    ('22222222-2222-2222-2222-222222222222', 'manage_matters', 'Manage matters and hearings'),
    ('33333333-3333-3333-3333-333333333333', 'manage_documents', 'Manage documents and OCR');

INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
    ('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222'),
    ('11111111-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333333'),
    ('22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222'),
    ('22222222-2222-2222-2222-222222222222', '33333333-3333-3333-3333-333333333333');

-- Seed sample users
INSERT INTO users (id, name, email, password_hash, role_id) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Admin User', 'admin@vakalat.com', '$2a$10$examplehash', '11111111-1111-1111-1111-111111111111'),
    ('22222222-2222-2222-2222-222222222222', 'Advocate User', 'advocate@vakalat.com', '$2a$10$examplehash', '22222222-2222-2222-2222-222222222222');

-- Seed sample clients
INSERT INTO clients (id, name, type, email, phone) VALUES
    ('11111111-1111-1111-1111-111111111111', 'John Doe', 'individual', 'john@example.com', '+1234567890'),
    ('22222222-2222-2222-2222-222222222222', 'Acme Corp', 'organization', 'contact@acme.com', '+1987654321');