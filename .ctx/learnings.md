# Learnings — SecurePay

## Format
Add here when a bug is resolved or something important is learned:
- Date
- What happened
- How it was resolved
- How to prevent it next time

---

## [2026-02-15] Initial
No learning records yet.

## [2026-02-15] Task Status Synchronization
It was noticed that files for a task marked "completed" in Tasks.json ("id: 1") were not on disk.
There might be inconsistency between task status and file system, always check the file system.
Solution: Files were recreated.

## [2026-02-15] Missing Proto Files
During API Gateway gRPC client implementation, it was noticed that proto files or generated codes for backend services were missing.
Solution: `PaymentServiceClient` interface was mocked temporarily. Real generated code should be used after backend services (`payment-service`, `account-service`) proto definitions are made.
In conventions, `proto/gen/go/` folder is expected but not present. But task 3.1 addressed creating proto files.

## [2026-02-15] Protoc Compiler Missing and Solution
Proto files were created under Task 3.1 but `protoc` was initially thought to be missing.
Later, plugins were installed with `go install` and `protoc` was detected on the system (v33.5).
`proto` folder was made a separate Go module and linked to `api-gateway` module with `go.work`.
Generated codes were successfully created and `grpc_clients.go` was updated.

## [2026-02-15] Minikube Docker Driver (WSL)
It was reiterated that Minikube should be started with `docker` driver on WSL. `virtualbox` or other drivers may cause issues on WSL.
Solution: `minikube start --driver=docker` command was used.

## [2026-02-15] PowerShell Command Chaining
`&&` operator does not work in PowerShell environment.
Solution: Execute commands separately or use `;` (or appropriate operator based on powershell version).

## [2026-02-15] Minikube and Docker Credential Helper Issue (WSL)
When `eval $(minikube docker-env)` is used in WSL environment, Docker build process gives `docker-credential-desktop.exe: exec format error`. Because even though Minikube environment is Linux-based, it tries to call Windows credential helper.
Solution: Build process can be done by pointing to an empty config file without credential helper using `DOCKER_CONFIG` environment variable. Alternatively, image tag can be changed to bypass cache.

## [2026-02-15] Minikube Image Caching
Minikube struggles to update the `latest` tag for images built locally with `imagePullPolicy: Never`. Old image ID might be used even if Pod is restarted.
Solution: Changing image tag (e.g., `v1.0.0`) is the most definitive solution.

## [2026-02-15] SPIFFE Socket Mount Method
Using `csi.spiffe.io` driver in Kubernetes environment can cause instability or path issues in some cases.
Solution: Mounting `/run/spire/agent-sockets` directory using `HostPath` volume is a more stable and reliable method.

## [2026-02-15] Service Dependencies and Startup
Ability to start services like API Gateway even if dependent backend services (Payment, Account) are not up yet facilitates development and testing processes.
Solution: Service connection (dial) errors in `main.go` were demoted from `Fatal` to `Warning` level, preventing app crash and making `/health` endpoint accessible.

## [2026-02-16] Docker Credential Helper Issue in Helm Installation (WSL)
During Helm chart installation (`helm install`), Docker credential helper error (`exec format error`) can also be received. Because Helm uses Docker configuration for local chart cache operations.
Solution: Credential helper can be disabled using `DOCKER_CONFIG` environment variable.

## [2026-02-16] PostgreSQL Connection Test
Using a temporary pod (`kubectl run --rm`) is practical to test database connections inside Kubernetes.
Example: `kubectl run test-db --rm -i --image=postgres:alpine --env="PGPASSWORD=password" -- sh -c 'psql -h my-postgres-postgresql -U securepay -d securepay -c "SELECT 1"'`

## [2026-02-16] PostgreSQL Connection and Migration
Using `kubectl cp` to copy files and then executing with `exec` is the most stable method to avoid symbolic link and path errors.
Also, PostgreSQL password installed with Helm is usually stored in a secret as base64 encoded.
Secret name is in `<release-name>-postgresql` format. Key `password` holds the specific user password, `postgres-password` holds the admin password.
Database name and username are determined during Helm installation, defaults might be `postgres` but if changed with `--set` (e.g., `securepay`), this should be noted.
Command: `kubectl get secret my-postgres-postgresql -o jsonpath="{.data.password}" | base64 -d`
