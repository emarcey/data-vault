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
    type TEXT REFERENCES admin.user_type(id)
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

COMMIT;