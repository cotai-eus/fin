#!/usr/bin/env bash
set -euo pipefail

# gen.sh - generate strong random secrets and self-signed certs for local/dev
# Usage: ./gen.sh

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
CERT_DIR="$ROOT_DIR/certs"
ENV_FILE="$ROOT_DIR/.env"
EXAMPLE_ENV="$ROOT_DIR/.env.example"

mkdir -p "$CERT_DIR"

timestamp() { date +%s; }
rand_hex() { local bytes=${1:-16}; openssl rand -hex "$bytes"; }

backup_env() {
  if [ -f "$ENV_FILE" ]; then
    cp "$ENV_FILE" "$ENV_FILE.bak.$(timestamp)"
    echo "Existing .env backed up to $ENV_FILE.bak.$(timestamp)"
  fi
}

echo "Generating secrets and certificates in: $ROOT_DIR"
backup_env

# Generate secrets
CORE_DB_PASSWORD=$(rand_hex 16)
CORE_DATABASE_URL="postgres://postgres:${CORE_DB_PASSWORD}@db-core:5432/lauratech?sslmode=disable"
ENCRYPTION_KEY=$(rand_hex 16)

KRATOS_DB_PASSWORD=$(rand_hex 16)
KRATOS_DSN="postgres://kratos:${KRATOS_DB_PASSWORD}@postgres-kratos:5432/kratos?sslmode=disable&max_conns=20&max_idle_conns=4"

KRATOS_COOKIE_SECRET=$(rand_hex 16)
KRATOS_CIPHER_SECRET=$(rand_hex 16)
KRATOS_SMTP_CONNECTION_URI="smtps://user:pass@mailslurper:1025/?skip_ssl_verify=true"

KRATOS_UI_COOKIE_SECRET=$(rand_hex 16)
KRATOS_UI_CSRF_SECRET=$(rand_hex 16)

ETCD_ROOT_PASSWORD=$(rand_hex 32)

# Write .env with organized sections matching .env.example
cat > "$ENV_FILE" <<EOF
# Core database
CORE_DB_PASSWORD=${CORE_DB_PASSWORD}
CORE_DATABASE_URL=${CORE_DATABASE_URL}

# Backend encryption
ENCRYPTION_KEY=${ENCRYPTION_KEY}

# Kratos database
KRATOS_DB_PASSWORD=${KRATOS_DB_PASSWORD}
KRATOS_DSN=${KRATOS_DSN}

# Kratos secrets
KRATOS_COOKIE_SECRET=${KRATOS_COOKIE_SECRET}
KRATOS_CIPHER_SECRET=${KRATOS_CIPHER_SECRET}
KRATOS_SMTP_CONNECTION_URI=${KRATOS_SMTP_CONNECTION_URI}

# Kratos selfservice UI
KRATOS_UI_COOKIE_SECRET=${KRATOS_UI_COOKIE_SECRET}
KRATOS_UI_CSRF_SECRET=${KRATOS_UI_CSRF_SECRET}

# ETCD
ETCD_ROOT_PASSWORD=${ETCD_ROOT_PASSWORD}
EOF

chmod 600 "$ENV_FILE"
echo "Wrote $ENV_FILE (permissions 600)"

### Generate a local CA and server certs (apisix, etcd, kratos)
if [ ! -f "$CERT_DIR/ca.key" ]; then
  echo "Generating CA (4096-bit RSA)"
  openssl req -x509 -nodes -newkey rsa:4096 -sha256 -days 3650 \
    -keyout "$CERT_DIR/ca.key" -out "$CERT_DIR/ca.crt" -subj "/CN=Local Dev CA"
  chmod 600 "$CERT_DIR/ca.key"
fi

gen_cert() {
  local name=$1
  local sans=$2
  local key="$CERT_DIR/${name}.key"
  local crt="$CERT_DIR/${name}.crt"
  local csr="$CERT_DIR/${name}.csr"
  local extf

  if [ -f "$key" ] && [ -f "$crt" ]; then
    echo "Certificate for $name already exists, skipping"
    return
  fi

  extf=$(mktemp)
  printf "%s" "subjectAltName=${sans}" > "$extf"

  openssl req -new -nodes -newkey rsa:2048 -keyout "$key" -subj "/CN=${name}" -out "$csr"
  openssl x509 -req -in "$csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial \
    -out "$crt" -days 3650 -sha256 -extfile "$extf"

  rm -f "$csr" "$extf"
  chmod 600 "$key"
  echo "Generated cert: $crt"
}

gen_cert "apisix" "DNS:apisix, DNS:localhost, IP:127.0.0.1"
gen_cert "etcd" "DNS:etcd, DNS:localhost, IP:127.0.0.1"
gen_cert "kratos" "DNS:kratos, DNS:localhost, IP:127.0.0.1"

echo "Certificates and keys written to: $CERT_DIR"

### Ensure .gitignore contains .env
GITIGNORE="$ROOT_DIR/.gitignore"
if [ -f "$GITIGNORE" ]; then
  if ! grep -qxF ".env" "$GITIGNORE"; then
    echo ".env" >> "$GITIGNORE"
    echo "Appended .env to $GITIGNORE"
  fi
else
  echo ".env" > "$GITIGNORE"
  echo "Created $GITIGNORE"
fi

echo "Done. Next: run 'docker compose build --no-cache' from $ROOT_DIR and 'docker compose up -d'"
echo "Keep $ENV_FILE private (it's ignored by .gitignore)."

exit 0