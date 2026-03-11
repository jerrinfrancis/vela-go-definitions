# Codebase Structure

**Analysis Date:** 2026-03-11

## Directory Layout

```
vela-go-definitions/
├── cmd/
│   └── register/                    # CLI entry point for definition registry
│       └── main.go
├── components/                      # ComponentDefinition implementations
│   ├── webservice.go                # Long-running scalable workload
│   ├── daemon.go                    # Daemon workload
│   ├── task.go                      # One-off task workload
│   ├── crontask.go                  # Cron-scheduled task workload
│   ├── statefulset.go               # Stateful workload
│   ├── k8s_objects.go               # Raw K8s objects component
│   ├── ref_objects.go               # Referenced K8s objects
│   ├── shared_helpers.go            # Helper functions for components
│   ├── schemas.go                   # Shared schema definitions
│   ├── *_test.go                    # Unit tests (co-located)
│   └── components_suite_test.go     # Test suite setup
├── traits/                          # TraitDefinition implementations (56 definitions)
│   ├── env.go                       # Environment variable trait
│   ├── affinity.go                  # Pod affinity trait
│   ├── scaler.go                    # Resource scaling trait
│   ├── labels.go                    # Pod labels trait
│   ├── sidecar.go                   # Sidecar container trait
│   ├── gateway.go                   # Network gateway trait
│   ├── json_patch.go                # JSON patch trait
│   ├── *_test.go                    # Unit tests (co-located)
│   └── traits_suite_test.go         # Test suite setup
├── policies/                        # PolicyDefinition implementations (12 definitions)
│   ├── override.go                  # Override policy
│   ├── garbage_collect.go           # Garbage collection policy
│   ├── apply_once.go                # Apply once policy
│   ├── read_only.go                 # Read-only policy
│   ├── common.go                    # Shared policy helpers
│   ├── *_test.go                    # Unit tests (co-located)
│   └── policies_suite_test.go       # Test suite setup
├── workflowsteps/                   # WorkflowStepDefinition implementations (35 definitions)
│   ├── apply_component.go           # Apply component step
│   ├── apply_deployment.go          # Apply deployment step
│   ├── build_push_image.go          # Build and push image step
│   ├── apply_terraform_config.go    # Terraform config step
│   ├── apply_terraform_provider.go  # Terraform provider step
│   ├── collect_service_endpoints.go # Collect endpoints step
│   ├── check_metrics.go             # Metrics check step
│   └── *_test.go                    # Unit tests (co-located)
├── test/
│   ├── e2e/                         # Ginkgo E2E test suite
│   │   ├── e2e_suite_test.go        # Suite setup and initialization
│   │   ├── helpers_test.go          # Shared E2E test utilities
│   │   ├── component_e2e_test.go    # Dynamic component E2E tests
│   │   ├── trait_e2e_test.go        # Dynamic trait E2E tests
│   │   ├── policy_e2e_test.go       # Dynamic policy E2E tests
│   │   └── workflowstep_e2e_test.go # Dynamic workflow step E2E tests
│   └── builtin-definition-example/  # YAML fixtures for E2E tests
│       ├── components/              # Application YAML for component tests
│       │   ├── webservice.yaml
│       │   ├── daemon.yaml
│       │   ├── task.yaml
│       │   ├── crontask.yaml
│       │   ├── statefulset.yaml
│       │   ├── k8s-objects.yaml
│       │   └── ref-objects.yaml
│       ├── traits/                  # Application YAML for trait tests
│       ├── policies/                # Application YAML for policy tests
│       └── workflowsteps/           # Application YAML for workflow step tests
├── module.yaml                      # KubeVela DefinitionModule manifest
├── go.mod                           # Module definition with kubevela fork
├── go.sum                           # Dependency checksums
├── Makefile                         # Build commands
├── README.md                        # Project documentation
└── CLAUDE.md                        # Development guidance
```

## Directory Purposes

**cmd/register/:**
- Purpose: CLI entry point for JSON definition registry export
- Contains: `main.go` with blank imports of all definition packages
- Key files: `main.go`
- Entry point for: `go run ./cmd/register` command

**components/:**
- Purpose: Component definitions (containerized workload types)
- Contains: Constructor functions and templates, co-located unit tests
- Key files: `webservice.go`, `daemon.go`, `task.go`, `statefulset.go`, `crontask.go`, `k8s_objects.go`, `ref_objects.go`, `shared_helpers.go`
- Patterns: Each file has matching `*_test.go` with Ginkgo tests

**traits/:**
- Purpose: Trait definitions (workload enhancement/patching)
- Contains: 56+ trait definitions with container mutation support
- Key files: `env.go`, `affinity.go`, `scaler.go`, `labels.go`, `sidecar.go`, `gateway.go`
- Patterns: `PatchContainer` config for standard mutations, custom CUE blocks for complex logic

**policies/:**
- Purpose: Policy definitions (deployment-time configuration)
- Contains: 12 policy definitions (override, garbage-collect, apply-once, read-only)
- Key files: `override.go`, `garbage_collect.go`, `apply_once.go`, `read_only.go`
- Patterns: Helper types for structured policy parameters

**workflowsteps/:**
- Purpose: Workflow step definitions (delivery pipeline stages)
- Contains: 35+ step definitions for application deployment
- Key files: `apply_component.go`, `build_push_image.go`, `apply_terraform_config.go`
- Patterns: Lightweight parameter definitions with CUE templates

**test/e2e/:**
- Purpose: Ginkgo v2 E2E test suite with dynamic test generation
- Contains: Test helpers, K8s client initialization, YAML fixture loading
- Key files: `e2e_suite_test.go`, `helpers_test.go`, `component_e2e_test.go`
- Patterns: Each test type (component, trait, policy, workflowstep) has dedicated test file with dynamic test generation from YAML fixtures

**test/builtin-definition-example/:**
- Purpose: YAML Application fixtures for E2E tests
- Contains: Test Applications demonstrating each definition type
- Structure: Organized by definition type (components/, traits/, policies/, workflowsteps/)
- Usage: Dynamically discovered and run in parallel by E2E tests

## Key File Locations

**Entry Points:**
- `cmd/register/main.go`: CLI entry point — outputs JSON registry
- `test/e2e/e2e_suite_test.go`: E2E test suite initialization
- Components: `components/<name>.go` (e.g., `components/webservice.go`)
- Traits: `traits/<name>.go` (e.g., `traits/env.go`)
- Policies: `policies/<name>.go` (e.g., `policies/override.go`)
- WorkflowSteps: `workflowsteps/<name>.go` (e.g., `workflowsteps/apply_component.go`)

**Configuration:**
- `module.yaml`: DefinitionModule manifest (maintainers, categories, placement)
- `go.mod`: Module definition and kubevela fork reference
- `Makefile`: Build targets (test-unit, test-e2e, tidy, cleanup)

**Core Logic:**
- Component templates: `components/<name>.go` templateFunc (private function)
- Trait patches: `traits/<name>.go` with `PatchContainerConfig`
- Policy logic: `policies/<name>.go` with helper types
- WorkflowStep definitions: `workflowsteps/<name>.go` minimal parameter setup

**Testing:**
- Unit tests: `<dir>/<name>_test.go` (co-located with definition)
- E2E tests: `test/e2e/<type>_e2e_test.go`
- Fixtures: `test/builtin-definition-example/<type>/`

## Naming Conventions

**Files:**
- Definition files: `<snake_case_name>.go` (e.g., `webservice.go`, `apply_component.go`)
- Test files: `<snake_case_name>_test.go` (e.g., `webservice_test.go`)
- Helper files: `<purpose>_helper.go` or `shared_helpers.go`
- Suite files: `<type>_suite_test.go` (e.g., `components_suite_test.go`)

**Directories:**
- Plural package names: `components/`, `traits/`, `policies/`, `workflowsteps/`
- Fixture organization: `test/builtin-definition-example/<type>/`

**Functions:**
- Constructor functions: PascalCase, exported (e.g., `Webservice()`, `Env()`, `Override()`)
- Template functions: camelCase, private (e.g., `webserviceTemplate()`, `envTemplate()`)
- Helper functions: camelCase (e.g., `readAppFromFile()`, `sanitizeForNamespace()`)
- Init functions: `init()` - triggers registration

**Types:**
- Definition types (from defkit): `*defkit.ComponentDefinition`, `*defkit.TraitDefinition`, `*defkit.PolicyDefinition`, `*defkit.WorkflowStepDefinition`
- Custom types: PascalCase (e.g., `PatchContainerConfig`)

## Where to Add New Code

**New Component:**
- File: `components/<new_component>.go`
- Tests: `components/<new_component>_test.go`
- E2E Fixture: `test/builtin-definition-example/components/<new_component>.yaml`
- Pattern: Copy `components/webservice.go` structure (constructor + template + init)

**New Trait:**
- File: `traits/<new_trait>.go`
- Tests: `traits/<new_trait>_test.go`
- E2E Fixture: `test/builtin-definition-example/traits/<new_trait>.yaml`
- Pattern: If modifying containers, use `PatchContainerConfig`; otherwise, use direct template

**New Policy:**
- File: `policies/<new_policy>.go`
- Tests: `policies/<new_policy>_test.go`
- E2E Fixture: `test/builtin-definition-example/policies/<new_policy>.yaml`
- Pattern: Define helper types for complex parameters, use `defkit.NewPolicy()`

**New WorkflowStep:**
- File: `workflowsteps/<new_step>.go`
- Tests: `workflowsteps/<new_step>_test.go` (optional for simple definitions)
- E2E Fixture: `test/builtin-definition-example/workflowsteps/<new_step>.yaml`
- Pattern: Minimal parameter setup with `Scope("Application")` or `Scope("Component")`

**New Utility/Helper:**
- File: `<package>/helpers.go` or `<package>/shared_<purpose>.go`
- Usage: Reference from definition files via package-level function calls

## Special Directories

**cmd/register/:**
- Purpose: Single-file entry point for registry export
- Generated: No
- Committed: Yes
- Role: Must import all definition packages (blank imports)

**.archive/:**
- Purpose: Historical versions and migration artifacts
- Generated: No
- Committed: Yes
- Role: Reference for deprecated patterns (not active)

**.claude/:**
- Purpose: Claude Code customizations and agent configurations
- Generated: No
- Committed: Yes
- Role: Development guidance for Claude

**.devcontainer/:**
- Purpose: VS Code Dev Container configuration
- Generated: No
- Committed: Yes
- Role: Standardized development environment

**test/builtin-definition-example/:**
- Purpose: Shared E2E test fixtures (Application YAML)
- Generated: No
- Committed: Yes
- Role: Dynamic source for E2E test generation (discovered at runtime)

---

*Structure analysis: 2026-03-11*
