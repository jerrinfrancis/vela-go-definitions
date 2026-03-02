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

package workflowsteps_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/oam-dev/vela-go-definitions/workflowsteps"
)

var _ = Describe("Deploy WorkflowStep", func() {
	It("should have correct name and CUE output", func() {
		step := workflowsteps.Deploy()

		Expect(step.GetName()).To(Equal("deploy"))
		Expect(step.GetDescription()).To(Equal("A powerful and unified deploy step for components multi-cluster delivery with policies."))

		cue := step.ToCue()

		Expect(cue).To(ContainSubstring(`type: "workflow-step"`))
		Expect(cue).To(ContainSubstring(`"category": "Application Delivery"`))
		Expect(cue).To(ContainSubstring(`"scope": "Application"`))
		Expect(cue).To(ContainSubstring(`auto: *true | bool`))
		Expect(cue).To(ContainSubstring(`policies:`))
		Expect(cue).To(ContainSubstring(`parallelism: *5 | int`))
		Expect(cue).To(ContainSubstring(`ignoreTerraformComponent: *true | bool`))
		Expect(cue).To(ContainSubstring(`multicluster.#Deploy`))
		Expect(cue).To(ContainSubstring(`builtin.#Suspend`))
	})
})

var _ = Describe("Suspend WorkflowStep", func() {
	It("should have correct name and CUE output", func() {
		step := workflowsteps.Suspend()

		Expect(step.GetName()).To(Equal("suspend"))

		cue := step.ToCue()

		Expect(cue).To(ContainSubstring(`type: "workflow-step"`))
		Expect(cue).To(ContainSubstring(`"category": "Process Control"`))
		Expect(cue).To(ContainSubstring(`builtin.#Suspend`))
		Expect(cue).To(ContainSubstring(`duration?:`))
		Expect(cue).To(ContainSubstring(`message?:`))
	})
})

var _ = Describe("ApplyComponent WorkflowStep", func() {
	It("should have correct name and CUE output", func() {
		step := workflowsteps.ApplyComponent()

		Expect(step.GetName()).To(Equal("apply-component"))

		cue := step.ToCue()

		Expect(cue).To(ContainSubstring(`type: "workflow-step"`))
		Expect(cue).To(ContainSubstring(`"category": "Application Delivery"`))
		Expect(cue).To(ContainSubstring(`"scope": "Application"`))
		Expect(cue).To(ContainSubstring(`component:`))
		Expect(cue).To(ContainSubstring(`cluster:`))
		Expect(cue).To(ContainSubstring(`namespace:`))
	})
})

var _ = Describe("All WorkflowSteps Registered", func() {
	type stepEntry struct {
		name  string
		toCue func() string
	}

	allSteps := []stepEntry{
		{"deploy", func() string { return workflowsteps.Deploy().ToCue() }},
		{"suspend", func() string { return workflowsteps.Suspend().ToCue() }},
		{"apply-component", func() string { return workflowsteps.ApplyComponent().ToCue() }},
	}

	for _, tc := range allSteps {
		It("should produce valid CUE for "+tc.name, func() {
			cue := tc.toCue()
			Expect(cue).NotTo(BeEmpty())
			Expect(cue).To(ContainSubstring("{"))
			Expect(cue).To(ContainSubstring("}"))
		})
	}
})
