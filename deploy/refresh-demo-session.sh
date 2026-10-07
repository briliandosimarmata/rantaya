#!/bin/sh
set -eu

cd "$(dirname "$0")/.."
demo_env_file=${DEMO_ENV_FILE:-/opt/rantaya-demo/.env}
compose() {
    docker compose --env-file "$demo_env_file" -f deploy/compose.demo.yaml "$@"
}

# Check the running API, not just a requested configuration file.
compose exec -T api sh -c '
    test "$APP_ENV" = demo &&
    test "$DEMO_LOGIN" = true &&
    test "$SEED_DEMO" = true
' </dev/null
compose exec -T db psql -X --set=ON_ERROR_STOP=1 -U ruang -d ruang < deploy/refresh-demo-session.sql
