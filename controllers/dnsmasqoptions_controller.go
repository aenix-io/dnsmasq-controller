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

	dnsmasqv1beta1 "github.com/kvaps/dnsmasq-controller/api/v1beta1"
	"github.com/kvaps/dnsmasq-controller/pkg/conf"
	"github.com/kvaps/dnsmasq-controller/pkg/util"
)

// DnsmasqOptionsReconciler reconciles a DnsmasqOptions object
type DnsmasqOptionsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dnsmasq.kvaps.cf,resources=dnsmasqoptions,verbs=get;list;watch

func (r *DnsmasqOptionsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("dnsmasqconfiguration", req.NamespacedName)
	config := conf.GetConfig()

	configFile := config.DnsmasqConfDir + "/" + req.Namespace + "-" + req.Name + ".conf"
	tmpConfigFile := config.DnsmasqConfDir + "/." + req.Namespace + "-" + req.Name + ".conf.tmp"

	res := &dnsmasqv1beta1.DnsmasqOptions{}
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

	// Write options
	var configData string
	for _, o := range res.Spec.Options {
		if o.Key == "dhcp-range" && !config.EnableDHCP {
			continue
		}
		configData += o.Key + "="
		configValues := ""
		for _, v := range o.Values {
			configValues += "," + v
		}
		configData += configValues[1:]
		configData += "\n"
	}
	configBytes := []byte(configData)

	configWritten, err := util.WriteConfig(configFile, tmpConfigFile, configBytes)
	if err != nil {
		r.Log.Error(err, "Failed to update "+configFile)
		return ctrl.Result{}, nil
	}

	if configWritten {
		if err = util.TestConfig(tmpConfigFile); err != nil {
			//os.Remove(tmpConfigFile)
			r.Log.Error(err, "Config "+tmpConfigFile+" is invalid!")
			return ctrl.Result{}, nil
		}

		if err = os.Rename(tmpConfigFile, configFile); err != nil {
			os.Remove(tmpConfigFile)
			r.Log.Error(err, "Failed to move "+tmpConfigFile+" to "+configFile)
			return ctrl.Result{}, nil
		}
		r.Log.Info("Written " + configFile)
		config.Generation++
	}

	return ctrl.Result{}, nil

}

func (r *DnsmasqOptionsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsmasqv1beta1.DnsmasqOptions{}).
		Complete(r)
}
