# Plano de Integração Frontend - Next.js + Ory Kratos + APISIX

## Visão Geral da Arquitetura

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Browser (localhost:9080)                    │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    APISIX Gateway (porta 9080)                      │
│  ┌─────────────┬────────────────┬─────────────────┬──────────────┐ │
│  │  /api/*     │ /.ory/kratos/* │    /auth/*      │     /*       │ │
│  │ (protected) │   (público)    │   (público)     │  (público)   │ │
│  └──────┬──────┴───────┬────────┴────────┬────────┴──────┬───────┘ │
└─────────┼──────────────┼─────────────────┼───────────────┼─────────┘
          │              │                 │               │
          ▼              ▼                 ▼               ▼
     ┌─────────┐   ┌──────────┐    ┌────────────┐   ┌───────────┐
     │ Backend │   │  Kratos  │    │ Kratos UI  │   │ Frontend  │
     │  :8080  │   │  :4433   │    │   :3000    │   │   :3000   │
     └─────────┘   └──────────┘    └────────────┘   └───────────┘
```

---

## 1. Pontos de Falha Críticos

### 1.1 Conflito de Portas (CRÍTICO)

**Problema:**
- Frontend Dockerfile expõe porta `3000`
- docker-compose mapeia `frontend:3000 → host:80` (incorreto)
- Kratos UI também usa porta `3000` internamente

**Impacto:** Container não inicia ou APISIX não consegue rotear.

**Solução:**
```yaml
# docker-compose.yml
frontend:
  build:
    context: ../front
  expose:
    - "3000"  # Apenas expõe internamente, não mapeia para host
  networks:
    - intranet
  # REMOVER: ports: - "3000:80"
```

---

### 1.2 Variáveis de Ambiente Incorretas (CRÍTICO)

**Problema:**
- `.env.local` aponta para Ory Cloud: `https://dazzling-wilbur-dglb1ak9pu.projects.oryapis.com`
- Deve apontar para Kratos local via APISIX

**Impacto:** Frontend tenta autenticar com Ory Cloud em vez do Kratos local.

**Solução - Criar `.env.docker`:**
```env
# URLs do Gateway (acessível do browser)
NEXT_PUBLIC_ORY_SDK_URL=http://127.0.0.1:9080/.ory/kratos/public
NEXT_PUBLIC_BACKEND_URL=http://127.0.0.1:9080/api

# URLs internas (server-side dentro do Docker)
ORY_SDK_URL=http://kratos:4433
BACKEND_API_URL=http://backend:8080
```

**docker-compose.yml:**
```yaml
frontend:
  build:
    context: ../front
  environment:
    - NEXT_PUBLIC_ORY_SDK_URL=http://127.0.0.1:9080/.ory/kratos/public
    - ORY_SDK_URL=http://kratos:4433
    - BACKEND_API_URL=http://backend:8080
```

---

### 1.3 Rewrites do Next.js Conflitam com APISIX (ALTO)

**Problema:**
- `next.config.ts` define rewrites para `/self-service/:path*` → Ory SDK
- Quando Next.js roda atrás do APISIX, os rewrites podem conflitar

**Impacto:** Fluxos de autenticação falham com 404 ou loops de redirect.

**Solução - Atualizar `next.config.ts`:**
```typescript
// next.config.ts
const nextConfig: NextConfig = {
  output: "standalone",
  reactStrictMode: true,

  async rewrites() {
    const oryUrl = process.env.ORY_SDK_URL || "http://kratos:4433";

    return [
      {
        source: "/self-service/:path*",
        destination: `${oryUrl}/self-service/:path*`,
      },
      // Proxy API calls internamente (server-side)
      {
        source: "/api/kratos/:path*",
        destination: `${oryUrl}/:path*`,
      },
    ];
  },
};
```

---

### 1.4 Cookies Cross-Origin (CRÍTICO)

**Problema:**
- Cookie `ory_kratos_session` é setado por Kratos em `127.0.0.1:4433`
- Browser acessa via `127.0.0.1:9080`
- Cookies não são enviados devido a SameSite/Domain mismatch

**Impacto:** Sessão nunca é reconhecida, usuário sempre "deslogado".

**Solução - Kratos `kratos.yml`:**
```yaml
# Já configurado, mas verificar:
cookies:
  domain: "127.0.0.1"
  path: "/"
  same_site: Lax  # Permite cookies em navegação top-level

serve:
  public:
    base_url: http://127.0.0.1:9080/.ory/kratos/public/
    cors:
      enabled: true
      allowed_origins:
        - http://127.0.0.1:9080
        - http://localhost:9080
      allow_credentials: true
```

---

### 1.5 Session Verification Server-Side (ALTO)

**Problema:**
- `src/core/ory/session.ts` chama Kratos diretamente
- Dentro do container, deve usar URL interna (`http://kratos:4433`)
- Fora do container, usa URL do gateway

**Impacto:** `getOrySession()` falha com `ECONNREFUSED` ou timeout.

**Solução - Atualizar `session.ts`:**
```typescript
// src/core/ory/session.ts
import { Configuration, FrontendApi } from "@ory/client";

const getOryClient = () => {
  // Server-side: usar URL interna do Docker
  // Client-side: usar URL do gateway
  const baseUrl = typeof window === "undefined"
    ? process.env.ORY_SDK_URL || "http://kratos:4433"
    : process.env.NEXT_PUBLIC_ORY_SDK_URL || "http://127.0.0.1:9080/.ory/kratos/public";

  return new FrontendApi(
    new Configuration({
      basePath: baseUrl,
      baseOptions: {
        withCredentials: true,
      },
    })
  );
};
```

---

### 1.6 Dockerfile Multi-Stage Build (MÉDIO)

**Problema:**
- Build-time não tem acesso às variáveis de ambiente do Docker Compose
- `NEXT_PUBLIC_*` são inlined no build

**Impacto:** URLs hardcoded no bundle, impossível mudar em runtime.

**Solução - Build com ARGs:**
```dockerfile
# front/Dockerfile
FROM oven/bun:alpine AS builder
WORKDIR /app

# Build args para variáveis públicas
ARG NEXT_PUBLIC_ORY_SDK_URL=http://127.0.0.1:9080/.ory/kratos/public
ARG NEXT_PUBLIC_BACKEND_URL=http://127.0.0.1:9080/api

ENV NEXT_PUBLIC_ORY_SDK_URL=$NEXT_PUBLIC_ORY_SDK_URL
ENV NEXT_PUBLIC_BACKEND_URL=$NEXT_PUBLIC_BACKEND_URL

# ... resto do build
```

**docker-compose.yml:**
```yaml
frontend:
  build:
    context: ../front
    args:
      NEXT_PUBLIC_ORY_SDK_URL: http://127.0.0.1:9080/.ory/kratos/public
      NEXT_PUBLIC_BACKEND_URL: http://127.0.0.1:9080/api
```

---

### 1.7 Health Check do Frontend (BAIXO)

**Problema:**
- Sem healthcheck, APISIX pode rotear para container não-pronto

**Solução:**
```yaml
# docker-compose.yml
frontend:
  healthcheck:
    test: ["CMD", "wget", "-q", "--spider", "http://localhost:3000"]
    interval: 10s
    timeout: 5s
    retries: 3
    start_period: 30s
```

---

## 2. Matriz de Compatibilidade de URLs

| Contexto | ORY_SDK_URL | BACKEND_API_URL |
|----------|-------------|-----------------|
| Browser (client-side) | `http://127.0.0.1:9080/.ory/kratos/public` | `http://127.0.0.1:9080/api` |
| Next.js Server (SSR) | `http://kratos:4433` | `http://backend:8080` |
| Server Actions | `http://kratos:4433` | `http://backend:8080` |

---

## 3. Ordem de Inicialização

```
1. postgres-kratos     (healthcheck: pg_isready)
2. kratos-migrate      (depends_on: postgres healthy)
3. kratos              (depends_on: migrate completed)
4. backend             (sem dependências)
5. frontend            (healthcheck: wget)
6. etcd                (sem dependências)
7. apisix              (depends_on: etcd, kratos, frontend?)
```

**Recomendação:** Adicionar `depends_on` para garantir ordem:
```yaml
apisix:
  depends_on:
    etcd:
      condition: service_started
    kratos:
      condition: service_started
    frontend:
      condition: service_healthy
    backend:
      condition: service_started
```

---

## 4. Checklist de Debugging

### 4.1 Testar Conectividade Interna
```bash
# Entrar no container frontend
docker compose exec frontend sh

# Testar Kratos
wget -qO- http://kratos:4433/health/ready

# Testar Backend
wget -qO- http://backend:8080/health
```

### 4.2 Testar Cookies
```bash
# Login e capturar cookie
curl -c cookies.txt -b cookies.txt \
  http://127.0.0.1:9080/.ory/kratos/public/self-service/login/browser

# Verificar sessão com cookie
curl -b cookies.txt \
  http://127.0.0.1:9080/.ory/kratos/public/sessions/whoami
```

### 4.3 Testar Rota Protegida
```bash
# Sem cookie (deve retornar 401)
curl -i http://127.0.0.1:9080/api/test

# Com cookie válido
curl -i -b cookies.txt http://127.0.0.1:9080/api/test
```

---

## 5. Arquivos a Modificar

| Arquivo | Ação | Prioridade |
|---------|------|------------|
| `docker/docker-compose.yml` | Corrigir portas, adicionar env vars | CRÍTICO |
| `front/Dockerfile` | Adicionar ARGs para build | CRÍTICO |
| `front/next.config.ts` | Atualizar rewrites | ALTO |
| `front/src/core/ory/session.ts` | URL dinâmica server/client | ALTO |
| `docker/kratos/kratos.yml` | Verificar cookie domain | MÉDIO |
| `docker/apisix/apisix.yaml` | Ajustar upstream frontend | MÉDIO |

---

## 6. Variáveis de Ambiente Finais

```env
# Frontend Container
NEXT_PUBLIC_ORY_SDK_URL=http://127.0.0.1:9080/.ory/kratos/public
NEXT_PUBLIC_BACKEND_URL=http://127.0.0.1:9080/api
ORY_SDK_URL=http://kratos:4433
BACKEND_API_URL=http://backend:8080
NODE_ENV=production
```

---

## 7. Validação Final

```bash
# 1. Build e subir
cd docker && docker compose up --build -d

# 2. Verificar logs de erro
docker compose logs frontend | grep -i error
docker compose logs apisix | grep -i error

# 3. Testar fluxo completo
# a) Acessar http://127.0.0.1:9080/auth/registration
# b) Criar conta
# c) Verificar redirect para dashboard
# d) Acessar rota protegida /api/*
```
