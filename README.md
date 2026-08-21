# Mensageforia

> Gere mensagens motivacionais automaticamente com IA, salve em Markdown + SQLite e faça push pro GitHub — todo dia, sem fazer nada.

---

## Como funciona

```
┌─────────────────────────────────────────────────────────────┐
│                     Mensageforia                            │
│                                                             │
│  ⏰ 01:00 ──► 🎲 Sorteio (1–10 mensagens no dia)           │
│                                                             │
│  🕐 Horários ──► 🧠 Ollama gera mensagem motivacional       │
│                   │                                         │
│                   ├──► 💾 SQLite (histórico)                 │
│                   ├──► 📝 messages/YYYY-MM-DD-HH.md          │
│                   └──► 🚀 git commit + push pro GitHub       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Features

- **Sorteio diário** — às 01:00 um número aleatório (1–10) define quantas mensagens o dia terá
- **Janela configurável** — horários distribuídos igualmente (padrão 08:00–20:00)
- **Persistência dupla** — SQLite (histórico estruturado) + Markdown (arquivos legíveis)
- **Backup automático** — cada mensagem é committada e enviada pro GitHub
- **Hot reload** — altere a janela de horários sem rebuildar o container
- **Zero dependência externa** — scheduler nativo em Go (sem cron library)

---

## Stack

| Componente | Tecnologia |
|---|---|
| Backend | Go 1.26 |
| LLM | Ollama (llama3.2:1b) |
| Banco | SQLite (modernc.org/sqlite) |
| Agendamento | Scheduler nativo Go |
| Container | Docker + Docker Compose |
| Deploy | VPS / Home Server |

---

## Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose v2
- [Git](https://git-scm.com/)
- Conta no GitHub com um **Personal Access Token** (fine-grained)
  - Permissão: `Contents → Read and write`
  - Criar em: [GitHub Settings → Developer settings → Fine-grained tokens](https://github.com/settings/tokens?type=beta)

---

## Quick Start

```bash
# 1. Clone o repositório
git clone https://github.com/RichardtJustke/mensageforia.git
cd mensageforia

# 2. Configure o token
cp .env.example .env
nano .env   # cole seu GITHUB_TOKEN

# 3. Suba os containers
docker compose up -d --build

# 4. Baixe o modelo (uma vez só)
docker compose exec ollama ollama pull llama3.2:1b

# 5. Teste na hora
docker compose run --rm app ./mensageforia --once
```

Pronto! O bot vai gerar uma mensagem, salvar no SQLite + Markdown e fazer push pro GitHub.

---

## Configuração

Todas as variáveis ficam no arquivo `.env` (não vai pro git).

| Variável | Obrigatória | Default | Descrição |
|---|---|---|---|
| `GITHUB_TOKEN` | ✅ | — | PAT do GitHub com permissão Contents: write |
| `OLLAMA_MODEL` | | `llama3.2:1b` | Modelo do Ollama |
| `OLLAMA_BASE_URL` | | `http://ollama:11434` | Endereço do servidor Ollama |
| `TZ` | | `America/Sao_Paulo` | Fuso horário |
| `MSG_WINDOW_START` | | `08:00` | Início da janela de mensagens |
| `MSG_WINDOW_END` | | `20:00` | Fim da janela de mensagens |
| `GIT_AUTHOR_NAME` | | `MensageForia Bot` | Nome no commit |
| `GIT_AUTHOR_EMAIL` | | `bot@mensageforia.local` | Email no commit |

---

## Scheduler

O scheduler funciona em dois momentos:

### 🎲 Sorteio (01:00)

Às 01:00 da manhã, o sistema sorteia um número de **1 a 10** que define quantas mensagens serão geradas naquele dia.

### ⏰ Distribuição

Os horários são distribuídos igualmente na janela configurável:

| Sorteio | Horários gerados |
|---|---|
| 1 mensagem | 08:00 |
| 3 mensagens | 08:00, 14:00, 20:00 |
| 5 mensagens | 08:00, 11:00, 14:00, 17:00, 20:00 |
| 10 mensagens | 08:00, 09:20, 10:40, ..., 20:00 |

### Alterar a janela

Edite o `.env` e reinicie:

```bash
# Exemplo: mensagens entre 06:00 e 22:00
MSG_WINDOW_START=06:00
MSG_WINDOW_END=22:00

docker compose restart app
```

---

## Estrutura do projeto

```
mensageforia/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── git/git.go              # Commit + push via os/exec
│   ├── message/message.go      # Orquestração: prompt → Ollama → SQLite → md → git
│   ├── ollama/trigger.go       # Cliente HTTP pro Ollama
│   ├── scheduler/scheduler.go  # Sorteio diário + distribuição de horários
│   └── storage/sqlite.go       # Schema + persistência SQLite
├── data/                       # Banco SQLite (gitignored)
├── messages/                   # Mensagens .md (commitadas pelo bot)
├── Dockerfile                  # Build multi-stage
├── docker-compose.yml          # Orquestração dos containers
├── entrypoint.sh               # Configura git + injeta token
├── .env.example                # Template das variáveis
└── go.mod
```

---

## Comandos úteis

```bash
# Ver logs em tempo real
docker compose logs -f app

# Gerar mensagem na hora (testar)
docker compose run --rm app ./mensageforia --once

# Ver mensagens no SQLite
python3 -c "
import sqlite3
for row in sqlite3.connect('data/mensageforia.db').execute('SELECT id, theme, timestamp FROM messages'):
    print(row)
"

# Ver sortear do dia (logs)
docker compose logs app | grep sorteio

# Reiniciar após mudar .env
docker compose restart app

# Rebuild completo
docker compose up -d --build
```

---

## Troubleshooting

| Problema | Solução |
|---|---|
| `Permission denied` no push | Verifique se o token tem Contents: Read and write |
| Ollama não responde | `docker compose logs ollama` — espere o healthcheck passar |
| Mensagem não aparece no GitHub | Verifique os logs: `docker compose logs app \| grep push` |
| Porta 11434 ocupada | O compose usa 11435 no host — conflito com Ollama local resolvido |
| Container não reinicia após reboot | `docker compose up -d` — o `restart: unless-stopped` cuida do resto |

---

## Licença

Projeto pessoal — sem licença formal.
