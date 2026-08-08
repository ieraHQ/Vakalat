-- Seed roles and permissions
INSERT INTO roles (id, name, description) VALUES
    ('11111111-1111-1111-1111-111111111111', 'admin', 'Administrator with full access'),
    ('22222222-2222-2222-2222-222222222222', 'advocate', 'Advocate with matter access'),
    ('33333333-3333-3333-3333-333333333333', 'staff', 'Staff with limited access');

INSERT INTO permissions (id, name, description) VALUES
    ('11111111-1111-1111-1111-111111111111', 'manage_users', 'Manage users and roles'),
    ('22222222-2222-2222-2222-222222222222', 'manage_matters', 'Manage matters and hearings'),
    ('33333333-3333-3333-3333-333333333333', 'manage_documents', 'Manage documents and OCR'),
    ('44444444-4444-4444-4444-444444444444', 'manage_clients', 'Manage clients and contacts'),
    ('55555555-5555-5555-5555-555555555555', 'manage_search', 'Use full-text and semantic search'),
    ('66666666-6666-6666-6666-666666666666', 'manage_ai', 'Use the AI assistant (summarize/ask/draft)');

INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111'),
    ('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222'),
    ('11111111-1111-1111-1111-111111111111', '33333333-3333-3333-3333-333333333333'),
    ('11111111-1111-1111-1111-111111111111', '44444444-4444-4444-4444-444444444444'),
    ('11111111-1111-1111-1111-111111111111', '55555555-5555-5555-5555-555555555555'),
    ('11111111-1111-1111-1111-111111111111', '66666666-6666-6666-6666-666666666666'),
    ('22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222'),
    ('22222222-2222-2222-2222-222222222222', '33333333-3333-3333-3333-333333333333'),
    ('22222222-2222-2222-2222-222222222222', '44444444-4444-4444-4444-444444444444'),
    ('22222222-2222-2222-2222-222222222222', '55555555-5555-5555-5555-555555555555'),
    ('22222222-2222-2222-2222-222222222222', '66666666-6666-6666-6666-666666666666');

-- Seed sample users.
-- Dev-only password for both accounts: "ChangeMe123!" — hashed with the app's
-- Argon2id + per-password salt scheme (backend/api/auth/password.go), NOT a
-- placeholder string. Rotate before any non-local use.
INSERT INTO users (id, name, email, password_hash, role_id) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Admin User', 'admin@vakalat.com', 'q6JGwV3dxUTkm+ZBZ+YzuA.O7YgcutFVlk24iXuXUydez8gGCqseX585vLiPC8kRA8', '11111111-1111-1111-1111-111111111111'),
    ('22222222-2222-2222-2222-222222222222', 'Advocate User', 'advocate@vakalat.com', 'q6JGwV3dxUTkm+ZBZ+YzuA.O7YgcutFVlk24iXuXUydez8gGCqseX585vLiPC8kRA8', '22222222-2222-2222-2222-222222222222');

-- Seed sample clients
INSERT INTO clients (id, name, type, email, phone) VALUES
    ('11111111-1111-1111-1111-111111111111', 'John Doe', 'individual', 'john@example.com', '+1234567890'),
    ('22222222-2222-2222-2222-222222222222', 'Acme Corp', 'organization', 'contact@acme.com', '+1987654321');