#!/bin/sh
# Entrypoint do container mensageforia.
# Prepara o git (identidade + remote autenticado) antes de executar o binário.
set -eu

GIT_NAME="${GIT_AUTHOR_NAME:-MensageForia Bot}"
GIT_EMAIL="${GIT_AUTHOR_EMAIL:-bot@mensageforia.local}"

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
fi

exec "$@"
