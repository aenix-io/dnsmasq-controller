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

// DnsmasqDhcpHost holds the mapping between Macs and IP that will be added to dnsmasq dhcp-hosts file.
type DhcpHost struct {
	Macs      []string `json:"macs,omitempty"`
	ClientIDs []string `json:"clientIDs,omitempty"`
	SetTags   []string `json:"setTags,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	IP        string   `json:"ip,omitempty"`
	Hostname  string   `json:"hostname,omitempty"`
	LeaseTime string   `json:"leaseTime,omitempty"`
	Ignore    bool     `json:"ignore,omitempty"`
}

// DhcpHostsSpec defines the desired state of DhcpHosts
type DhcpHostsSpec struct {
	Controller string     `json:"controller,omitempty"`
	Hosts      []DhcpHost `json:"hosts"`
}

// DhcpHostsStatus defines the observed state of DhcpHosts
type DhcpHostsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// DhcpHosts is the Schema for the dhcphosts API
type DhcpHosts struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DhcpHostsSpec   `json:"spec,omitempty"`
	Status DhcpHostsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DhcpHostsList contains a list of DhcpHosts
type DhcpHostsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DhcpHosts `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DhcpHosts{}, &DhcpHostsList{})
}
