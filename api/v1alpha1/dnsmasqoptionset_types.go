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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DnsmasqOptionSetSpecOption defines the desired state of DnsmasqOptionSetSpec
type DnsmasqOptionSetSpecOption struct {
	// +kubebuilder:validation:Enum=dhcp-range;dhcp-host;dhcp-userclass;dhcp-circuitid;dhcp-remoteid;dhcp-subscrid;dhcp-ignore;dhcp-broadcast;mx-host;dhcp-boot;dhcp-option;dhcp-option-force;server;rev-server;local;domain;dhcp-vendorclass;alias;dhcp-vendorclass;srv-host;txt-record;ptr-record;bootp-dynamic;dhcp-mac;dhcp-ignore-names;rebind-domain-ok;dhcp-match;dhcp-name-match;naptr-record;dhcp-generate-names;cname;pxe-service;add-mac;dhcp-duid;host-record;caa-record;dns-rr;auth-zone;synth-domain
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// DnsmasqOptionSetSpec defines the desired state of DnsmasqOptionSet
type DnsmasqOptionSetSpec struct {
	Controller string                       `json:"controller,omitempty"`
	Options    []DnsmasqOptionSetSpecOption `json:"options,omitempty"`
}

// DnsmasqOptionSetStatus defines the observed state of DnsmasqOptionSet
type DnsmasqOptionSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:resource:shortName=dnsmasqopt;dnsmasqopts
// +kubebuilder:object:root=true

// DnsmasqOptionSet is the Schema for the dnsmasqoptionsets API
type DnsmasqOptionSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DnsmasqOptionSetSpec   `json:"spec,omitempty"`
	Status DnsmasqOptionSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DnsmasqOptionSetList contains a list of DnsmasqOptionSet
type DnsmasqOptionSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsmasqOptionSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsmasqOptionSet{}, &DnsmasqOptionSetList{})
}
