## Summary

> _What does this PR do? One or two sentences — be concrete._

<!-- Example: Adds OTel-aware structured JSON logging to account-service and payment-service.
     trace_id / span_id are now injected automatically into every log record via slog.InfoContext(ctx, …). -->


## Motivation / Linked Issue

Closes #<!-- issue number -->

> _Why is this change needed? What problem does it solve?_


## Type of Change

<!-- Tick all that apply -->

- [ ]  Bug fix (non-breaking)
- [ ]  New feature (non-breaking)
- [ ]  Breaking change (existing behaviour changes)
- [ ]  Refactor / internal improvement
- [ ]  Dependency update
- [ ]  Documentation / comments only
- [ ]  CI / workflow change
- [ ]  Security improvement


## Services / Components Affected

<!-- Tick every service or layer this PR touches -->

- [ ] `api-gateway`
- [ ] `account-service`
- [ ] `payment-service`
- [ ] `notification-service`
- [ ] `proto` (shared protobuf definitions)
- [ ] `k8s` (Kubernetes manifests)
- [ ] `.github` (CI workflows / templates)
- [ ] Other: <!-- describe -->


## How Was It Tested?

<!-- Describe what you ran / checked. CI alone is not enough — add context. -->

```
# Example
go test -v ./internal/logger/...
grpcurl -plaintext ... PaymentService/InitiatePayment
kubectl logs -n default deploy/payment-service | jq .
```

**Test results:**

- [ ] Unit tests pass (`go test ./...`)
- [ ] `go vet ./...` clean
- [ ] `gofmt` clean
- [ ] Logs verified in JSON with `trace_id` / `span_id` fields populated
- [ ] Tested end-to-end in Minikube (if applicable)


## Observability Checklist

> _Fill this in for any PR that touches logging, tracing, or metrics._

- [ ] All new log calls use `slog.InfoContext(ctx, …)` / `slog.ErrorContext(ctx, …)` (not bare `slog.Info`)
- [ ] No `fmt.Println` or `log.Println` introduced
- [ ] New spans are named `<service>.<layer>.<operation>` (e.g. `account-service.handler.CheckBalance`)
- [ ] Errors are recorded on the span via `span.RecordError(err)` where relevant
- [ ] `trace_id` and `span_id` appear in log output when a span is active (verified manually or via test)

> _Not applicable? Delete this section._


## Security Checklist

- [ ] No secrets, credentials, or PII are logged
- [ ] SPIFFE/SPIRE mTLS is not weakened by this change
- [ ] No new network endpoints are exposed without authentication

> _Not applicable? Delete this section._


## Screenshots / Log Samples

> _Paste a sample JSON log line or a Jaeger trace screenshot if it helps reviewers._

```json
{
  "time": "2026-02-20T01:22:39+03:00",
  "level": "INFO",
  "msg": "Payment initiated successfully",
  "payment_id": "pay-123",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id":  "00f067aa0ba902b7",
  "trace_flags": "01"
}
```

> _Delete this section if not applicable._


## Reviewer Notes

> _Anything specific you want reviewers to focus on?_
> _E.g. "Please double-check the context propagation in kafka/consumer.go line 60."_


## Pre-merge Checklist (for the PR author)

- [ ] PR title follows Conventional Commits: `feat(account-service): …` / `fix(payment-service): …`
- [ ] Branch is up-to-date with `master`
- [ ] CI workflow `PR – Structured Logging` is green
- [ ] At least one reviewer approved
- [ ] `CHANGELOG.md` updated (if the project maintains one)
