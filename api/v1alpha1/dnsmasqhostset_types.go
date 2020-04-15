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

// DnsmasqHost holds the mapping between IP and hostnames that will be added to dnsmasq hosts file.
type DnsmasqHost struct {
	// IP address of the host file entry.
	IP string `json:"ip" protobuf:"bytes,1,opt,name=ip"`
	// Hostnames for the above IP address.
	Hostnames []string `json:"hostnames,omitempty" protobuf:"bytes,2,rep,name=hostnames"`
}

// DnsmasqHostSetSpec defines the desired state of DnsmasqHostSet
type DnsmasqHostSetSpec struct {
	Controller string        `json:"controller,omitempty"`
	Hosts      []DnsmasqHost `json:"hosts"`
}

// DnsmasqHostSetStatus defines the observed state of DnsmasqHostSet
type DnsmasqHostSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// DnsmasqHostSet is the Schema for the dnsmasqhostsets API
type DnsmasqHostSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DnsmasqHostSetSpec   `json:"spec,omitempty"`
	Status DnsmasqHostSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DnsmasqHostSetList contains a list of DnsmasqHostSet
type DnsmasqHostSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsmasqHostSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsmasqHostSet{}, &DnsmasqHostSetList{})
}
