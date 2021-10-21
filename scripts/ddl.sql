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

COMMENT ON TABLE admin.users IS 'Users stores information about each user, including their client_id & a hash of the secret used to generate an access token.';
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


CREATE TABLE admin.data_tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL,
    is_active BOOLEAN NOT NULL
);

COMMENT ON TABLE admin.data_tables IS 'Data tables stores all user created tables for data being stored. Kept separate from information schema so we can log who did what.';
CREATE UNIQUE INDEX uq__admin__data_tables__name ON admin.data_tables(name) WHERE is_active;

CREATE TRIGGER set_admin__data_tables_timestamp
    BEFORE UPDATE ON admin.data_tables
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE admin.data_table_permissions (
    user_id UUID REFERENCES admin.users(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    table_id UUID NOT NULL REFERENCES admin.data_tables(id),
    is_decrypt_allowed BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES admin.users(id) NOT NULL,
    updated_by UUID REFERENCES admin.users(id) NOT NULL,
    is_active BOOLEAN NOT NULL,
    PRIMARY KEY (user_id, table_id)
);

COMMENT ON TABLE admin.data_table_permissions IS 'Table permissions grants access for a user to query a given table. Access should only be granted by an admin';
COMMENT ON COLUMN admin.data_table_permissions.is_decrypt_allowed IS 'If true, user is allowed to decrypt results of query. If false, user can only get metadata';

CREATE TRIGGER set_admin__data_table_permissions_timestamp
    BEFORE UPDATE ON admin.data_table_permissions
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


CREATE TABLE admin.encrypted_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_id UUID NOT NULL REFERENCES admin.data_tables(id),
    row_id UUID NOT NULL,
    column_name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    hash_value TEXT NOT NULL
);

CREATE INDEX idx__admin__encrypted_keys__hash_value ON admin.encrypted_keys(hash_value);
CREATE UNIQUE INDEX uq__admin__encrypted_keys__table_row_column ON admin.encrypted_keys(table_id, row_id, column_name);

CREATE TABLE admin.encrypted_key_metadata_type (
    id TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TABLE admin.encrypted_key_metadata (
    encrypted_key_id UUID REFERENCES admin.encrypted_keys(id) NOT NULL,
    encrypted_key_metadata_type_id TEXT REFERENCES admin.encrypted_key_metadata_type(id) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY (encrypted_key_id, encrypted_key_metadata_type_id)
);


--raw secret: 9e5d5da6-24ed-42a9-a105-c531bed8175d
INSERT INTO admin.users (id, name, client_secret_hash, type) VALUES ('03b6f72c-f3f4-43d9-a705-17b326924d74', 'admin', 'sha256:9dbcc64bd8006081f719719f9be420059b907429710bc64215c843596c6d0d64', 'admin');

COMMIT;