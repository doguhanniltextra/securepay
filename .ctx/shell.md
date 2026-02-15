# SHELL PROTOCOL (WSL)

### 1. COMMAND CHAINING
**FORBIDDEN:** Using `&&` operator.
**MANDATORY:** Writing each command on a separate line.

```powershell
# FORBIDDEN
wsl kubectl get pods && wsl kubectl get svc

# MANDATORY
wsl kubectl get pods
wsl kubectl get svc
```

### 2. VISUALIZATION
**FORBIDDEN:** Using filtering tools like `grep`, `awk`, `sed`.
**MANDATORY:** Using `kubectl`'s own formatting options.

```powershell
# FORBIDDEN
wsl kubectl get pods | grep securepay

# MANDATORY
wsl kubectl get pods -n securepay
```
