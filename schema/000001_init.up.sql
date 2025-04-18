CREATE TABLE volume24 (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    volume24 NUMERIC NOT NULL,
    ethereum NUMERIC,
    bitcoin NUMERIC,
    solana NUMERIC,
    bsc NUMERIC,
    tron NUMERIC
);

CREATE TABLE fees24 (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    fees24 NUMERIC NOT NULL,
    ethereum NUMERIC,
    bitcoin NUMERIC,
    solana NUMERIC,
    bsc NUMERIC,
    tron NUMERIC
);

CREATE TABLE revenue24 (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    revenue24 NUMERIC NOT NULL,
    ethereum NUMERIC,
    bitcoin NUMERIC,
    solana NUMERIC,
    bsc NUMERIC,
    tron NUMERIC
);

CREATE TABLE tvl_Chains (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    tvl BIGINT NOT NULL
);

CREATE TABLE tvl_Protocols (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    tvl NUMERIC NOT NULL,

    -- Основные сети
    ethereum NUMERIC,
    solana NUMERIC,
    binance NUMERIC,
    bitcoin NUMERIC,
    tron NUMERIC,


    -- Побочные сети (для каждой основной сети)
    ethereum_borrowed NUMERIC,
    ethereum_pool2 NUMERIC,
    ethereum_staking NUMERIC,
    ethereum_treasury NUMERIC,
    ethereum_vesting NUMERIC,
    ethereum_own_tokens NUMERIC,

    solana_borrowed NUMERIC,
    solana_pool2 NUMERIC,
    solana_staking NUMERIC,
    solana_vesting NUMERIC,

    binance_borrowed NUMERIC,
    binance_pool2 NUMERIC,
    binance_staking NUMERIC,
    binance_vesting NUMERIC,

    bitcoin_staking NUMERIC,

    tron_borrowed NUMERIC,
    tron_pool2 NUMERIC,
    tron_staking NUMERIC
);

CREATE TABLE Mcap_Protocols (
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    mcap BIGINT NOT NULL
);

CREATE TABLE Coingecko (
    created_at NUMERIC NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    volume24 NUMERIC NOT NULL,
    mcap NUMERIC NOT NULL
);