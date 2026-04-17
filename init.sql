CREATE TABLE IF NOT EXISTS market_indicators (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    pair VARCHAR(20) NOT NULL,
    timestamp_unix BIGINT NOT NULL,
    mms_20 DECIMAL(18, 8) DEFAULT NULL,
    mms_50 DECIMAL(18, 8) DEFAULT NULL,
    mms_200 DECIMAL(18, 8) DEFAULT NULL,

    INDEX idx_pair_timestamp (pair, timestamp_unix)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;