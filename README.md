# Noosphere (Monorepo)

A private, secure, and personalized AI assistant platform modeled after industry-standard developer agents. This project runs entirely on local hardware, combining a lightweight machine learning sidecar with a high-performance orchestrator backend and dual client interfaces (Web and CLI).

## Architecture Overview

The system utilizes a decoupled, layered architecture to isolate heavy machine learning inference from core application and routing logic. This ensures that the underlying AI models can be upgraded or swapped out in the future without impacting user interfaces or data persistence layers.

                  ┌─────────────────┐
                  │ frontend-web    │ (Next.js / Browser)
                  └────────┬────────┘
                           │ HTTP / WebSockets
                           ▼
┌──────────────┐  ┌─────────────────┐  Internal HTTP  ┌────────────────────────┐
│ frontend-cli ├──►│   backend-api   ├───────────────►│  services/ai-inference │
└──────────────┘  │   (Go Server)   │  (Port 8000)   │  (Python / llama-cpp)  │
 (Bun / Ink)      └────────┬────────┘                 └────────────────────────┘
                           │ SQL / Vectors
                           ▼
                  ┌─────────────────┐
                  │  PostgreSQL DB  │ (with pgvector)
                  └─────────────────┘

### Core Components

*   apps/web-client (Frontend UI): A responsive, chat-focused web browser interface built with Next.js, TypeScript, and Tailwind CSS. It supports markdown rendering, syntax highlighting for code blocks, and real-time token streaming.
*   apps/cli-client (Terminal UI): A blazing-fast, interactive command-line developer assistant built using Bun and Ink (React for the CLI). It renders terminal layouts, loaders, and menus natively.
*   services/backend-api (Orchestrator): A highly concurrent REST API built in Go (Golang) using a strict Controller-Service-Repository pattern. It manages user state, rate limiting, conversation histories, and forwards prompt contexts to the inference engine.
*   services/ai-inference (The Brain): A lightweight Python microservice wrapping FastAPI around llama-cpp-python. It directly loads quantized .gguf models into local system RAM/VRAM and handles the raw text inference.

---

## Repository Structure

noosphere/
├── apps/
│   ├── web-client/            # Next.js frontend (Browser App)
│   └── cli-client/            # Bun & Ink frontend (Terminal App)
├── services/
│   ├── backend-api/           # Go API server (Orchestrator)
│   │   ├── cmd/api/main.go    # Server entry point
│   │   └── internal/          # Controller, Service, Repository layers
│   └── ai-inference/          # Python inference engine (Sidecar)
│       ├── models/            # Target directory for local .gguf models
│       └── app.py             # FastAPI runner
├── docker-compose.yml         # Local multi-container infrastructure orchestration
└── README.md

---

## Technical Stack & Rationale

| Layer | Technology | Selection Rationale |
| :--- | :--- | :--- |
| CLI Runtime | Bun | Instantaneous cold starts and native, configuration-free TypeScript execution. |
| CLI Layout | Ink / Yoga | Allows the use of React state and Flexbox grid alignment directly inside the terminal window. |
| Backend API | Go (Golang) | High-concurrency throughput via ultra-lightweight Goroutines (2KB) with a minimal idle memory footprint (~15MB). |
| ML Inference | Python | Industry-standard ecosystem (llama-cpp-python) providing clean high-level abstractions over low-level C++/CUDA compilation matrix math. |
| Database | PostgreSQL | Enterprise relational data handling coupled with pgvector for local document embeddings and vector searches. |

---

## Local Development Setup

### Prerequisites

*   Docker and Docker Compose
*   Go (1.22+)
*   Bun Runtime
*   Python (3.10+)

### Quick Start

1. Clone the repository:
   git clone https://github.com/astrostar41ust/Noosphere.git
   cd noosphere

2. Download an Open-Weight Model:
   Download your preferred quantized .gguf model (e.g., Llama 3 8B, Mistral 7B) from Hugging Face and place it in the models directory:
   mkdir -p services/ai-inference/models
   # Download your model file into this folder

3. Spin Up Infrastructure:
   Launch the database, Python inference sidecar, and Go orchestrator backend simultaneously using Docker Compose:
   docker compose up --build

4. Run the Web UI Interface:
   cd apps/web-client
   bun install
   bun run dev

5. Run the CLI Interface:
   cd apps/cli-client
   bun install
   bun run start

---

## Branching & Contribution Guidelines

This project strictly enforces Feature Branching workflows to preserve main branch stability.

*   main: Always stable and production-ready. Direct commits are restricted.
*   feature/*: All updates, components, and architectural enhancements must be developed on isolated feature branches and merged exclusively via approved Pull Requests.

# Workflow example:
git checkout main
git pull origin main
git checkout -b feature/your-feature-name