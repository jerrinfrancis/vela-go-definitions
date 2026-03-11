# External Integrations

**Analysis Date:** 2026-03-11

## APIs & External Services

**Kubernetes API:**
- Kubernetes cluster API (v1.31.10)
  - SDK/Client: `k8s.io/client-go` v0.31.10
  - Purpose: Create, update, and monitor Kubernetes resources (Deployments, StatefulSets, Jobs, etc.)
  - Auth: KUBECONFIG environment variable points to cluster credentials

**KubeVela Platform:**
- KubeVela Application API (core.oam.dev/v1beta1)
  - SDK/Client: `sigs.k8s.io/controller-runtime` v0.19.7
  - Purpose: Create Application CRDs, query definition status, manage application lifecycle
  - Auth: KUBECONFIG (same as Kubernetes cluster auth)
  - Custom fork: guidewire-oss/kubevela v0.0.0-20260310070415-25c33d481369

**Vela CLI:**
- KubeVela Command-line Interface
  - Built from: guidewire-oss/kubevela fork
  - Purpose: Install definitions (`vela def apply-module`), list definitions (`vela def list-module`)
  - Downloaded or built from source in CI workflows
  - Installation: `.github/actions/setup-vela-environment/action.yaml` installs from kubevela.io/script/install.sh

## Data Storage

**Databases:**
- Not applicable - this is a definition module with no persistent storage

**File Storage:**
- Local filesystem only
  - YAML fixtures: `test/builtin-definition-example/` (components, traits, policies, workflowsteps)
  - Source code: `components/`, `traits/`, `policies/`, `workflowsteps/` (Go code)
  - Module manifest: `module.yaml` (KubeVela DefinitionModule metadata)

**Caching:**
- None - definitions are compiled at build time via defkit.ToJSON()
- k3d cluster uses ephemeral storage for test runs (cleaned up after tests)

## Authentication & Identity

**Auth Provider:**
- Custom - KUBECONFIG-based (Kubernetes service account or user credentials)
  - Implementation: Standard Kubernetes client-go auth flow
  - Location: KUBECONFIG env var → cluster credentials
  - No external OAuth/OIDC integration

**Service Accounts:**
- Kubernetes default service account in vela-system namespace
  - Used by KubeVela controller to reconcile Applications
  - Created by KubeVela installer during `vela install`

## Monitoring & Observability

**Error Tracking:**
- None - errors logged to console/test output
- GitHub Actions: Errors surface in workflow step logs

**Logs:**
- console output (stdout/stderr)
  - Unit tests: go test verbose output
  - E2E tests: Ginkgo output + kubectl logs from pods
  - CI: GitHub Actions step logs + workflow summary

**Metrics:**
- Not applicable at module level
- KubeVela platform provides Prometheus metrics (indirect dependency)
- Test output: Application readiness polling with configurable timeout

## CI/CD & Deployment

**Hosting:**
- GitHub (github.com/oam-dev/vela-go-definitions)
- Fork reference: github.com/guidewire-oss/kubevela (KubeVela custom fork)

**CI Pipeline:**
- GitHub Actions (`.github/workflows/`)
  - Triggers: Push to main, PRs to main, manual workflow_dispatch
  - Runners: ubuntu-latest
  - Tests run in ephemeral k3d clusters

**Definition Installation:**
- Mechanism: `vela def apply-module . --conflict=overwrite`
  - Applied to: Running KubeVela cluster (requires KUBECONFIG)
  - Conflict strategy: Overwrites existing definitions with same name
- Source: ./cmd/register/main.go generates JSON output of all definitions via defkit.ToJSON()

**Build Pipeline:**
- Unit tests (no cluster): `go test ./components/... ./traits/... ./policies/... ./workflowsteps/...`
- E2E tests (with cluster): Ginkgo parallel test runner against live KubeVela cluster
  - Component tests: `.github/workflows/test-definitions.yaml` → `make test-e2e-components`
  - Trait tests: `make test-e2e-traits`
  - Policy tests: `make test-e2e-policies`
  - WorkflowStep tests: `make test-e2e-workflowsteps`

## Environment Configuration

**Required env vars:**
- KUBECONFIG - Path to kubeconfig file (for E2E tests)
  - Example: `~/.kube/master.yaml` or `$KUBECONFIG` in CI
  - Must point to KubeVela-enabled cluster

**Optional env vars:**
- TESTDATA_PATH - Path to test fixture YAML directory (default: `test/builtin-definition-example`)
- E2E_TIMEOUT - Timeout for E2E test suites (default: 10m, customizable)
- PROCS - Parallel process count for Ginkgo (default: 10)

**Secrets location:**
- None required - all auth via KUBECONFIG
- GitHub Actions: No secrets stored in code or configs
- Dev environment: KUBECONFIG points to local/remote cluster credentials

## Webhooks & Callbacks

**Incoming:**
- None - module does not expose HTTP endpoints

**Outgoing:**
- None - module does not make outbound API calls
- KubeVela controller (internal) handles Application reconciliation
- No external webhook notifications triggered

## Cluster Dependencies

**Required Cluster Services:**
- KubeVela core controller - vela-system namespace
  - Handles Application CRD reconciliation
  - Health checks via `kubectl wait --for=condition=available deployment/kubevela-vela-core`

**Multicluster Support:**
- Optional: cluster-gateway (from oam-dev/cluster-gateway v1.9.2)
  - Used by KubeVela for multicluster workload distribution
  - Not required for single-cluster definition testing

## Test Data & Fixtures

**Location:**
- `test/builtin-definition-example/components/` - Component definition test apps (webservice, worker, task, statefulset, etc.)
- `test/builtin-definition-example/trait/` - Trait definition test apps (env, affinity, scaler, hpa, etc.)
- `test/builtin-definition-example/policies/` - Policy definition test apps (override, garbage-collect, topology, etc.)
- `test/builtin-definition-example/workflowsteps/` - WorkflowStep definition test apps (deploy, apply-terraform, webhook, etc.)

**Format:**
- YAML Application manifests
- Each file tests a specific definition with realistic parameters
- Dynamically loaded and applied during E2E test runs
- Test execution: Parallel application submission via kubectl apply

---

*Integration audit: 2026-03-11*
