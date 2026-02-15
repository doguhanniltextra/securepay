# GITHUB PROTOCOL (CONVENTIONAL COMMITS)

## 1. COMMIT STRUCTURE
**MANDATORY:** Every commit must be in the following format:

```
<type>(<scope>): <subject>

<body>
```

## 2. TYPES

| Type | Description |
|------|-------------|
| **feat** | New feature |
| **fix** | Bug fix |
| **docs** | Documentation |
| **style** | Formatting, punctuation |
| **refactor** | Code restructuring |
| **test** | Adding tests |
| **chore** | Build, CI, tools |

## 3. SCOPES

| Scope | Description |
|-------|-------------|
| **payment-service** | Payment service |
| **account-service** | Account service |
| **notification-service** | Notification service |
| **api-gateway** | API gateway |
| **helm** | Helm installations |
| **minikube** | Minikube configuration |
| **ci** | CI/CD pipeline |

## 4. EXAMPLES

```
feat(payment-service): add idempotency check

- Redis idempotency key check added
- Duplicate processing with same key prevented
```

```
fix(account-service): fix balance calculation

- Optimistic locking added in PostgreSQL transaction
- Race condition fixed
```

```
docs(helm): update README.md

- Installation steps updated
- New parameters added
```
