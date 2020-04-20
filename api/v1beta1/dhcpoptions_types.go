/*


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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DhcpOption defines dhcp-option for dnsmasq
type DhcpOption struct {
	// +kubebuilder:validation:Pattern="^([0-9]+|option:.+|option6:.+)$"
	Key     string   `json:"key"`
	Values  []string `json:"values"`
	Tags    []string `json:"tags,omitempty"`
	Encap   string   `json:"encap,omitempty"`
	ViEncap string   `json:"viEncap,omitempty"`
	Vendor  string   `json:"leaseTime,omitempty"`
}

// DhcpOptionsSpec defines the desired state of DhcpOptions
type DhcpOptionsSpec struct {
	Controller string       `json:"controller,omitempty"`
	Options    []DhcpOption `json:"options"`
}

// DhcpOptionsStatus defines the observed state of DhcpOptions
type DhcpOptionsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// +kubebuilder:printcolumn:name="Controller",type="string",JSONPath=".spec.controller"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// DhcpOptions is the Schema for the dhcpoptions API
type DhcpOptions struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DhcpOptionsSpec   `json:"spec,omitempty"`
	Status DhcpOptionsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DhcpOptionsList contains a list of DhcpOptions
type DhcpOptionsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DhcpOptions `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DhcpOptions{}, &DhcpOptionsList{})
}
