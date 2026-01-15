# Fintech MVP MVP - Documentação do Projeto

Bem-vindo à documentação oficial do MVP da Fintech. Este projeto tem como objetivo validar hipóteses centrais de aquisição, retenção e transacionalidade em um ambiente controlado (PoC).

## Índice

1. [Checklist do MVP](./mvp-checklist.md) - Escopo funcional e requisitos mínimos.
2. [Arquitetura Técnica](./architecture.md) - Detalhes sobre a stack (Go, Next.js, Ory, Infra).
3. [Roadmap e Fases](./roadmap.md) - Cronograma sugerido para a Prova de Conceito.

## Visão Geral

O objetivo desta Prova de Conceito (PoC) é lançar um produto financeiro mínimo viável para validar a demanda e a operabilidade técnica.

### Pilares
- **Simplicidade:** Foco estrito no "Happy Path" do usuário.
- **Segurança:** Conformidade básica desde o dia 0 (Criptografia, 2FA).
- **Observabilidade:** Métricas claras para decisão de Go/No-Go.

## Stack Tecnológico Rápido

- **Backend:** Go (Gin framework)
- **Frontend:** Next.js 16 (App Router) + Bun
- **Banco de Dados:** PostgreSQL (Transacional) + Redis (Cache/Token blacklisting)
- **Autenticação:** Ory Kratos (Gerenciamento de Identidade) + APISIX (Gateway)
- **Infraestrutura:** Docker & Docker Compose
