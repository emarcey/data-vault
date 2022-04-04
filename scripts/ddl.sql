-- This DDL is for Postgres
BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS admin;
COMMENT ON SCHEMA admin IS 'Admin schema contains tables used for administration of service (e.g. users, permissions)';
CREATE SCHEMA IF NOT EXISTS data;
COMMENT ON SCHEMA data IS 'Data schema contains tables used for user-defined data storage';

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    returns trigger AS $$
BEGIN
    NEW.updated_at = NOW();
    return NEW;
END;
$$ LANGUAGE PLPGSQL;

CREATE TABLE admin.user_type (
    id TEXT PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

INSERT INTO admin.user_type VALUES ('admin');
INSERT INTO admin.user_type VALUES ('developer');

CREATE TABLE admin.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    name TEXT NOT NULL,
    client_secret_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    type TEXT REFERENCES admin.user_type(id),
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL
);

COMMENT ON TABLE admin.users IS 'Users stores information about each user, including their user_id & a hash of the secret used to generate an access token.';
COMMENT ON COLUMN admin.users.client_secret_hash IS 'A hash of the unique client secret generated for this user. Of the form "{encryptionMethod}|||{encryptedValue}"';

CREATE UNIQUE INDEX uq__admin__users__name ON admin.users(name) WHERE is_active;

CREATE TRIGGER set_admin__users_timestamp
    BEFORE UPDATE ON admin.users
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE admin.access_tokens (
    id_hash TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    user_id UUID REFERENCES admin.users(id) NOT NULL,
    is_latest BOOLEAN NOT NULL,
    invalid_at TIMESTAMPTZ NOT NULL
);

COMMENT ON TABLE admin.access_tokens IS 'Access tokens stores generated access tokens for temporary usage. A user may only have one valid access token at a time.';
COMMENT ON COLUMN admin.access_tokens.is_latest IS 'Keeps track of which access token is most recent.';
COMMENT ON COLUMN admin.access_tokens.invalid_at IS 'Datetime at which access token will no longer be usable.';

CREATE INDEX idx__admin__access_tokens__user_is_valid ON admin.access_tokens(user_id, is_latest);
CREATE INDEX idx__admin__access_tokens__invalid_at ON admin.access_tokens(invalid_at);
CREATE UNIQUE INDEX uq__admin__access_tokens__user_valid ON admin.access_tokens(user_id) WHERE is_latest;

CREATE TRIGGER set_admin__access_tokens_timestamp
    BEFORE UPDATE ON admin.access_tokens
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


CREATE TABLE admin.secrets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TRIGGER set_admin__secrets_timestamp
    BEFORE UPDATE ON admin.secrets
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMENT ON TABLE admin.secrets IS 'secrets stores all user created secrets for data being stored. Kept separate from information schema so we can log who did what.';
CREATE UNIQUE INDEX uq__admin__secrets__name ON admin.secrets(name) WHERE is_active;

CREATE TABLE admin.secret_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES admin.users(id) NOT NULL,
    secret_id UUID REFERENCES admin.secrets(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TRIGGER set_admin__secret_permissions_timestamp
    BEFORE UPDATE ON admin.secret_permissions
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMENT ON TABLE admin.secret_permissions IS 'secret permissions stores all secret access permissions';
CREATE UNIQUE INDEX uq__admin__secret_permissions__user_secret ON admin.secret_permissions(user_id, secret_id) WHERE is_active;

CREATE TABLE admin.user_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TRIGGER set_admin__user_groups_timestamp
    BEFORE UPDATE ON admin.user_groups
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMENT ON TABLE admin.user_groups IS 'user_groups stores all distinct user permission groups';
CREATE UNIQUE INDEX uq__admin__user_groups__name ON admin.user_groups(name) WHERE is_active;


CREATE TABLE admin.user_group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES admin.users(id) NOT NULL,
    user_group_id UUID REFERENCES admin.user_groups(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TRIGGER set_admin__user_group_members_timestamp
    BEFORE UPDATE ON admin.user_group_members
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMENT ON TABLE admin.user_group_members IS 'user_group_members stores the mapping of users to user groups';
CREATE UNIQUE INDEX uq__admin__user_group_members__user_secret ON admin.user_group_members(user_id, user_group_id) WHERE is_active;

CREATE TABLE admin.secret_group_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_group_id UUID REFERENCES admin.user_groups(id) NOT NULL,
    secret_id UUID REFERENCES admin.secrets(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TRIGGER set_admin__secret_group_permissions_timestamp
    BEFORE UPDATE ON admin.secret_group_permissions
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

COMMENT ON TABLE admin.secret_group_permissions IS 'secret group permissions stores all secret access permissions for user groups';
CREATE UNIQUE INDEX uq__admin__secret_group_permissions__user_group_secret ON admin.secret_group_permissions(user_group_id, secret_id) WHERE is_active;

COMMIT;
