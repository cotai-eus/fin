# ---------- BUILD ----------
FROM oven/bun:alpine AS builder
WORKDIR /app

# Build args para variáveis públicas (inlined no build)
ARG NEXT_PUBLIC_ORY_SDK_URL=http://127.0.0.1:9080/.ory/kratos/public
ARG NEXT_PUBLIC_BACKEND_URL=http://127.0.0.1:9080/api

ENV NEXT_PUBLIC_ORY_SDK_URL=$NEXT_PUBLIC_ORY_SDK_URL
ENV NEXT_PUBLIC_BACKEND_URL=$NEXT_PUBLIC_BACKEND_URL

# Dependências mínimas para Prisma e builds nativos
RUN apk add --no-cache libc6-compat openssl

# Instala dependências
COPY package.json bun.lockb* ./
RUN bun i

# Código-fonte
COPY . .

# Build do Next (standalone)
RUN bun run build
RUN apk del --no-cache zlib || true

# ---------- RUNTIME ----------
FROM oven/bun:alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0
EXPOSE 3000

# Apenas libs necessárias em runtime (sharp → vips, se usar)
RUN apk add --no-cache libc6-compat openssl
RUN apk del --no-cache zlib || true

# Copia o output standalone
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

# O standalone já inclui node_modules necessários
CMD ["bun", "server.js"]