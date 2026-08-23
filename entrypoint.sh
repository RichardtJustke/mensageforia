#!/bin/sh
# Entrypoint do container mensageforia.
# Prepara o git (identidade + remote autenticado) antes de executar o binário.
set -eu

GIT_NAME="${GIT_AUTHOR_NAME:-RichardtJustke}"
GIT_EMAIL="${GIT_AUTHOR_EMAIL:-rj.justke@gmail.com}"

if command -v git >/dev/null 2>&1 && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  # Evita erro de "dubious ownership" quando o .git vem de outro UID (build context)
  git config --global --add safe.directory /app

  # 1. Identidade do committer
  git config --global user.name "$GIT_NAME"
  git config --global user.email "$GIT_EMAIL"

  # 2. Injeta o GITHUB_TOKEN na URL do remote para autenticar o push
  if [ -n "${GITHUB_TOKEN:-}" ]; then
    REMOTE_URL=$(git remote get-url origin 2>/dev/null || true)
    if [ -n "$REMOTE_URL" ]; then
      # Remove credenciais antigas da URL, se houver
      CLEAN_URL=$(printf '%s' "$REMOTE_URL" | sed -E 's#https://[^/@]+@#https://#')
      case "$CLEAN_URL" in
        https://github.com/*)
          AUTH_URL=$(printf '%s' "$CLEAN_URL" | sed -E "s#^https://#https://x-access-token:${GITHUB_TOKEN}@#")
          git remote set-url origin "$AUTH_URL"
          ;;
        git@github.com:*)
          AUTH_URL=$(printf '%s' "$CLEAN_URL" | sed -E "s#^git@github.com:#https://x-access-token:${GITHUB_TOKEN}@github.com/#")
          git remote set-url origin "$AUTH_URL"
          ;;
      esac
    fi
  fi
  # 3. Commit + push de mensagens pendentes (previne perda em reinicialização)
  if [ -n "${GITHUB_TOKEN:-}" ]; then
    UNTRACKED=$(git ls-files --others --exclude-standard -- messages/ 2>/dev/null || true)
    MODIFIED=$(git diff --name-only -- messages/ 2>/dev/null || true)
    if [ -n "$UNTRACKED" ] || [ -n "$MODIFIED" ]; then
      echo "[entrypoint] Mensagens pendentes encontradas, fazendo commit+push..."
      git add messages/
      git commit -m "mensagem automática: pendências ($(date +%Y-%m-%d-%Hh%M))" || true
      git push || echo "[entrypoint] WARN: push falhou, tentando novamente após o app iniciar"
    else
      echo "[entrypoint] Nenhuma mensagem pendente."
    fi
  fi
fi

exec "$@"
