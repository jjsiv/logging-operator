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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FluentBitSpec defines the desired state of FluentBit
type FluentBitSpec struct {
	// FluentBit image to use for the DaemonSet.
	// +required
	Image string `json:"image,omitempty"`
	// Config defines FluentBit global properties (service section).
	// See https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/yaml/service-section
	Config           FluentBitConfig      `json:"config,omitempty"`
	PipelineSelector metav1.LabelSelector `json:"pipelineSelector,omitempty"`
}

// FluentBitConfig defines service level configuration for the FluentBit instance.
type FluentBitConfig struct {
	// Sets the flush time in seconds.nanoseconds.
	// The engine loop uses a flush timeout to define when to flush the records ingested by input plugins through the defined output plugins.
	Flush *float64 `json:"flush,omitempty"`
	// Sets the grace time in seconds as an integer value. The engine loop uses a grace timeout to define the wait time on exit.
	Grace *int64 `json:"grace,omitempty"`
	// Sets the logging verbosity level.
	// Possible values: off, error, warn, info, debug, and trace.
	// Values are cumulative. For example, if debug is set, it will include error, warning, info, and debug.
	// +kubebuilder:validation:Enum=off;error;warn;info;debug;trace
	LogLevel *string `json:"logLevel,omitempty"`
}

// FluentBitStatus defines the observed state of FluentBit.
type FluentBitStatus struct {
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=fb;flb;fbs;flbs
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"

// FluentBit defines the state of a FluentBit DaemonSet and its core configuration.
type FluentBit struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of FluentBit
	// +required
	Spec FluentBitSpec `json:"spec"`

	// status defines the observed state of FluentBit
	// +optional
	Status FluentBitStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// FluentBitList contains a list of FluentBit
type FluentBitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []FluentBit `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FluentBit{}, &FluentBitList{})
}
