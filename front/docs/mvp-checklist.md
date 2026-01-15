# Checklist Mínimo para MVP de Fintech (PoC)

Este documento define o escopo estrito para a Prova de Conceito.

## 1. Objetivo do PoC
- **Hipótese Principal:** Validar aquisição de clientes, fluidez do fluxo de pagamento e taxas de retenção inicial.
- **Metas:** Provar que usuários conseguem depositar, transferir e sacar dinheiro com fricção mínima e segurança adequada.

## 2. Requisitos Regulatórios & Compliance
- [ ] **Regime:** Definir modelo (ex: EMI - Electronic Money Institution ou parceria BaaS).
- [ ] **KYC/AML:** Implementação via Ory (Docker) + Integração de provedor de identidade simples.
- [ ] **Legal:** Política de privacidade e Termos de Uso redigidos e acessíveis no app.

## 3. Funcionalidades Core (Obrigatórias)
- [ ] **Cadastro (Onboarding):** Fluxo simplificado com verificação KYC.
- [ ] **Autenticação Segura:** Login com suporte a 2FA (Ory Kratos).
- [ ] **Wallet:** Visualização de saldo e status da conta.
- [ ] **Top-up:** Entrada de fundos via cartão de débito (Gateway simulado ou sandbox).
- [ ] **Transações:** 
  - Pagamentos P2P (entre usuários da plataforma).
  - Transferência para contas externas (Simulação de trilhos bancários).
- [ ] **Histórico:** Feed de transações com detalhes básicos.
- [ ] **Admin Dash:** Painel básico para monitoramento e gestão de usuários/logs.

## 4. Integrações Essenciais
- [ ] **Identity Proivder:** Ory (Kratos/Hydra) + PostgreSQL.
- [ ] **Gateway de Pagamentos:** Microserviço dedicado (interface para Stripe/Adyen ou Mock).
- [ ] **Core Bancário:** Integração com trilhos bancários (Legado ou BaaS).
- [ ] **Notificações:** Serviço de email transacional (ex: SendGrid, AWS SES ou SMTP gratuito).

## 5. Segurança & Risco
- [ ] **Criptografia:** TLS em trânsito (HTTPS) e criptografia em repouso para dados sensíveis.
- [ ] **Gestão de Segredos:** Uso de KMS ou Vault para chaves de API e banco.
- [ ] **Proteção:** Rate limiting no Gateway (APISIX) e regras básicas de detecção de fraude.
- [ ] **Segurança Ofensiva:** Scan de dependências e pentest básico automatizado.

## 6. Dados & Métricas (KPIs de Sucesso)
- [ ] **Negócio:** CAC, ARPU, Conversão no Onboarding, Volume TPV.
- [ ] **Técnico:** Taxa de erros (Error Rate), Latência de API, Uptime.
- [ ] **Auditoria:** Logs imutáveis de transações financeiras e alterações de conta.

## 7. Equipe Mínima
- 1 Product Manager (PM)
- 1 Desenvolvedor Backend (Go)
- 1 Desenvolvedor Frontend (Next.js/React)
- 1 DevOps/Infraestrutura
