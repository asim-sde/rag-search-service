# RAG Search Service

A High-Level System Design (HLD) implementation of a Retrieval-Augmented Generation (RAG) backend utilizing document embeddings and vector similarity search.

## Tech Stack
- **Go (Golang)**: High performance, lightweight backend routing using `go-chi/chi`.
- **PostgreSQL & `pgvector`**: Utilizes HNSW indexed similarity search to ingest text documents and perform lightning-fast nearest-neighbor (`<->`) distance lookups directly inside PostgreSQL.
- **Redis**: Token Bucket rate limiting built with `go-redis` to protect the search and ingestion APIs.
- **Docker & Docker Compose**: Containerized database and cache layer.

## Features
- **Vector Similarity Search:** Ingests text documents, translates them into high-dimensional vector embeddings, and executes distance queries over the pgvector extension natively using raw SQL strings for max compatibility.
- **Token Bucket Rate Limiter:** Protects the ingestion and search REST APIs by mapping client usage to short-lived sliding windows utilizing Redis `INCR` and `EXPIRE`.
- **Embedding Pipeline:** Provides an isolated `EmbeddingService` abstraction that mocks text-to-vector transformations, designed to easily plug into external models (e.g., OpenAI API).

## Setup Instructions

### 1. Prerequisites
- Docker & Docker Compose
- Go 1.23+

### 2. Local Infrastructure setup
Spin up PostgreSQL (with pgvector) and Redis containers:
```bash
docker-compose up -d postgres redis
```

### 3. Running the Server
Start the Go application server:
```bash
go run cmd/server/main.go
```

## API Endpoints

### 1. Ingest Document
```http
POST /api/documents
X-Client-Id: user-123
Content-Type: application/json

{
    "content": "A high performance distributed system design requires caching."
}
```

### 2. Search Similar Documents
```http
POST /api/search
X-Client-Id: user-123
Content-Type: application/json

{
    "query": "distributed caching"
}
```
