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

package controllers

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsmasqv1alpha1 "github.com/kvaps/dnsmasq-controller/api/v1alpha1"
	"github.com/kvaps/dnsmasq-controller/pkg/conf"
	"github.com/kvaps/dnsmasq-controller/pkg/util"
)

// DnsmasqHostSetReconciler reconciles a DnsmasqHostSet object
type DnsmasqHostSetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dnsmasq.kvaps.cf,resources=dnsmasqhostsets,verbs=get;list;watch

func (r *DnsmasqHostSetReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("dnshost", req.NamespacedName)
	config := conf.GetConfig()

	configFile := config.DnsmasqConfDir + "/hosts/" + req.Namespace + "-" + req.Name

	res := &dnsmasqv1alpha1.DnsmasqHostSet{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, res)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found
			if _, err := os.Stat(configFile); !os.IsNotExist(err) {
				os.Remove(configFile)
				r.Log.Info("Removed " + configFile)
				config.Generation++
			}
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if res.Spec.Controller != config.ControllerName {
		if _, err := os.Stat(configFile); !os.IsNotExist(err) {
			// Controller name has been changed
			os.Remove(configFile)
			r.Log.Info("Removed " + configFile)
			config.Generation++
		}
		return ctrl.Result{}, nil
	}

	// Write hosts
	var configData string
	for _, h := range res.Spec.Hosts {
		configData += h.IP
		for _, hostname := range h.Hostnames {
			configData += " " + hostname
		}
		configData += "\n"
	}
	configBytes := []byte(configData)

	configWritten, err := util.WriteConfig(configFile, configFile, configBytes)
	if err != nil {
		r.Log.Error(err, "Failed to update "+configFile)
		return ctrl.Result{}, nil
	}

	if configWritten {
		r.Log.Info("Written " + configFile)
		config.Generation++
	}

	return ctrl.Result{}, nil
}

func (r *DnsmasqHostSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsmasqv1alpha1.DnsmasqHostSet{}).
		Complete(r)
}
