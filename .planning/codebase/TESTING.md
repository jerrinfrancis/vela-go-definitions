# Testing Patterns

**Analysis Date:** 2026-03-11

## Test Framework

**Runner:**
- `go test` for unit tests
- Ginkgo v2 (github.com/onsi/ginkgo/v2 v2.23.3) for E2E tests
- Config: Makefile-based test execution with Ginkgo CLI

**Assertion Library:**
- Gomega (github.com/onsi/gomega v1.36.2)
- Custom matchers from defkit: `github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers`

**Run Commands:**
```bash
make test-unit                      # Run all unit tests
make test-e2e                       # Run all E2E tests (requires cluster)
make test-e2e-components            # E2E tests for component definitions
make test-e2e-traits                # E2E tests for trait definitions
make test-e2e-policies              # E2E tests for policy definitions
make test-e2e-workflowsteps         # E2E tests for workflow step definitions
make install-ginkgo                 # Install Ginkgo CLI
make cleanup-e2e-namespaces         # Delete e2e test namespaces
```

**Environment Variables:**
- `TESTDATA_PATH`: Path to YAML test fixtures (default: `test/builtin-definition-example`)
- `E2E_TIMEOUT`: E2E test timeout (default: `10m`)
- `PROCS`: Ginkgo parallelism (default: `10`)

## Test File Organization

**Location:**
- Unit tests: co-located with implementation (`traits/env.go` → `traits/env_test.go` in `traits_test` package)
- E2E tests: centralized in `test/e2e/` directory
- Test fixtures: `test/builtin-definition-example/{type}/` (e.g., `test/builtin-definition-example/components/`)

**Naming:**
- Unit test files: `{definition}_test.go` (e.g., `webservice_test.go`, `labels_test.go`)
- Suite files: `{type}_suite_test.go` (e.g., `components_suite_test.go`)
- E2E test files: `{type}_e2e_test.go` (e.g., `component_e2e_test.go`, `trait_e2e_test.go`)
- Helper files: `helpers_test.go` for shared E2E utilities

**Structure:**
```
components/
├── webservice.go              # Definition
├── webservice_test.go         # Unit test
├── components_suite_test.go   # Suite setup
└── shared_helpers.go          # Shared utility functions

test/e2e/
├── component_e2e_test.go      # Component E2E tests
├── trait_e2e_test.go          # Trait E2E tests
├── policy_e2e_test.go         # Policy E2E tests
├── workflowstep_e2e_test.go   # WorkflowStep E2E tests
├── e2e_suite_test.go          # Ginkgo suite setup
└── helpers_test.go            # E2E helper functions

test/builtin-definition-example/
├── components/
│   ├── webservice.yaml
│   ├── daemon.yaml
│   └── ...
├── traits/
│   ├── env.yaml
│   └── ...
└── ...
```

## Test Structure

**Suite Organization:**

Unit test pattern (Ginkgo v2 with dot imports):
```go
var _ = Describe("Webservice Component", func() {
    Describe("Webservice()", func() {
        It("should create a webservice component definition", func() {
            comp := components.Webservice()
            Expect(comp.GetName()).To(Equal("webservice"))
        })

        It("should have correct workload type", func() {
            comp := components.Webservice()
            workload := comp.GetWorkload()
            Expect(workload.APIVersion()).To(Equal("apps/v1"))
        })
    })
})
```

E2E test pattern with dynamic test generation:
```go
var _ = Describe("Component Definition E2E Tests", Label("components"), func() {
    ctx := context.Background()

    Context("when testing component definitions", func() {
        testDataPath := filepath.Join(getTestDataPath(), "components")

        // Dynamic test generation from YAML files
        When("applying component applications", func() {
            for _, file := range func() []string {
                testPath := filepath.Join(getTestDataPath(), "components")
                f, _ := listYAMLFiles(testPath)
                return f
            }() {
                file := file
                It(fmt.Sprintf("should run %s", filepath.Base(file)), func() {
                    // Test implementation
                })
            }
        })
    })
})
```

**Patterns:**
- Setup: `BeforeSuite()` for one-time initialization (e.g., `initK8sClient()`)
- Setup per-test: `BeforeEach()` for per-test setup (e.g., `policy := policies.ApplyOnce()`)
- Test labels: `Label("components")`, `Label("traits")`, `Label("policies")`, `Label("workflowsteps")` for filtering with `--label-filter`
- Assertions: Gomega `Expect().To()`, `Expect().NotTo()` chains
- Output: `GinkgoWriter.Printf()` for test logging

## Mocking

**Framework:**
- No external mocking library (testify/mock, gomock not used)
- Kubernetes client mock: actual `controller-runtime` client against kubeconfig
- JSON unmarshaling: standard Go `encoding/json` with error checks

**Patterns:**
```go
// E2E: Real K8s client initialization
cfg, err := config.GetConfig()
k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})

// Unit: Direct function calls with test data
comp := components.Webservice()
Expect(comp.GetName()).To(Equal("webservice"))
```

**What to Mock:**
- Kubernetes resources: use real client with E2E test cluster
- File I/O: actual `os.ReadFile()` with test fixture files in `test/builtin-definition-example/`

**What NOT to Mock:**
- Definition builder functions (test them directly)
- defkit API calls (test actual generated output)
- Kubernetes API calls (use live cluster in E2E tests)

## Fixtures and Factories

**Test Data:**
```go
// E2E test: Dynamic YAML loading and patching
app, err := readAppFromFile(file)
app.SetNamespace(uniqueNs)
updateAppNamespaceReferences(app, uniqueNs)

// Unit test: Direct constructor calls
comp := components.Webservice()
tpl := defkit.NewTemplate()
comp.GetTemplate()(tpl)
Expect(tpl.GetOutput()).To(BeDeployment())
```

**Location:**
- YAML fixtures: `test/builtin-definition-example/{type}/` (e.g., `test/builtin-definition-example/components/webservice.yaml`)
- Programmatic fixtures: Generated in-memory via defkit constructors
- Test namespaces: Created dynamically with unique names (`e2e-{sanitized-app-name}`) per test

**Fixture Strategy:**
- Each YAML fixture is a complete KubeVela Application manifest
- Tests read fixture → set namespace → patch namespace references → apply to cluster
- Parallel tests use unique namespaces to avoid conflicts
- Cleanup: `make cleanup-e2e-namespaces` deletes all `e2e-*` namespaces

## Coverage

**Requirements:** No coverage threshold enforced in CI

**View Coverage:**
```bash
go test -cover ./components/... ./traits/... ./policies/... ./workflowsteps/...
```

## Test Types

**Unit Tests:**
- Scope: Definition structure validation, parameter correctness, CUE output validation
- Approach: Direct constructor calls, assertion on `.GetName()`, `.GetDescription()`, `.ToCue()`
- Pattern: One test file per definition (e.g., `webservice_test.go` tests `webservice.go`)
- Examples in: `components/webservice_test.go`, `traits/labels_test.go`, `policies/apply_once_test.go`

Unit test assertions check:
- Definition name and description
- Parameter presence and types (using `.HaveParamNamed()` matcher)
- Generated CUE output contains expected substrings (`.ContainSubstring()`)
- Workload type and API version
- Auxiliary output keys (`.HaveKey()` on Outputs map)

**Integration Tests:**
- Not heavily used (defkit has its own integration test matchers)
- Some pattern matching in unit tests validates generated CUE structure

**E2E Tests:**
- Framework: Ginkgo v2 with dynamic test generation from YAML fixtures
- Cluster requirement: Live KubeVela cluster with definitions pre-registered
- Approach: Load YAML Application → create namespace → apply to cluster → poll for readiness → verify status
- Timeout: 5 minutes per application (configurable via `E2E_TIMEOUT`)
- Polling: 5-second intervals for status checks (constants: `AppRunningTimeout`, `PollInterval`)

E2E test flow (`test/e2e/component_e2e_test.go` pattern):
1. Load YAML Application from `test/builtin-definition-example/components/`
2. Sanitize app name for namespace
3. Create unique namespace `e2e-{app-name}`
4. Delete app if exists (clean slate)
5. Create Application in cluster
6. Poll for `app.Status.Phase == "running"`
7. Assert deployment replicas match expected count
8. Clean up namespace

## Common Patterns

**Async Testing:**
- E2E uses polling with timeout: `Eventually()` with `WithTimeout()` and `WithPolling()`
- Example: wait for Application phase to be "running" within 5 minutes, polling every 5 seconds

**Error Testing:**
- Unit tests: `Expect(trait.ToCue()).To(ContainSubstring("..."))`
- E2E tests: `Expect(err).NotTo(HaveOccurred(), "context message")`
- Explicit error unwrapping: `if err != nil && !errors.IsAlreadyExists(err) { ... }`

**Parallelism:**
- E2E tests parallelized by Ginkgo: `ginkgo --procs=10` (default in Makefile)
- Each test gets isolated namespace (`e2e-{app-name}`) for safe parallel execution
- Dynamic test generation ensures no test name collisions

**Test Helpers:**
- `initK8sClient()`: Initialize controller-runtime client once per suite
- `readAppFromFile()`: Parse multi-doc YAML, extract Application
- `updateAppNamespaceReferences()`: Patch namespace in app spec and nested properties
- `sanitizeForNamespace()`: Convert app name to valid Kubernetes namespace
- `listYAMLFiles()`: List fixture files from test directory

---

*Testing analysis: 2026-03-11*
