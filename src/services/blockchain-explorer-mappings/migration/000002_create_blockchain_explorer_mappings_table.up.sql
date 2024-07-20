CREATE TABLE IF NOT EXISTS blockchain_explorer_mappings
(
    blockchain_id    bigserial NOT NULL,
    blockchain_explorer_id      bigserial NOT NULL,
    PRIMARY KEY (blockchain_id, blockchain_explorer_id),
    FOREIGN KEY (blockchain_id) REFERENCES blockchains (id),
    FOREIGN KEY (blockchain_explorer_id) REFERENCES blockchain_explorers (id)
    );
