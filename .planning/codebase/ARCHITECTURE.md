# Architecture

**Analysis Date:** 2026-03-11

## Pattern Overview

**Overall:** Definition Registry Pattern (Self-Registering Module)

**Key Characteristics:**
- Package-based definition modules that auto-register via `init()` functions
- Fluent API builder pattern for declarative definition construction
- Plugin discovery through blank imports into central registry
- Separation of definition builders (Go) from execution templates (CUE)

## Layers

**Definition Construction Layer:**
- Purpose: Define Kubernetes abstractions (components, traits, policies, workflow steps) using fluent API
- Location: `components/`, `traits/`, `policies/`, `workflowsteps/`
- Contains: Definition constructor functions using `defkit.NewComponent()`, `defkit.NewTrait()`, `defkit.NewPolicy()`, `defkit.NewWorkflowStep()`
- Depends on: `defkit` package from kubevela/pkg/definition
- Used by: Registry entry point and E2E tests

**Template Layer:**
- Purpose: CUE template functions that define actual Kubernetes resource behavior
- Location: Inline in definition constructors (e.g., `webserviceTemplate()` function)
- Contains: `*defkit.Template` functions that build runtime manifests using CUE
- Depends on: KubeVela context (`vela := defkit.VelaCtx()`) and resource builders
- Used by: KubeVela controller at runtime

**Parameter Definition Layer:**
- Purpose: Schema definition for user-facing parameters
- Location: Parameter definitions within constructors (e.g., `defkit.String("image")`)
- Contains: `defkit.String`, `defkit.Int`, `defkit.Bool`, `defkit.Enum`, `defkit.Array`, `defkit.Object`, `defkit.StringKeyMap`
- Depends on: Parameter type builders and validation rules
- Used by: Schema validation and Application parsing

**Registry Layer:**
- Purpose: Central aggregation point for all definitions
- Location: `cmd/register/main.go`
- Contains: Blank imports of all definition packages to trigger `init()` registration
- Depends on: `defkit.ToJSON()` for serialization
- Used by: CLI tool to output definition catalog

**Testing Layer:**
- Purpose: Unit and E2E validation of definitions
- Location: `test/` (E2E) and `*_test.go` alongside definitions (unit)
- Contains: Ginkgo BDD tests with dynamic YAML fixture-based test generation
- Depends on: `controller-runtime` client, Kubernetes cluster, KubeVela CRDs
- Used by: CI/CD and local validation

## Data Flow

**Definition Registration Flow:**

1. Definition constructor called (e.g., `Webservice()`)
2. Fluent API builds parameter schema using `defkit.String()`, `defkit.Array()`, etc.
3. Template function assigned via `Template(func(tpl *defkit.Template) {...})`
4. Definition returned as `*defkit.ComponentDefinition` (or Trait/Policy/WorkflowStep)
5. `init()` function calls `defkit.Register(def)` to add to global registry
6. Package import triggers `init()` automatically

**Application Deployment Flow:**

1. User creates Application CRD with component/trait references
2. KubeVela controller parses Application properties against definition schemas
3. Controller invokes template function for each definition
4. Template function receives `defkit.VelaCtx()` with application context
5. Template builds Kubernetes resources using CUE and outputs via `tpl.Output(resource)`
6. Controller applies generated manifests to cluster

**State Management:**

- Definition state: Immutable, defined at registration time via `defkit` builders
- Application state: Mutable KubeVela Application CRD with component/trait specs
- Runtime state: Generated Kubernetes resources (Deployments, StatefulSets, etc.) managed by controllers

## Key Abstractions

**ComponentDefinition:**
- Purpose: Describes containerized workload archetypes (Deployment, StatefulSet, Job, CronJob)
- Examples: `webservice`, `daemon`, `task`, `crontask`, `statefulset`
- Pattern: Constructor function returns `*defkit.ComponentDefinition` with workload type and template

**TraitDefinition:**
- Purpose: Describes capabilities that can patch or enhance components
- Examples: `env`, `affinity`, `scaler`, `labels`, `sidecar`, `gateway`
- Pattern: Uses `PatchContainer` config for mutations or direct CUE patches

**PolicyDefinition:**
- Purpose: Describes policies applied during deployment
- Examples: `override`, `garbage-collect`, `apply-once`, `read-only`
- Pattern: No workload coupling; applied globally to applications

**WorkflowStepDefinition:**
- Purpose: Describes steps in application delivery workflows
- Examples: `apply-component`, `deploy`, `build-push-image`, `apply-terraform-config`
- Pattern: Scoped to application lifecycle (`Scope("Application")`)

## Entry Points

**CLI Entry Point:**
- Location: `cmd/register/main.go`
- Triggers: `go run ./cmd/register`
- Responsibilities:
  - Imports all definition packages (blank imports trigger `init()`)
  - Calls `defkit.ToJSON()` to serialize registry
  - Outputs JSON to stdout for external consumption

**Definition Constructors:**
- Location: Package-level functions in `components/`, `traits/`, `policies/`, `workflowsteps/`
- Examples: `Webservice()`, `Env()`, `Override()`, `ApplyComponent()`
- Responsibilities:
  - Build definition schema
  - Configure parameters
  - Assign template function
  - Return definition object

**Template Functions:**
- Location: Private functions within definition constructors (e.g., `webserviceTemplate()`)
- Triggered: During KubeVela reconciliation
- Responsibilities:
  - Receive application context via `defkit.VelaCtx()`
  - Access user parameters via context
  - Construct Kubernetes resources
  - Output via `tpl.Output(resource)`

## Error Handling

**Strategy:** Schema validation at definition registration time; runtime validation at Application deployment

**Patterns:**
- Parameter validation: `defkit.String("param").Required()` enforces required fields during schema validation
- Type enforcement: `defkit.Enum("type").Values("A", "B")` restricts values at schema level
- Template errors: CUE errors in template functions fail Application reconciliation with clear messages
- Container mutation errors: Trait `PatchContainer` logic includes conditional error checks in CUE (e.g., `if len(_matchContainers_) == 0 { err: "container not found" }`)

## Cross-Cutting Concerns

**Logging:** Not explicitly handled in definition constructors; delegated to KubeVela controller at runtime

**Validation:** Two-stage approach:
- Schema-level: Parameter type and required validation via `defkit` builders
- Application-level: YAML validation against generated CUE schema during Application admission

**Authentication:** Handled by KubeVela controller; definitions are cluster-aware (e.g., cluster parameters in workflow steps)

**Multi-container Support:** Traits support multi-container patching via `containerName` parameter and loop constructs in CUE templates

---

*Architecture analysis: 2026-03-11*
