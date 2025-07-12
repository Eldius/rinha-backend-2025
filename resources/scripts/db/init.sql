
CREATE TABLE payment_info (
    id SERIAL PRIMARY KEY,
    correlation_id VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    requested_at TIMESTAMP NOT NULL,
    message TEXT,
    provider VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_info_correlation_id ON payment_info(correlation_id);
CREATE INDEX idx_payment_info_status ON payment_info(status);
CREATE INDEX idx_payment_info_provider ON payment_info(provider);
