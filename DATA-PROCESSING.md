# Data Processing Agreement

**Effective date:** March 2026
**Service:** memo-fast (MCP server for project memory)
**Provider:** Codexor

## 1. Scope

This document describes how memo-fast processes data when you use the service. By subscribing and configuring your credentials, you accept these terms.

## 2. What we process

When you invoke memo-fast tools, the service receives:

- Source code files (content and file paths)
- Git commit metadata (SHA, subject, date, changed files)
- Search queries you send to the MCP tools

This data is processed **in runtime only** to perform extraction, embedding, and vector database operations.

## 3. What we store

**Nothing.** memo-fast is stateless. Your data passes through our infrastructure during request processing and is discarded when the request completes. We do not persist, cache, log, or copy your source code or commit content.

Your vectors and metadata are stored exclusively in **your own** vector database account (Pinecone, Qdrant, or ChromaDB), using credentials you provide.

## 4. BYOK / BYOS model

memo-fast operates on a Bring Your Own Keys / Bring Your Own Storage model:

- **Your keys**: You provide your own vector database API keys. We use them only during request processing.
- **Your storage**: Embeddings and metadata are written to your database account. You control access, retention, and deletion.
- **Your responsibility**: You are responsible for managing your vector database account, including costs, backups, and access control.

## 5. Embeddings

Source code is converted to vector embeddings using a local ONNX model (BAAI/bge-base-en-v1.5) running on our infrastructure. Embeddings are generated in-memory and sent directly to your vector database. No external embedding APIs (OpenAI, etc.) are used.

## 6. What we do log

- Request count and latency (for billing and monitoring)
- Error codes (without request content)
- Subscriber ID (for authentication)

We do **not** log source code, file contents, commit bodies, or query text.

## 7. Third-party services

| Service | Purpose | Data shared |
|---------|---------|-------------|
| MCPize | Gateway, billing, authentication | Subscriber ID, request count |
| Cloud compute | Processing infrastructure | In-memory request processing |
| Your vector DB | Storage (your account) | Embeddings + metadata |

## 8. Data deletion

Since we store no user data, there is nothing to delete on our side. To remove your indexed data, delete the collection/namespace in your vector database account.

## 9. Security

- All communication is encrypted in transit (TLS/HTTPS)
- Service-to-service authentication via identity tokens
- Subscriber credentials are encrypted at rest by MCPize
- No shared storage between subscribers

## 10. Changes

We may update this document. Changes take effect for new subscribers immediately. Existing subscribers are notified via MCPize.

## Contact

For questions about data processing: israel@codexor.dev
