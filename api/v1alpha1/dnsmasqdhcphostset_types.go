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

// DnsmasqDhcpHost holds the mapping between Macs and IP that will be added to dnsmasq dhcp-hosts file.
type DnsmasqDhcpHost struct {
	Macs      []string `json:"macs,omitempty"`
	ClientIDs []string `json:"clientIDs,omitempty"`
	SetTags   []string `json:"setTags,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	IP        string   `json:"ip,omitempty"`
	Hostname  string   `json:"hostname,omitempty"`
	LeaseTime string   `json:"leaseTime,omitempty"`
	Ignore    bool     `json:"ignore,omitempty"`
}

// DnsmasqDhcpHostSetSpec defines the desired state of DnsmasqDhcpHostSet
type DnsmasqDhcpHostSetSpec struct {
	Controller string            `json:"controller,omitempty"`
	Hosts      []DnsmasqDhcpHost `json:"hosts"`
}

// DnsmasqDhcpHostSetStatus defines the observed state of DnsmasqDhcpHostSet
type DnsmasqDhcpHostSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// DnsmasqDhcpHostSet is the Schema for the dnsmasqdhcphostsets API
type DnsmasqDhcpHostSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DnsmasqDhcpHostSetSpec   `json:"spec,omitempty"`
	Status DnsmasqDhcpHostSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DnsmasqDhcpHostSetList contains a list of DnsmasqDhcpHostSet
type DnsmasqDhcpHostSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsmasqDhcpHostSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsmasqDhcpHostSet{}, &DnsmasqDhcpHostSetList{})
}
