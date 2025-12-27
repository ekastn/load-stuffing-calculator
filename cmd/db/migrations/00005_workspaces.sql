-- +goose Up
-- +goose StatementBegin

-- Workspaces
CREATE TABLE workspaces (
    workspace_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type TEXT NOT NULL CHECK (type IN ('personal', 'organization')),
    name TEXT NOT NULL,
    owner_user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_workspaces_owner ON workspaces(owner_user_id);

-- Workspace memberships
CREATE TABLE members (
    member_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(workspace_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (workspace_id, user_id)
);

CREATE INDEX idx_members_workspace ON members(workspace_id);
CREATE INDEX idx_members_user ON members(user_id);

-- Workspace invites
CREATE TABLE invites (
    invite_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(workspace_id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE RESTRICT,
    token_hash TEXT NOT NULL,
    invited_by_user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
    expires_at TIMESTAMPTZ NOT NULL,
    accepted_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_invites_workspace ON invites(workspace_id);
CREATE INDEX idx_invites_email ON invites(email);

-- Platform members (e.g. founder)
CREATE TABLE platform_members (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Active workspace tracking for refresh tokens
ALTER TABLE refresh_tokens
    ADD COLUMN workspace_id UUID REFERENCES workspaces(workspace_id) ON DELETE SET NULL;

CREATE INDEX idx_refresh_tokens_workspace ON refresh_tokens(workspace_id);

-- Workspace scoping for master data and plans.
-- Note: keep workspace_id nullable for trial/guest plans.
ALTER TABLE containers
    ADD COLUMN workspace_id UUID REFERENCES workspaces(workspace_id) ON DELETE CASCADE;

ALTER TABLE products
    ADD COLUMN workspace_id UUID REFERENCES workspaces(workspace_id) ON DELETE CASCADE;

ALTER TABLE load_plans
    ADD COLUMN workspace_id UUID REFERENCES workspaces(workspace_id) ON DELETE CASCADE;

CREATE INDEX idx_containers_workspace ON containers(workspace_id);
CREATE INDEX idx_products_workspace ON products(workspace_id);
CREATE INDEX idx_load_plans_workspace ON load_plans(workspace_id);

-- Backfill legacy resources: route old records into a single synthetic workspace.
-- This keeps the database consistent for existing installs.
DO $$
DECLARE
    legacy_owner UUID;
    legacy_workspace UUID;
    legacy_admin_role UUID;
BEGIN
    SELECT user_id INTO legacy_owner
    FROM users
    ORDER BY created_at ASC
    LIMIT 1;

    IF legacy_owner IS NOT NULL THEN
        INSERT INTO workspaces (type, name, owner_user_id)
        VALUES ('organization', 'Legacy Workspace', legacy_owner)
        RETURNING workspace_id INTO legacy_workspace;

        UPDATE containers SET workspace_id = legacy_workspace WHERE workspace_id IS NULL;
        UPDATE products SET workspace_id = legacy_workspace WHERE workspace_id IS NULL;
        UPDATE load_plans SET workspace_id = legacy_workspace WHERE workspace_id IS NULL AND created_by_type <> 'guest';

        -- Best-effort: if roles were seeded already, grant access to the first user.
        SELECT role_id INTO legacy_admin_role FROM roles WHERE name = 'admin' LIMIT 1;
        IF legacy_admin_role IS NOT NULL THEN
            INSERT INTO members (workspace_id, user_id, role_id)
            VALUES (legacy_workspace, legacy_owner, legacy_admin_role)
            ON CONFLICT DO NOTHING;
        END IF;
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE load_plans DROP COLUMN IF EXISTS workspace_id;
ALTER TABLE products DROP COLUMN IF EXISTS workspace_id;
ALTER TABLE containers DROP COLUMN IF EXISTS workspace_id;

ALTER TABLE refresh_tokens DROP COLUMN IF EXISTS workspace_id;

DROP TABLE IF EXISTS platform_members;
DROP TABLE IF EXISTS invites;
DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS workspaces;

-- +goose StatementEnd
