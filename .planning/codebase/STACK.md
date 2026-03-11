# Technology Stack

**Analysis Date:** 2026-03-11

## Languages

**Primary:**
- Go 1.23.8 - KubeVela Definition Module implementation, component/trait/policy/workflow definitions

## Runtime

**Environment:**
- Go 1.23.8 (from `go.mod`)

**Package Manager:**
- Go modules (go.mod)
- Lockfile: `go.sum` (present)

## Frameworks

**Core:**
- github.com/oam-dev/kubevela (custom fork: guidewire-oss/kubevela) - KubeVela platform and defkit API for definition authoring
  - Replace directive: `github.com/oam-dev/kubevela => github.com/guidewire-oss/kubevela v0.0.0-20260310070415-25c33d481369`
  - defkit package: `github.com/oam-dev/kubevela/pkg/definition/defkit` - fluent API for defining ComponentDefinition, TraitDefinition, PolicyDefinition, WorkflowStepDefinition

**Kubernetes Clients:**
- k8s.io/api v0.31.10 - Kubernetes API types (Deployment, StatefulSet, Job, CronJob, etc.)
- k8s.io/apimachinery v0.31.10 - Kubernetes common types and utilities
- k8s.io/client-go v0.31.10 - Kubernetes client library for API interactions
- sigs.k8s.io/controller-runtime v0.19.7 - KubeVela Application and Definition CRD clients

**Testing:**
- github.com/onsi/ginkgo/v2 v2.23.3 - BDD testing framework for E2E tests
- github.com/onsi/gomega v1.36.2 - Assertion library (used with Ginkgo)

**YAML/Config:**
- sigs.k8s.io/yaml v1.4.0 - YAML marshaling/unmarshaling for KubeVela objects

**Build/Dev:**
- Make-based build (Makefile in root)
- Go CLI tools (go test, go mod, go run)
- Ginkgo CLI (`make install-ginkgo` installs ginkgo binary)

## Key Dependencies

**Critical:**
- github.com/oam-dev/kubevela (custom fork) - Provides defkit API for programmatic definition authoring
- github.com/onsi/ginkgo/v2 v2.23.3 - E2E test runner with parallel execution support
- github.com/onsi/gomega v1.36.2 - Assertions for E2E tests
- k8s.io/client-go v0.31.10 - Kubernetes cluster interaction

**Infrastructure/Observability:**
- github.com/kubevela/pkg v1.9.3-0.20251028181209-ef6824214171 (indirect) - KubeVela shared utilities
- github.com/kubevela/workflow v0.6.3-0.20251125110424-924e73add777 (indirect) - KubeVela workflow engine
- github.com/oam-dev/cluster-gateway v1.9.2-0.20250629203450-2b04dd452b7a (indirect) - Multicluster support
- go.opentelemetry.io/otel v1.28.0 (indirect) - OpenTelemetry tracing support (via KubeVela)
- github.com/prometheus/client_golang v1.20.5 (indirect) - Prometheus metrics (via KubeVela)

**Supporting:**
- github.com/google/uuid v1.6.0 - UUID generation for test resources
- go.etcd.io/etcd/client/v3 v3.5.16 (indirect) - etcd client for cluster config management (via KubeVela)
- google.golang.org/grpc v1.67.1 (indirect) - gRPC for cluster-gateway communication

## Configuration

**Environment:**
- KUBECONFIG - Points to active Kubernetes cluster (required for E2E tests)
  - Default dev setup: `~/.kube/master.yaml`
- TESTDATA_PATH - Path to YAML fixture files for E2E tests
  - Default: `test/builtin-definition-example`
- E2E_TIMEOUT - Timeout for E2E test suites (default: 10m)
- PROCS - Number of parallel Ginkgo processes for E2E tests (default: 10)

**Build:**
- `Makefile` - Primary build configuration
  - `go.mod` - Module dependencies and replace directives
  - `go.sum` - Dependency lock file

## Platform Requirements

**Development:**
- Go 1.23.8+
- kubectl CLI (for K8s cluster interaction)
- k3d or Kind (for local Kubernetes cluster in CI/testing)
- Docker (for k3d cluster creation)
- Ginkgo CLI (installed via `make install-ginkgo`)

**CI/CD:**
- GitHub Actions (`.github/workflows/` configurations)
- Ubuntu Linux (ubuntu-latest runner)
- k3d for ephemeral K3s cluster creation
- Vela CLI (built from KubeVela source or downloaded from kubevela.io)

**Production:**
- Kubernetes v1.31.10+ cluster
- KubeVela v0.0.0-20260310070415-25c33d481369 (from guidewire-oss/kubevela fork)
- Definitions installed via `vela def apply-module` command

## CI/CD Integration

**GitHub Actions:**
- `unit-tests.yaml` - Runs on push to main, PRs to main
  - `go test` across all definition packages
  - Go version read from go.mod

- `test-definitions.yaml` - Runs on workflow_dispatch, go.mod/go.sum changes
  - Parallel E2E test suites (components, traits, policies, workflowsteps)
  - Sets up k3d cluster, installs KubeVela from replace directive source
  - Installs definitions via `vela def apply-module`

**Custom Actions:**
- `.github/actions/setup-vela-environment/action.yaml` - Composite action
  - Sets up Go, k3d cluster
  - Downloads and installs Vela CLI
  - Builds Vela CLI from source (from go.mod replace directive)
  - Installs definitions from module
  - Installs Ginkgo CLI

## Environment Configuration

**Required env vars (for E2E):**
- KUBECONFIG - Path to kubeconfig for target Kubernetes cluster

**Optional env vars:**
- TESTDATA_PATH - Override test data location (default: test/builtin-definition-example)
- E2E_TIMEOUT - Override E2E test timeout (default: 10m)
- PROCS - Override parallel process count (default: 10)

**Secrets location:**
- No external API keys or secrets required
- GitHub Actions: CI workflow runs in github.com/oam-dev/vela-go-definitions or fork
- Dev environment: Requires valid KUBECONFIG pointing to KubeVela cluster

---

*Stack analysis: 2026-03-11*
