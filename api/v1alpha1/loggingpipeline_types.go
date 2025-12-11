/*
Copyright 2025.

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

package v1alpha1

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/jjsiv/logging-operator/api/v1alpha1/input"
	"github.com/jjsiv/logging-operator/api/v1alpha1/output"
	"github.com/jjsiv/logging-operator/internal/fluentbit"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LoggingPipelineSpec defines the desired state of LoggingPipeline
type LoggingPipelineSpec struct {
	Inputs  *PipelineInputs  `json:"inputs,omitempty"`
	Outputs *PipelineOutputs `json:"outputs,omitempty"`
}

type PipelineInputs struct {
	Tail *input.Tail `json:"tail,omitempty"`
}

func (pi *PipelineInputs) ToFluentBitInputs(opts *input.BuildOptions) []fluentbit.Input {
	var inputs []fluentbit.Input
	if pi.Tail != nil {
		inputs = append(inputs, pi.Tail.ToFluentBitInput(opts))
	}

	return inputs
}

type PipelineOutputs struct {
	Stdout *output.Stdout `json:"stdout,omitempty"`
}

func (po *PipelineOutputs) ToFluentBitOutputs(opts *output.BuildOptions) []fluentbit.Output {
	var outputs []fluentbit.Output
	if po.Stdout != nil {
		outputs = append(outputs, po.Stdout.ToFluentBitOutput(opts))
	}

	return outputs
}

// LoggingPipelineStatus defines the observed state of LoggingPipeline.
type LoggingPipelineStatus struct {
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=lp;lps
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"

// LoggingPipeline is the Schema for the loggingpipelines API
type LoggingPipeline struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of LoggingPipeline
	// +required
	Spec LoggingPipelineSpec `json:"spec"`

	// status defines the observed state of LoggingPipeline
	// +optional
	Status LoggingPipelineStatus `json:"status,omitzero"`
}

// GenerateTag generates a unique tag for this pipeline
func (lp *LoggingPipeline) GenerateTag() string {
	input := fmt.Sprintf("%s-%s", lp.Namespace, lp.Name)
	hash := sha256.Sum256([]byte(input))
	tag := hex.EncodeToString(hash[:])[:16]

	return tag + ".*"
}

// +kubebuilder:object:root=true

// LoggingPipelineList contains a list of LoggingPipeline
type LoggingPipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []LoggingPipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoggingPipeline{}, &LoggingPipelineList{})
}
