-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blocks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    operation_id UUID NOT NULL REFERENCES operations(id) ON DELETE CASCADE,
    block_type VARCHAR(20) NOT NULL CHECK (block_type IN ('header', 'footer')),
    platform VARCHAR(20) NOT NULL CHECK (platform IN ('wordpress', 'tilda', 'bitrix', 'html5', 'unknown')),
    content JSONB NOT NULL,
    html TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_blocks_operation_id ON blocks(operation_id);
CREATE INDEX IF NOT EXISTS idx_blocks_block_type ON blocks(block_type);
CREATE INDEX IF NOT EXISTS idx_blocks_platform ON blocks(platform);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blocks CASCADE;
-- +goose StatementEnd