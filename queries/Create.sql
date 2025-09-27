CREATE TABLE IF NOT EXISTS blocks (
    hash        BYTEA PRIMARY KEY,      
    prev_hash   BYTEA,                 
    data        BYTEA NOT NULL,             
    timestamp   BIGINT NOT NULL,           
    nonce       INTEGER NOT NULL,          
    CONSTRAINT fk_prev FOREIGN KEY (prev_hash) REFERENCES blocks(hash)
);

CREATE TABLE IF NOT EXISTS tail (
    hash BYTEA PRIMARY KEY
);

