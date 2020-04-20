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

// DnsmasqOption defines option for dnsmasq
type DnsmasqOption struct {
	// +kubebuilder:validation:Enum=dhcp-range;dhcp-host;dhcp-userclass;dhcp-circuitid;dhcp-remoteid;dhcp-subscrid;dhcp-ignore;dhcp-broadcast;mx-host;dhcp-boot;dhcp-option;dhcp-option-force;server;rev-server;local;domain;dhcp-vendorclass;alias;dhcp-vendorclass;srv-host;txt-record;ptr-record;bootp-dynamic;dhcp-mac;dhcp-ignore-names;rebind-domain-ok;dhcp-match;dhcp-name-match;naptr-record;dhcp-generate-names;cname;pxe-service;add-mac;dhcp-duid;host-record;caa-record;dns-rr;auth-zone;synth-domain
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// DnsmasqOptionsSpec defines the desired state of DnsmasqOptions
type DnsmasqOptionsSpec struct {
	Controller string          `json:"controller,omitempty"`
	Options    []DnsmasqOption `json:"options"`
}

// DnsmasqOptionsStatus defines the observed state of DnsmasqOptions
type DnsmasqOptionsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// DnsmasqOptions is the Schema for the dnsmasqoptions API
type DnsmasqOptions struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DnsmasqOptionsSpec   `json:"spec,omitempty"`
	Status DnsmasqOptionsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DnsmasqOptionsList contains a list of DnsmasqOptions
type DnsmasqOptionsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsmasqOptions `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsmasqOptions{}, &DnsmasqOptionsList{})
}
