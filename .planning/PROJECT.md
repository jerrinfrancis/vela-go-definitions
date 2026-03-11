# vela-go-definitions defkit API Documentation

## What This Is

A single-page HTML developer reference for the `defkit` fluent Go API used in `vela-go-definitions`. It documents every method a developer needs to create ComponentDefinitions, TraitDefinitions, PolicyDefinitions, and WorkflowStepDefinitions — complete with side-by-side Go and CUE examples and modern navigation.

## Core Value

Any new developer can open one HTML page and understand how to write a defkit definition from scratch, without reading source code.

## Requirements

### Validated

- ✓ defkit Go API for Components, Traits, Policies, WorkflowSteps — existing
- ✓ Parameter builders: String, Int, Bool, Enum, Struct, Array, List, Object, StringKeyMap — existing
- ✓ Template system: tpl.Output, tpl.Patch, tpl.UsePatchContainer — existing
- ✓ Resource builders: NewResource, Set, SetIf, ForEach, etc. — existing
- ✓ Value expressions: Lit, Reference, Interpolation, conditions — existing
- ✓ VelaCtx: Name, AppName, Namespace, Revision — existing

### Active

- [ ] Single-page HTML documentation generated into docs/index.html
- [ ] All primary API methods documented with description, signature, Go example, CUE equivalent
- [ ] Sidebar navigation with section anchors
- [ ] Side-by-side Go + CUE code blocks per method
- [ ] Modern, professional design (dark sidebar, clean typography)
- [ ] Definition builders section (NewComponent, NewTrait, NewPolicy, NewWorkflowStep)
- [ ] Parameter builders + chain methods documented
- [ ] Template methods documented
- [ ] Resource builder methods documented
- [ ] Value expressions and conditions documented
- [ ] VelaCtx documented
- [ ] Full worked example per definition type (Component, Trait, Policy, WorkflowStep)

### Out of Scope

- Auto-generation from Go source/AST — hand-crafted HTML for full design control
- Multi-page site — single page for simplicity
- Helper/internal methods — only public API used in definitions

## Context

- Defkit comes from a custom KubeVela fork: `github.com/guidewire-oss/kubevela`
- The API compiles to CUE — every method corresponds to CUE constructs
- Definitions self-register via `init()` and are exported via `defkit.ToJSON()`
- 4 definition types: Component (workload), Trait (patch), Policy (control), WorkflowStep (automation)
- Parameter types map to CUE types: String→string, Int→int, Bool→bool, Enum→disjunction, Struct→struct, Array→[...T]

## Constraints

- **Output**: `docs/index.html` — single file, no external dependencies (all CSS/JS inline or CDN)
- **Build**: No build toolchain — pure HTML/CSS/JS
- **Scope**: Document creation-time API only — no runtime/internal methods
- **Token limit**: Generate in phases to avoid 32k output token limit

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Single-page HTML | Simpler to maintain, shareable as single file | — Pending |
| Side-by-side Go + CUE | Helps developers understand the mapping | — Pending |
| docs/ directory | GitHub Pages compatible | — Pending |
| Hand-crafted HTML | Full control over design, no build toolchain needed | — Pending |

---
*Last updated: 2026-03-11 after initialization*
