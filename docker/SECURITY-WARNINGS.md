# Security Warnings - Development Environment

**⚠️ DO NOT USE IN PRODUCTION**

This document lists all hardcoded credentials and security configurations that are acceptable for development but **MUST** be changed before production deployment.

---

## 🔐 Hardcoded Credentials

### PostgreSQL Databases

#### db-core (Main Database)
- **Username**: `postgres`
- **Password**: `postgres`
- **Location**: `docker/docker-compose.yml` line ~30

#### postgres-kratos (Kratos Identity DB)
- **Username**: `kratos`
- **Password**: `secret`
- **Location**: `docker/docker-compose.yml` line ~60

### Ory Kratos

#### Cookie Secret
- **Value**: `changeme-please-use-random`
- **Purpose**: Session cookie encryption
- **Location**: `docker/kratos/kratos.yml`

#### CSRF Secret
- **Value**: `changeme-csrf-random-secret`
- **Purpose**: CSRF token generation
- **Location**: `docker/kratos/kratos.yml`

#### DSN (Database Connection)
- **Value**: `postgres://kratos:secret@postgres-kratos:5432/kratos?sslmode=disable`
- **Location**: `docker/kratos/kratos.yml`

### APISIX

#### Admin Key
- **Value**: `edd1c9f034335f136f87ad84b625c8f1`
- **Purpose**: APISIX Admin API access
- **Location**: `docker/apisix/config.yaml`

### Backend Encryption

#### Encryption Key
- **Value**: `32cff1bed37bfe6f2d9ad7ae6b95534d`
- **Purpose**: Card data encryption (AES-256)
- **Location**: `docker/docker-compose.yml` (ENCRYPTION_KEY env var)

---

## ⚠️ Security Issues

### 1. Development Mode

**Ory Kratos** is running in development mode with:
- Plaintext secrets (changeme*)
- Disabled CSRF protection in some flows
- Insecure cookie settings (sslmode=disable)

### 2. No TLS/SSL

All services communicate over HTTP:
- Backend API: http://backend:8080
- APISIX Gateway: http://127.0.0.1:9080
- Kratos: http://kratos:4433

### 3. PostgreSQL Security

- Root user with simple password
- Unencrypted connections (`sslmode=disable`)
- No connection pooling limits
- No query timeouts

### 4. APISIX Configuration

- Hardcoded admin key in config file
- No rate limiting on admin API
- All routes publicly accessible (except auth)

### 5. Backend Configuration

- Static encryption key
- No key rotation mechanism
- Secrets in environment variables (visible in docker inspect)

---

## ✅ Action Required Before Production

### Immediate (Critical)

1. **Generate Strong Secrets**
   ```bash
   # Generate 32-byte random secret
   openssl rand -base64 32
   ```

2. **Use Secret Management**
   - Migrate to AWS Secrets Manager / Azure Key Vault / HashiCorp Vault
   - Never commit secrets to Git
   - Use environment-specific secret files (.env.production)

3. **Enable TLS/SSL**
   - Obtain SSL certificates (Let's Encrypt)
   - Configure NGINX/APISIX for HTTPS termination
   - Enable PostgreSQL SSL connections

4. **Secure Database**
   - Create dedicated users with minimal privileges
   - Use strong passwords (16+ characters, mixed case, symbols)
   - Enable SSL connections (`sslmode=require`)
   - Configure connection limits

5. **Disable Kratos Dev Mode**
   - Update `kratos.yml` to production settings
   - Enable CSRF protection
   - Use secure cookie settings (httpOnly, secure, sameSite)

### Medium Priority

6. **Implement Key Rotation**
   - Backend encryption key rotation strategy
   - APISIX admin key rotation
   - Kratos secrets rotation

7. **Network Security**
   - Restrict APISIX admin API to localhost only
   - Use Docker networks to isolate services
   - Configure firewall rules (UFW/iptables)

8. **Monitoring & Alerts**
   - Set up failed login attempt monitoring
   - Alert on unauthorized API access
   - Track secret usage/rotation

### Low Priority (Nice to Have)

9. **Compliance**
   - PCI-DSS compliance validation (for card data)
   - LGPD/GDPR compliance audit
   - Regular security penetration testing

10. **Infrastructure**
    - Use managed database services (RDS, Cloud SQL)
    - Enable automatic backups with encryption
    - Implement disaster recovery plan

---

## 🛠️ How to Generate Secure Secrets

### For Kratos Cookies/CSRF
```bash
openssl rand -base64 32 | head -c 32
```

### For Backend Encryption Key
```bash
openssl rand -hex 32
```

### For PostgreSQL Passwords
```bash
openssl rand -base64 24 | tr -d "=+/" | cut -c1-20
```

### For APISIX Admin Key
```bash
openssl rand -hex 32
```

---

## 📚 References

- [OWASP Top 10 - Security Misconfiguration](https://owasp.org/Top10/A05_2021-Security_Misconfiguration/)
- [Ory Kratos Production Checklist](https://www.ory.sh/docs/kratos/guides/production)
- [APISIX Security Best Practices](https://apisix.apache.org/docs/apisix/terminology/secret/)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/auth-pg-hba-conf.html)

---

**Last Updated**: 2026-01-16  
**Review Frequency**: Before each deployment
