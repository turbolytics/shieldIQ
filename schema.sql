CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tenants
CREATE TABLE tenants
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ      DEFAULT now()
);

-- Webhooks (stream source)
CREATE TABLE webhooks
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id  UUID NOT NULL REFERENCES tenants (id),
    name       TEXT NOT NULL,
    secret     TEXT NOT NULL,
    source     TEXT NOT NULL, -- e.g., 'github'
    events     JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ      DEFAULT now()
);

-- Notification Channels (e.g., Slack)
CREATE TABLE notification_channels
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id  UUID  NOT NULL REFERENCES tenants (id),
    name       TEXT  NOT NULL,
    type       TEXT  NOT NULL, -- 'slack', 'webhook'
    config     JSONB NOT NULL,
    created_at TIMESTAMPTZ      DEFAULT now()
);

-- Rules
CREATE TABLE rules
(
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id               UUID NOT NULL REFERENCES tenants (id),
    name                    TEXT NOT NULL,
    sql                     TEXT NOT NULL,
    severity                TEXT NOT NULL,
    notification_channel_id UUID NOT NULL REFERENCES notification_channels (id),
    enabled                 BOOLEAN          DEFAULT true,
    created_at              TIMESTAMPTZ      DEFAULT now()
);

-- Alerts
CREATE TABLE alerts
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id  UUID  NOT NULL REFERENCES tenants (id),
    rule_id    UUID  NOT NULL REFERENCES rules (id),
    message    TEXT  NOT NULL,
    severity   TEXT  NOT NULL,
    data       JSONB NOT NULL,
    created_at TIMESTAMPTZ      DEFAULT now()
);

insert into tenants (id, name) VALUES ('00000000-0000-0000-0000-000000000000', 'root');