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
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dnsmasqv1alpha1 "github.com/kvaps/dnsmasq-controller/api/v1alpha1"
	"github.com/kvaps/dnsmasq-controller/pkg/conf"
)

// DnsmasqOptionSetReconciler reconciles a DnsmasqOptionSet object
type DnsmasqOptionSetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dnsmasq.kvaps.cf,resources=dnsmasqoptionsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dnsmasq.kvaps.cf,resources=dnsmasqoptionsets/status,verbs=get;update;patch

func (r *DnsmasqOptionSetReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("dnsmasqoptionset", req.NamespacedName)
	config := conf.GetConfig()

	configFile := config.DnsmasqConfDir + "/" + req.Namespace + "-" + req.Name + ".conf"
	tmpConfigFile := config.DnsmasqConfDir + "/." + req.Namespace + "-" + req.Name + ".conf.tmp"

	res := &dnsmasqv1alpha1.DnsmasqOptionSet{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, res)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found
			if _, err := os.Stat(configFile); !os.IsNotExist(err) {
				r.Log.Info("Deleting config " + configFile)
				os.Remove(configFile)
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
			r.Log.Info("Deleting config " + configFile)
			os.Remove(configFile)
			config.Generation++
		}
		return ctrl.Result{}, nil
	}

	var configData string
	for _, opt := range res.Spec.Options {
		configData += opt.Key + "=" + opt.Value + "\n"
	}
	configBytes := []byte(configData)

	// If file exists check hash
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		hasher := md5.New()
		f, err := os.Open(configFile)
		if err != nil {
			r.Log.Error(err, "Failed to read config "+configFile)
			return ctrl.Result{}, nil
		}
		defer f.Close()
		if _, err := io.Copy(hasher, f); err != nil {
			r.Log.Error(err, "Can not calculate hash of "+configFile)
			return ctrl.Result{}, nil
		}
		oldHash := hasher.Sum(nil)[:16]

		hasher = md5.New()
		hasher.Write(configBytes)
		newHash := hasher.Sum(nil)

		if bytes.Equal(oldHash, newHash) {
			return ctrl.Result{}, nil
		}
		f.Close()
	}

	r.Log.Info("Writing config " + configFile)
	err = ioutil.WriteFile(tmpConfigFile, configBytes, 0644)
	if err != nil {
		r.Log.Error(err, "Failed to write "+tmpConfigFile)
		return ctrl.Result{}, nil
	}

	if err = testConfig(tmpConfigFile); err != nil {
		//os.Remove(tmpConfigFile)
		r.Log.Error(err, "Config "+tmpConfigFile+" is invalid!")
		return ctrl.Result{}, nil
	}

	if err = os.Rename(tmpConfigFile, configFile); err != nil {
		os.Remove(tmpConfigFile)
		r.Log.Error(err, "Failed to move "+tmpConfigFile+" to "+configFile)
		return ctrl.Result{}, nil
	}

	config.Generation++

	return ctrl.Result{}, nil
}

func (r *DnsmasqOptionSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsmasqv1alpha1.DnsmasqOptionSet{}).
		Complete(r)
}

func testConfig(f string) error {
	var stderr bytes.Buffer
	dnsmasqBinary, err := exec.LookPath("dnsmasq")
	if err != nil {
		return err
	}
	cmd := &exec.Cmd{
		Path:   dnsmasqBinary,
		Args:   []string{"dnsmasq", "--test", "--conf-file=" + f},
		Stderr: &stderr,
	}
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf(string(stderr.Bytes()))
	}
	return err
}
