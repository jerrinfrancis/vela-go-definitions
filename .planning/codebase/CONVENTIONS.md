# Coding Conventions

**Analysis Date:** 2026-03-11

## Naming Patterns

**Files:**
- Definition files: `{name}.go` (e.g., `webservice.go`, `env.go`, `apply_once.go`)
- Test files: `{name}_test.go` in the same package (e.g., `webservice_test.go`)
- Test suites: `{type}_suite_test.go` for Ginkgo suite setup (e.g., `components_suite_test.go`, `e2e_suite_test.go`)
- Helper files: `shared_helpers.go`, `helpers_test.go`

**Functions:**
- Definition constructors: PascalCase returning specific types (e.g., `Webservice()`, `ApplyOnce()`, `Env()`)
- Each package has an `init()` function that auto-registers definitions via `defkit.Register()`
- Helper functions: camelCase, often simple action verbs (e.g., `readAppFromFile()`, `sanitizeForNamespace()`)
- Private helpers: lowercase prefixed (e.g., `updateAppNamespaceReferences()`)

**Variables:**
- Package-level constants: UPPERCASE for timeout/polling values (e.g., `AppRunningTimeout`, `PollInterval`)
- Package-level vars: camelCase, often initialized to nil (e.g., `k8sClient`)
- Local variables: camelCase (e.g., `app`, `ns`, `comp`)
- Underscores for blank imports: `_ "package/path"` to trigger side effects (registration)

**Types:**
- Structs: PascalCase (e.g., `Application`, `Component`)
- defkit-created types: descriptive names matching CUE parameter names (e.g., `volumeMounts`, `containerPort`)

## Code Style

**Formatting:**
- Standard Go formatting with `gofmt`
- Line length: idiomatic Go conventions (typically 80-120 characters)
- Imports grouped by: standard library, external packages, local packages (separated by blank lines)
- License header: Apache 2.0 standard 14-line header on all source files

**Linting:**
- Governed by Go's standard linters
- No detected linting config files (`.golangci.yml`, `golangci-lint` not configured in CI)
- Unit tests run with `-race` flag to detect race conditions

## Import Organization

**Order:**
1. Standard library packages (`fmt`, `os`, `strings`, etc.)
2. External packages from Go ecosystem (`github.com/onsi/...`, `k8s.io/...`)
3. Local packages from same module (`github.com/oam-dev/vela-go-definitions/...`)

**Path Aliases:**
- Dot imports used only for Ginkgo/Gomega (`. "github.com/onsi/ginkgo/v2"`, `. "github.com/onsi/gomega"`)
- Matcher imports: `. "github.com/oam-dev/kubevela/pkg/definition/defkit/testing/matchers"`
- Blank imports for registration triggers: `_ "github.com/oam-dev/vela-go-definitions/components"`

## Error Handling

**Patterns:**
- Explicit error checks: `if err != nil { ... }`
- Return errors with context: `fmt.Errorf("message: %w", err)` for wrapping
- For E2E tests: `Expect(err).NotTo(HaveOccurred(), "message")` using Gomega assertions
- Silent error ignoring permitted only in test context for side effects (e.g., `_ = k8sClient.Delete(ctx, app)`)
- Kubernetes API error checking: `if !errors.IsAlreadyExists(err)` and `apierrors.IsNotFound(err)` patterns

## Logging

**Framework:** Standard output via `GinkgoWriter` in E2E tests

**Patterns:**
- E2E tests use `GinkgoWriter.Printf()` for test output (e.g., `GinkgoWriter.Printf("Creating namespace %s...\n", uniqueNs)`)
- Test context messages use string interpolation with `\n` for newlines
- GitHub Actions integration: test summary appended to `$GITHUB_STEP_SUMMARY`

## Comments

**When to Comment:**
- Function-level documentation: single-line comments above exported functions explaining purpose
- Complex logic: clarify intent (e.g., "Each app has unique name, so namespace based on app name is unique per test")
- API usage: explain non-obvious defkit patterns (e.g., PatchContainer fluent API usage, CustomPatchContainerBlock content)
- CUE code blocks: explain complex conditionals or transformations

**JSDoc/TSDoc:**
- Not used (Go uses single-line and block comments)
- Exported functions have single-line summary comments (e.g., `// Webservice creates a webservice component definition.`)
- Multi-line descriptions use adjacent comment blocks

## Function Design

**Size:**
- Definition constructors: 10-50 lines typically (defkit fluent API chains)
- Helper functions: 5-30 lines (utility functions)
- E2E test blocks: 20-100 lines (Ginkgo test bodies)

**Parameters:**
- defkit functions take no parameters — chain configuration fluently
- Helper functions accept `tpl *defkit.Template` as context (e.g., `ContainerMountsHelper(tpl, volumeMounts)`)
- Test helpers often return `*defkit.HelperVar` for template composition

**Return Values:**
- Definition functions return typed pointers: `*defkit.ComponentDefinition`, `*defkit.TraitDefinition`, `*defkit.PolicyDefinition`
- Helper functions return `*defkit.HelperVar` for chainable composition
- E2E helpers return tuples: `([]string, error)`, `(*v1beta1.Application, error)`
- Test blocks return nothing (use Gomega `Expect()` for assertions)

## Module Design

**Exports:**
- Definition constructors are exported (PascalCase): `Webservice()`, `ApplyOnce()`, `Env()`
- Package-level registration happens in unexported `init()` (runs automatically on import)
- Helper functions exported if reusable across packages: `ContainerMountsHelper()`, `PodVolumesHelper()`

**Barrel Files:**
- No explicit barrel files (index.go pattern not used)
- Import entire packages to trigger all `init()` registrations (e.g., `_ "github.com/oam-dev/vela-go-definitions/components"`)
- `cmd/register/main.go` imports all definition packages by blank import

**defkit API Patterns:**

All definitions follow a consistent fluent API pattern:

```go
func Webservice() *defkit.ComponentDefinition {
    param := defkit.String("name").Required().Description("...")

    return defkit.NewComponent("webservice").
        Description("...").
        Workload("apps/v1", "Deployment").
        Params(param, ...).
        Template(func(tpl *defkit.Template) {
            // Template implementation
        })
}
```

**Parameter Types:**
- Scalar: `defkit.String()`, `defkit.Int()`, `defkit.Bool()`
- Collections: `defkit.Array()`, `defkit.List()`, `defkit.StringList()`, `defkit.StringKeyMap()`
- Complex: `defkit.Object()`, `defkit.Struct()`
- Enums: Use `.Values()` not `.Enum()` for consistency (recent defkit refactor)

**Method Chains:**
- Parameter builders: `.Required()`, `.Default()`, `.Description()`, `.Short()`, `.Ignore()`, `.Optional()`
- Definition builders: `.Description()`, `.Workload()`, `.AppliesTo()`, `.Template()`, `.Helper()`, `.Params()`
- Template builders: `.Helper()`, `.Output()`, `.Outputs()`, `.UsePatchContainer()`

---

*Convention analysis: 2026-03-11*
