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

CREATE TABLE IF NOT EXISTS outputs (
    txid      BYTEA NOT NULL,
    idx       INT NOT NULL, 
    recipient TEXT NOT NULL,
    amount    INT NOT NULL,
    PRIMARY KEY (txid, idx)
);

CREATE TABLE IF NOT EXISTS transactions (
    txid      BYTEA NOT NULL PRIMARY KEY,
    blockhash   BYTEA 
);

CREATE TABLE IF NOT EXISTS inputs (
    txid      BYTEA NOT NULL,  
    idx       INT NOT NULL,   
    prev_txid BYTEA NOT NULL, 
    prev_idx  INT NOT NULL,
    owner TEXT NOT NULL,
    PRIMARY KEY (txid, idx),
    FOREIGN KEY (prev_txid, prev_idx) REFERENCES outputs(txid, idx)
);