CREATE TABLE IF NOT EXISTS banks (
    id                  VARCHAR(10) NOT NULL PRIMARY KEY
    , title             VARCHAR(50) NOT NULL
    , webcheckout_url   VARCHAR(50) NOT NULL DEFAULT ''
    , active            BOOLEAN     NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS orders (
    id                      BIGSERIAL       PRIMARY KEY
    , order_id              VARCHAR(50)     NOT NULL
    , customer_id           VARCHAR(40)     NOT NULL
    , customer_name         VARCHAR(100)    NOT NULL
    , customer_phone        VARCHAR(13)     NOT NULL
    , customer_notif_token  VARCHAR(100)    NOT NULL
    , delivery_address      VARCHAR(100)    NOT NULL
    , partner_id            INTEGER         NOT NULL
    , partner_title         VARCHAR(50)     NOT NULL
    , partner_brand         VARCHAR(50)     NOT NULL
    , pickup_address        VARCHAR(100)    NOT NULL DEFAULT ''
    , deliverer_id          VARCHAR(40)     NOT NULL DEFAULT ''
    , total_amount          BIGINT          NOT NULL
    , paid_amount           BIGINT          NOT NULL DEFAULT 0
    , paytype               VARCHAR(10)     NOT NULL
    , products              JSONB           NOT NULL
    , status                VARCHAR(10)     NOT NULL DEFAULT 'pending'
    -- , deadline           TIMESTAMPTZ     NOT NULL
    , created_at            TIMESTAMPTZ     NOT NULL DEFAULT now()
    , updated_at            TIMESTAMPTZ     NOT NULL DEFAULT now()
);

INSERT INTO banks (
    id, title, webcheckout_url
) VALUES (
    'visa', 'Visa Inc.', 'Visa Webcheckout URL'
), (
    'km', 'Tajik National Payment System: ""Korti Milli""', 'KM Webcheckout URL'
), (
    'paypal', 'PayPal Holdings, Inc.', 'PayPal Webcheckout URL'
);
