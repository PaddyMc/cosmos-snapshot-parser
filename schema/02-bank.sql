/* ---- SUPPLY ---- */

CREATE TYPE COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE supply
(
    height BIGINT  NOT NULL,
    denom  TEXT,
    amount TEXT
);
