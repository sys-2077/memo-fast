# memo-fast

Your AI agent forgets everything between sessions. memo-fast fixes that.

It indexes your codebase -- functions, classes, git history, dependency graphs -- into a vector database you own. Next time your agent opens the project, it picks up exactly where you left off.

```
"what changed while I was away?"     -->  memory_resume
"how does the auth flow work?"       -->  memory_query
"explain OrderProcessor in depth"    -->  memory_deepdive
"show me the architecture topology"  -->  memory_spatial_map
```

## How it works

memo-fast runs as an MCP server. Your agent connects to it and gets 8 tools for searching, navigating, and understanding your codebase through semantic memory.

Your code is processed into embeddings and stored in a vector database **you control**. Pinecone, Qdrant, or ChromaDB -- your account, your keys, your data. memo-fast processes everything in-memory and stores nothing on its side.

```
MCP client (Claude Code, Cursor, etc.)
    |
    | 8 tools via MCP protocol
    v
memo-fast (extract + embed)  -->  Your vector DB
```

## Tools

*CLI* = also available via the CLI method (see below)

| Tool | What it does | |
|------|-------------|:---:|
| `memory_resume` | Recovers context from recent git commits. Start every session here. | |
| `memory_resume_extended` | Like resume, but also resolves the dependency graph of every touched entity. | |
| `memory_refresh` | Context for files you are currently editing (max 10 files, 60s cooldown). | *CLI* |
| `memory_query` | Semantic search. Ask in natural language, get code, commits, and files. | |
| `memory_deepdive` | Finds an entity by name, pulls its source and direct dependencies. | |
| `memory_spatial_map` | Full dependency graph with PageRank ranking and community clustering. | |
| `memory_enrich` | Generates semantic descriptions for code entities (improves search quality). | |
| `memory_index` | Index files and commits (max 30 files per call). | *CLI* |

## Getting started

1. Subscribe at [mcpize.com/memo-fast](https://mcpize.com)
2. Add your vector DB credentials (Pinecone, Qdrant, or ChromaDB) in the dashboard
3. Add memo-fast to your MCP client config
4. Open your project and ask your agent anything about the codebase

## The CLI Method

For initial indexing and automatic updates, memo-fast includes a lightweight CLI. It reads your project files and git history directly from disk and sends them to the server for processing.

### Install

```bash
# macOS (Apple Silicon)
curl -fsSL https://github.com/codexor/memo-fast-product/releases/latest/download/memo-fast-darwin-arm64.tar.gz | tar xz -C /usr/local/bin

# macOS (Intel)
curl -fsSL https://github.com/codexor/memo-fast-product/releases/latest/download/memo-fast-darwin-amd64.tar.gz | tar xz -C /usr/local/bin

# Linux
curl -fsSL https://github.com/codexor/memo-fast-product/releases/latest/download/memo-fast-linux-amd64.tar.gz | tar xz -C /usr/local/bin

# Windows
# Download memo-fast-windows-amd64.zip from GitHub Releases and add to PATH
```

### Usage

```bash
memo-fast init                  # One-time setup per project
memo-fast index                 # Index all files + recent commits
memo-fast index --incremental   # Index only the last commit
memo-fast hook install          # Auto-index on every git commit
```

Once the hook is installed, your memory updates itself after every commit. No manual re-indexing needed.

## Credentials

You bring your own storage. Configure one of these backends:

| Credential | Required | Description |
|------------|----------|-------------|
| `VECTOR_BACKEND` | No | `pinecone` (default), `qdrant`, or `chroma` |
| `PINECONE_API_KEY` | If Pinecone | Your API key |
| `PINECONE_HOST` | If Pinecone | Your index host URL |
| `PINECONE_INDEX` | No | Index name (default: `memo-fast`) |
| `QDRANT_URL` | If Qdrant | Your Qdrant Cloud URL |
| `QDRANT_API_KEY` | If Qdrant | Your API key |

## Data processing

memo-fast is stateless. Your source code is processed in-memory to generate embeddings, then written directly to your vector database. We do not store, cache, or log your code. See [DATA-PROCESSING.md](DATA-PROCESSING.md) for details.

## License to Use

See [MCPize terms](https://mcpize.com/terms).
