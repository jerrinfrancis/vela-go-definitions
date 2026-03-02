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

package policies_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/policies"
)

var _ = Describe("Topology Policy", func() {
	It("should have correct name and CUE output", func() {
		policy := policies.Topology()

		Expect(policy.GetName()).To(Equal("topology"))
		Expect(policy.GetDescription()).To(Equal("Describe the destination where components should be deployed to."))

		cue := policy.ToCue()

		Expect(cue).To(ContainSubstring(`type: "policy"`))
		Expect(cue).To(ContainSubstring(`clusters?:`))
		Expect(cue).To(ContainSubstring(`clusterLabelSelector?:`))
		Expect(cue).To(ContainSubstring(`allowEmpty?:`))
		Expect(cue).To(ContainSubstring(`namespace?:`))
	})
})

var _ = Describe("Override Policy", func() {
	It("should have correct name and CUE output", func() {
		policy := policies.Override()

		Expect(policy.GetName()).To(Equal("override"))

		cue := policy.ToCue()

		Expect(cue).To(ContainSubstring(`type: "policy"`))
		Expect(cue).To(ContainSubstring(`#PatchParams`))
		Expect(cue).To(ContainSubstring(`name?:`))
		Expect(cue).To(ContainSubstring(`type?:`))
		Expect(cue).To(ContainSubstring(`properties?:`))
		Expect(cue).To(ContainSubstring(`traits?:`))
		Expect(cue).To(ContainSubstring(`disable: *false | bool`))
		Expect(cue).To(ContainSubstring(`components?:`))
		Expect(cue).To(ContainSubstring(`selector?:`))
	})
})

var _ = Describe("GarbageCollect Policy", func() {
	It("should have correct name and CUE output", func() {
		policy := policies.GarbageCollect()

		Expect(policy.GetName()).To(Equal("garbage-collect"))

		cue := policy.ToCue()

		Expect(cue).To(ContainSubstring(`type: "policy"`))
		Expect(cue).To(ContainSubstring(`#GarbageCollectPolicyRule`))
		Expect(cue).To(ContainSubstring(`#ResourcePolicyRuleSelector`))
		Expect(cue).To(ContainSubstring(`applicationRevisionLimit?:`))
		Expect(cue).To(ContainSubstring(`keepLegacyResource: *false | bool`))
		Expect(cue).To(ContainSubstring(`continueOnFailure: *false | bool`))
		Expect(cue).To(ContainSubstring(`rules?:`))
		Expect(cue).To(ContainSubstring(`strategy: *"onAppUpdate"`))
		Expect(cue).To(ContainSubstring(`componentNames?:`))
		Expect(cue).To(ContainSubstring(`componentTypes?:`))
	})
})

var _ = Describe("All Policies Registered", func() {
	type policyEntry struct {
		name  string
		toCue func() string
	}

	allPolicies := []policyEntry{
		{"topology", func() string { return policies.Topology().ToCue() }},
		{"override", func() string { return policies.Override().ToCue() }},
		{"garbage-collect", func() string { return policies.GarbageCollect().ToCue() }},
	}

	for _, tc := range allPolicies {
		It("should produce valid CUE for "+tc.name, func() {
			cue := tc.toCue()
			Expect(cue).NotTo(BeEmpty())
			Expect(cue).To(ContainSubstring("{"))
			Expect(cue).To(ContainSubstring("}"))
		})
	}
})
