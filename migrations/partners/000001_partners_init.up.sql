CREATE TABLE IF NOT EXISTS partners (
    id                  SERIAL          PRIMARY KEY
    , title             VARCHAR(50)     NOT NULL
    , brand             VARCHAR(20)     NOT NULL
    , phone             VARCHAR(13)     NOT NULL
    , email             VARCHAR(20)     NOT NULL
    , address           VARCHAR(100)    NOT NULL
    , api_url           VARCHAR(100)    NOT NULL
    , verified          BOOLEAN         NOT NULL DEFAULT FALSE
    , enabled           BOOLEAN         NOT NULL DEFAULT FALSE
    , created_at        TIMESTAMPTZ     NOT NULL DEFAULT now()
    , updated_at        TIMESTAMPTZ     NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS products (
    id              SERIAL          PRIMARY KEY
    , title         VARCHAR(100)    NOT NULL
    , description   VARCHAR(256)    NOT NULL DEFAULT ''
    , picture_url   VARCHAR(50)     NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS available (
    product_id      INT         NOT NULL REFERENCES products (id)
    , partner_id    INT         NOT NULL REFERENCES partners (id)
    , price         BIGINT      NOT NULL
    , active        BOOLEAN     NOT NULL DEFAULT TRUE
    , UNIQUE (product_id, partner_id)
);

INSERT INTO partners (
    title
    , brand
    , phone
    , email
    , address
    , api_url
    , verified
    , enabled
) VALUES (
    'Partner 1 Title'
    , 'Partner 1 Brand'
    , '(123)456-7890'
    , 'test@mail.tst'
    , 'Partner 1 Address'
    , 'Partner 1 API URL'
    , TRUE
    , TRUE
), (
    'Partner 2 Title'
    , 'Partner 2 Brand'
    , '(234)567-8901'
    , 'test@mail.tst'
    , 'Partner 2 Address'
    , 'Partner 2 API URL'
    , TRUE
    , TRUE
);

INSERT INTO products (
    title
    , description
    , picture_url
) VALUES (
    'Product 1 Title'
    , 'Product 1 Description'
    , 'Product 1 Picture URL'
), (
    'Product 2 Title'
    , 'Product 2 Description'
    , 'Product 2 Picture URL'
), (
    'Product 3 Title'
    , 'Product 3 Description'
    , 'Product 3 Picture URL'
), (
    'Product 4 Title'
    , 'Product 4 Description'
    , 'Product 4 Picture URL'
);

INSERT INTO available (
    product_id
    , partner_id
    , price
) VALUES (
    1, 1, 100
), (
    1, 2, 120
), (
    3, 1, 1500
), (
    3, 2, 1200
), (
    4, 2, 500
), (
    4, 1, 500
), (
    2, 1, 350
), (
    2, 2, 350
);
