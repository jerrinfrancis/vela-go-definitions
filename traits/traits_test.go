/*
Copyright 2025 The KubeVela Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package traits_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/traits"
)

var _ = Describe("All Traits Registered", func() {
	type traitEntry struct {
		name   string
		toCue  func() string
	}

	allTraits := []traitEntry{
		{"scaler", func() string { return traits.Scaler().ToCue() }},
		{"labels", func() string { return traits.Labels().ToCue() }},
		{"annotations", func() string { return traits.Annotations().ToCue() }},
		{"expose", func() string { return traits.Expose().ToCue() }},
		{"sidecar", func() string { return traits.Sidecar().ToCue() }},
		{"env", func() string { return traits.Env().ToCue() }},
		{"resource", func() string { return traits.Resource().ToCue() }},
		{"affinity", func() string { return traits.Affinity().ToCue() }},
		{"hpa", func() string { return traits.HPA().ToCue() }},
		{"init-container", func() string { return traits.InitContainer().ToCue() }},
		{"service-account", func() string { return traits.ServiceAccount().ToCue() }},
		{"gateway", func() string { return traits.Gateway().ToCue() }},
		{"service-binding", func() string { return traits.ServiceBinding().ToCue() }},
		{"startup-probe", func() string { return traits.StartupProbe().ToCue() }},
		{"securitycontext", func() string { return traits.SecurityContext().ToCue() }},
		{"container-image", func() string { return traits.ContainerImage().ToCue() }},
	}

	for _, tc := range allTraits {
		It("should produce valid CUE for "+tc.name, func() {
			cue := tc.toCue()
			Expect(cue).NotTo(BeEmpty())
			Expect(cue).To(ContainSubstring("{"))
			Expect(cue).To(ContainSubstring("}"))
		})
	}
})

// PatchFieldBuilderPatterns verifies that the PatchField builder methods
// (.IsSet(), .NotEmpty(), .Default(), .Int(), .Bool(), .StringArray(), .Target(), .Strategy())
// used in the three PatchContainer-based traits produce the correct CUE output patterns.
var _ = Describe("PatchField Builder Patterns", func() {
	Context("IsSet generates != _|_ guard and optional param syntax", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.StartupProbe().ToCue()

			// .IsSet() alone → optional param (field?: type) + guarded in PatchContainer body
			Expect(cue).To(ContainSubstring(`exec?: {`))
			Expect(cue).To(ContainSubstring(`if _params.exec != _|_`))

			// .Int().IsSet() → optional int param + guarded
			Expect(cue).To(ContainSubstring(`terminationGracePeriodSeconds?: int`))
			Expect(cue).To(ContainSubstring(`if _params.terminationGracePeriodSeconds != _|_`))

			// .Int().IsSet().Default("0") → default value in param + guarded in PatchContainer body
			Expect(cue).To(ContainSubstring(`initialDelaySeconds: *0 | int`))
			Expect(cue).To(ContainSubstring(`if _params.initialDelaySeconds != _|_`))
		})
	})

	Context("Default without IsSet generates unguarded assignment", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.SecurityContext().ToCue()

			// .Bool().Default("false") → default value, no guard in PatchContainer body
			Expect(cue).To(ContainSubstring(`allowPrivilegeEscalation: *false | bool`))
			Expect(cue).To(ContainSubstring(`allowPrivilegeEscalation: _params.allowPrivilegeEscalation`))

			// .Int().IsSet() → optional, guarded
			Expect(cue).To(ContainSubstring(`runAsUser?: int`))
			Expect(cue).To(ContainSubstring(`if _params.runAsUser != _|_`))
		})
	})

	Context("Target remaps param name to different container field", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.SecurityContext().ToCue()

			// .Target("add") maps addCapabilities param → add field in container
			Expect(cue).To(ContainSubstring(`addCapabilities?: [...string]`))
			Expect(cue).To(ContainSubstring(`add: _params.addCapabilities`))

			// .Target("drop") maps dropCapabilities param → drop field in container
			Expect(cue).To(ContainSubstring(`dropCapabilities?: [...string]`))
			Expect(cue).To(ContainSubstring(`drop: _params.dropCapabilities`))
		})
	})

	Context("NotEmpty generates != empty string guard", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.ContainerImage().ToCue()

			// .NotEmpty() → guarded with != "" in PatchContainer body
			Expect(cue).To(ContainSubstring(`if _params.imagePullPolicy != ""`))

			// .NotEmpty() should NOT make the field optional (no ? suffix)
			Expect(cue).To(ContainSubstring(`imagePullPolicy: *""`))
			Expect(cue).NotTo(ContainSubstring(`imagePullPolicy?:`))
		})
	})

	Context("Strategy generates patchStrategy annotation", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.ContainerImage().ToCue()

			// .Strategy("retainKeys") → // +patchStrategy=retainKeys annotation
			Expect(cue).To(ContainSubstring(`// +patchStrategy=retainKeys`))
		})
	})

	Context("StringArray generates typed array", func() {
		It("should produce correct CUE patterns", func() {
			cue := traits.SecurityContext().ToCue()

			// .StringArray().IsSet() → optional typed array
			Expect(cue).To(ContainSubstring(`addCapabilities?: [...string]`))
			Expect(cue).To(ContainSubstring(`dropCapabilities?: [...string]`))
		})
	})
})
