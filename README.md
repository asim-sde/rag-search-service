# RAG Search Service

A High-Level System Design (HLD) implementation of a Retrieval-Augmented Generation (RAG) backend utilizing document embeddings and vector similarity search.

## Tech Stack
- **Java 17 & Spring Boot**
- **PostgreSQL & `pgvector`** (for HNSW indexed similarity search)
- **Redis** (for Token Bucket rate limiting)
- **Docker & Docker Compose**

## Features
- **Vector Similarity Search:** Ingests text documents, translates them into high-dimensional vector embeddings, and performs nearest-neighbor (`<->`) distance lookups directly inside PostgreSQL.
- **Token Bucket Rate Limiter:** Protects the ingestion and search REST APIs by mapping client usage to short-lived sliding windows utilizing Redis `INCR` and `EXPIRE`.
- **Embedding Pipeline:** Provides an isolated `EmbeddingService` abstraction that mocks text-to-vector transformations, designed to easily plug into external models (e.g., OpenAI API).

## Running Locally

1. Start the PostgreSQL (with pgvector) and Redis containers:
   ```bash
   docker-compose up -d postgres redis
   ```
2. Run the Spring Boot application:
   ```bash
   ./mvnw spring-boot:run
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
