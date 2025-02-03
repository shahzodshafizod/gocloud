CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS agents (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4()
    , name          VARCHAR(255)    NOT NULL
    , priority      INTEGER         NOT NULL
    , created_at    TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP
    , updated_at    TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO agents (id, name, priority)
VALUES ('87b77778-6a51-4ef7-a9cd-e2eec44aefaf', 'orders', 1);
