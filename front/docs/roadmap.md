# Roadmap e Cronograma - PoC Fintech

O desenvolvimento será dividido em fases agressivas para garantir a entrega rápida de valor e validação.

## Cronograma Geral: 4 a 7 Semanas

### Fase 0: Definição & Fundação (1 Semana)
**Objetivo:** Alinhamento e Setup Inicial.
- [ ] Definição e "congelamento" dos requisitos de hipótese.
- [ ] Decisão sobre provedores externos (KYC, Pagamentos).
- [ ] Setup do repositório (Monorepo vs Polyrepo).
- [ ] Configuração do ambiente Docker local (Ory, Postgres, APISIX).
- [ ] Design dos contratos de API (Swagger/OpenAPI).

### Fase 1: Desenvolvimento Core (2 a 4 Semanas)
**Objetivo:** Funcionalidades críticas "ponta a ponta".

#### Semana 1-2 (Backend Heavy)
- [ ] Implementação do Auth Flow (Integração Ory Kratos).
- [ ] Estrutura do Ledger no PostgreSQL.
- [ ] Endpoints básicos: Saldo, Extrato.
- [ ] Serviço de Mock de Pagamentos (Top-up/Transferência).

#### Semana 3-4 (Frontend Integration)
- [ ] Telas de Login/Registro (Next.js + Ory UI).
- [ ] Dashboard do usuário.
- [ ] Integração dos formulários com APIs de Backend.
- [ ] Fluxo de KYC básico no frontend.

### Fase 2: Estabilização & Deploy (1 a 2 Semanas)
**Objetivo:** Preparar para uso real (Friends & Family).
- [ ] Testes de carga leves e Pentest básico.
- [ ] Configuração de Observabilidade (Logs, Métricas).
- [ ] Ajustes finais de UX/UI.
- [ ] Deploy em ambiente de Staging/Prod (Cloud Provider).
- [ ] **Go-Live do PoC.**

### Critérios de Sucesso para Encerramento
1. **Onboarding:** X usuários reais completaram o cadastro com sucesso.
2. **Performance:** Tempo médio de checkout/transação < Z segundos.
3. **Segurança:** Zero vulnerabilidades críticas conhecidas.
4. **Negócio:** Volume transacionado atingiu o mínimo esperado para validação.
