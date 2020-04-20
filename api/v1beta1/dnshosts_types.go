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

// DnsHost holds the mapping between IP and hostnames that will be added to dnsmasq hosts file.
type DnsHost struct {
	// IP address of the host file entry.
	IP string `json:"ip" protobuf:"bytes,1,opt,name=ip"`
	// Hostnames for the above IP address.
	Hostnames []string `json:"hostnames,omitempty" protobuf:"bytes,2,rep,name=hostnames"`
}

// DnsHostsSpec defines the desired state of DnsHosts
type DnsHostsSpec struct {
	Controller string    `json:"controller,omitempty"`
	Hosts      []DnsHost `json:"hosts"`
}

// DnsHostsStatus defines the observed state of DnsHosts
type DnsHostsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// DnsHosts is the Schema for the dnshosts API
type DnsHosts struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DnsHostsSpec   `json:"spec,omitempty"`
	Status DnsHostsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DnsHostsList contains a list of DnsHosts
type DnsHostsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsHosts `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsHosts{}, &DnsHostsList{})
}
