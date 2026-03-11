# Codebase Concerns

**Analysis Date:** 2026-03-11

## Tech Debt

**ForkedDependency - KubeVela on Custom Fork:**
- Issue: Project uses a custom fork of KubeVela (`github.com/guidewire-oss/kubevela`) via `go.mod` replace directive instead of upstream
- Files: `go.mod:16`
- Impact: Maintenance burden if upstream KubeVela releases new APIs or fixes critical bugs. Fork divergence increases over time. Integration with other KubeVela tools/plugins may break if they depend on upstream APIs
- Fix approach: Evaluate if custom fork features can be contributed upstream or if a minimal fork wrapper can reduce maintenance. Document reasons for fork and track upstream changes for eventual merge back

**UnstableAPIVersions:**
- Issue: Primary dependency `github.com/oam-dev/kubevela` pinned to `v0.0.0` with replace to pre-release version (`v0.0.0-20260310070415-25c33d481369`)
- Files: `go.mod:6,16`
- Impact: No semantic versioning provides no guarantee of stability. Definitions built against this unstable API may break silently in production when the underlying defkit API changes
- Fix approach: Track upstream KubeVela releases. Create semver tags for this module aligned with KubeVela versions. Test against multiple KubeVela versions in CI

## Test Coverage Gaps

**WorkflowSteps - Minimal Unit Test Coverage:**
- What's not tested: 31 workflow step definition files have only 2 dedicated unit test files (`workflowsteps_test.go` and `workflowsteps_suite_test.go`)
- Files: `workflowsteps/*.go` (31 definitions) vs `workflowsteps/workflowsteps_test.go`
- Risk: Most workflowstep definitions (notification, request, terraform, cloud resources, etc.) only get parameter validation via E2E tests. Bugs in parameter definition, defaults, or CUE generation won't be caught until integration testing
- Coverage: Unit tests cover Deploy, Suspend, ApplyComponent only. Other 28 definitions untested at unit level
- Priority: High

**Traits - Partial Unit Test Coverage:**
- What's not tested: 29 trait definitions with only 25 dedicated test files. ~4 traits have no corresponding unit test
- Files: `traits/*.go` (29 definitions) with sparse unit tests
- Risk: Trait parameter validation and default values not unit-tested
- Priority: Medium

**Components - Partial Unit Test Coverage:**
- What's not tested: 10 component definitions with 7 unit test files. 3 components missing dedicated tests
- Files: `components/*.go` (10 definitions) with 7 test files
- Risk: Parameter validation and schema generation for untested components not verified
- Priority: Medium

**E2E Test Fixture Coverage:**
- What's not tested: 76 YAML fixture files in `test/builtin-definition-example/` cover basic scenarios only. Missing edge cases: invalid parameters, conflicting configurations, permission errors, resource exhaustion
- Files: `test/builtin-definition-example/<type>/` directories
- Risk: Integration tests may pass with valid configs but fail at runtime with malformed parameters or edge cases
- Priority: Medium

## Fragile Areas

**WorkflowSteps - Large Definitions with Complex Parameters:**
- Files: `workflowsteps/notification.go` (409 lines), `workflowsteps/request.go`, `workflowsteps/build_push_image.go`
- Why fragile: Notification uses nested `OneOf` parameter variants with schema references. Manual schema ref strings like `WithSchemaRef("TextType")` are prone to typos. No compile-time validation that references match defined types
- Safe modification: Add unit tests to validate all `WithSchemaRef` references exist. Use constants for schema names
- Test coverage: Gaps - only basic CUE structure validation in `workflowsteps_test.go`

**Parameter Definition Chaining - defkit Fluent API:**
- Files: All definition files use chained `.WithFields()`, `.Values()`, `.Description()` calls
- Why fragile: Long fluent chains like `defkit.Array("ports").WithFields(...).WithFields(...).WithFields(...)` are difficult to refactor. A single typo in field name breaks the entire parameter structure. No syntax validation until E2E runtime
- Safe modification: Extract complex parameter definitions into helper functions (e.g., `func portsParameter() *defkit.ParamSpec`). Test helpers independently
- Test coverage: No unit tests validating parameter names, defaults, or structure before CUE generation

**Custom KubeVela Fork Integration:**
- Files: `go.mod:16` replace directive points to `github.com/guidewire-oss/kubevela`. All imports from `github.com/oam-dev/kubevela` resolve to fork
- Why fragile: If fork falls behind upstream, defkit API changes won't be caught. If upstream introduces breaking changes to parameter types, definitions won't compile
- Safe modification: Pin fork to specific commit tags. Run tests against both upstream and fork versions
- Test coverage: CI only tests against fork, not upstream

## Missing Critical Features

**No API Versioning for Definitions:**
- Problem: Definitions have no version metadata or migration path. If a definition parameter is renamed or removed, users' existing applications break
- Blocks: Safe rolling updates. Backward compatibility guarantees. Definition upgrade paths
- Priority: High

**No Parameter Validation at Definition Time:**
- Problem: defkit allows invalid parameter structures (e.g., conflicting defaults, invalid enum values) that only surface at E2E test time
- Blocks: Early error detection. Faster development feedback loop
- Priority: Medium

**No Definition Schema Documentation Generation:**
- Problem: CUE definitions are generated at runtime but not exported as OpenAPI specs or JSON Schema for documentation/IDE tooling
- Blocks: IDE autocompletion. Documentation generation. CLI parameter discovery
- Priority: Medium

## Scaling Limits

**Serial Definition Registration:**
- Current approach: All definitions register via `init()` functions at startup. Registration order is undefined. Large modules with 100+ definitions may have slow startup
- Limit: Registration is O(n) with definition count. No lazy loading
- Scaling path: Implement definition registry with lazy loading. Move registration to module load time instead of init()
- Priority: Low (current scale ~60 definitions, not a bottleneck yet)

**Flat Package Structure for Large Definition Sets:**
- Current capacity: 56 Go files across 4 directories (components, traits, policies, workflowsteps)
- Scaling limit: IDE indexing and grep performance degrade with 100+ files per directory
- Scaling path: Organize definitions into sub-packages by category (e.g., `traits/networking/`, `traits/resource-management/`)
- Priority: Low

## Dependencies at Risk

**Kubernetes API Version Lock:**
- Risk: `k8s.io/api v0.31.10` is already at a specific patch version. If Kubernetes APIs deprecate fields used in definitions, this breaks
- Impact: Definitions target KubeVela's supported K8s versions. Check KubeVela's K8s compatibility matrix before updating
- Migration plan: Update K8s API versions in lockstep with KubeVela releases. Test against multiple K8s versions in CI
- Files: `go.mod:10`

**Controller-Runtime Version:**
- Risk: `sigs.k8s.io/controller-runtime v0.19.7` may lag upstream. If controller-runtime introduces breaking changes to client APIs, defkit may not expose them
- Impact: New KubeVela controller features may require newer controller-runtime
- Migration plan: Keep controller-runtime aligned with KubeVela's version
- Files: `go.mod:12`

## Known Limitations

**Manual E2E Cleanup Required:**
- Issue: E2E tests create namespaces (`e2e-<app-name>`) that don't auto-clean if tests abort
- Files: `Makefile:62-73` provides manual cleanup via `cleanup-e2e-namespaces` target
- Impact: Namespace accumulation in live clusters if tests fail. Requires manual intervention or CI cleanup hooks
- Workaround: Makefile provides `force-cleanup-e2e-namespaces` to remove stuck finalizers
- Priority: Low (CI workflows should call cleanup, but manual verification needed)

**go.mod Default Parallelism vs Makefile:**
- Issue: Makefile defaults to `PROCS=10` for E2E tests, but documentation says `PROCS=4`
- Files: `Makefile:18` sets `PROCS ?= 10`, `CLAUDE.md:29` documents default as `PROCS=10`
- Impact: E2E tests may run with unexpected parallelism if user doesn't set variable. Can cause resource exhaustion on small clusters
- Workaround: Explicitly set `PROCS` variable in CI or local runs
- Priority: Low

## Security Considerations

**No Secrets Management in Definition Parameters:**
- Risk: Definitions like `notification.go` accept secret references via `secretRef`, but validation happens at runtime in KubeVela, not in definitions
- Files: `workflowsteps/notification.go:24-40` uses `stringValueOrSecretRef` helper
- Current mitigation: defkit doesn't validate secret names or keys at definition time. Validation deferred to KubeVela runtime
- Recommendations: Document that secrets must exist in cluster before application deployment. Add examples showing proper secret setup. Consider adding unit tests validating secret reference structure

**No Audit Trail for Definition Changes:**
- Risk: Module.yaml maintainers and categories aren't versioned. No changelog for definition parameter changes
- Files: `module.yaml`
- Current mitigation: None - users relying on module updates may get breaking changes without notice
- Recommendations: Implement definition versioning. Document breaking changes in release notes. Use semantic versioning for module releases

## Performance Considerations

**Large Definition Files with Fluent API Chains:**
- Problem: `components/webservice.go` (514 lines), `components/statefulset.go` (513 lines) use deeply nested parameter definitions
- Files: `components/webservice.go:26-150+`, `components/statefulset.go`
- Cause: Fluent API requires chains of method calls. No way to split definitions across files
- Improvement path: Extract parameter definitions into separate functions (`func webserviceImageParam()`, `func webservicePortsParam()`). Reduces per-file complexity

**E2E Test Parallelism May Exhaust Resources:**
- Problem: Default `PROCS=10` with `E2E_TIMEOUT=10m` may overload small clusters
- Files: `Makefile:18`, `test/e2e/*.go`
- Cause: Each test creates an Application in isolated namespace. 10 parallel tests = 10 concurrent Applications
- Improvement path: Add cluster capacity checks before tests. Scale PROCS based on cluster size. Implement test queueing for resource-constrained environments

## Process & Documentation Concerns

**Incomplete E2E Test Documentation:**
- Problem: README.md and CLAUDE.md don't document E2E test failure modes or debugging steps
- Files: `README.md:93-127`, `CLAUDE.md:29`
- Impact: New contributors can't debug failing E2E tests. Manual troubleshooting required
- Fix approach: Add E2E debugging guide. Document common failure patterns (namespace stuck, cluster overload, timeout). Include kubectl commands for inspection

**Missing CI Workflow Documentation:**
- Problem: `.github/workflows/test-definitions.yaml` and `unit-tests.yaml` exist but not referenced in README
- Files: `.github/workflows/`
- Impact: Users don't know what CI gates apply. PRs may fail due to undocumented CI checks
- Fix approach: Document all CI workflows in CONTRIBUTING.md. Link from README

---

*Concerns audit: 2026-03-11*
