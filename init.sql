CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS documents (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    embedding vector(3) -- Using a tiny 3-dimensional vector for mocked embeddings
);

-- Create an HNSW index for fast nearest neighbor search
CREATE INDEX ON documents USING hnsw (embedding vector_l2_ops);
