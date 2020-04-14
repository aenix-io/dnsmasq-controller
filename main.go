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

package main

import (
	"flag"
	"os"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	logrzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	dnsmasqv1alpha1 "github.com/kvaps/dnsmasq-controller/api/v1alpha1"
	"github.com/kvaps/dnsmasq-controller/controllers"
	"github.com/kvaps/dnsmasq-controller/pkg/conf"
	"github.com/kvaps/dnsmasq-controller/pkg/server"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = dnsmasqv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	config := conf.GetConfig()
	flag.StringVar(&config.WatchNamespace, "watch-namespace", "", "Namespace the controller watches for updates to Kubernetes objects."+
		" All namespaces are watched if this parameter is left empty.")
	flag.StringVar(&config.ControllerName, "controller", "", "Name of the controller this controller satisfies.")
	flag.StringVar(&config.MetricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&config.EnableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.IntVar(&config.SyncDelay, "sync-delay", 1, "Time in seconds to syncronise dnsmasq configuration.")
	flag.StringVar(&config.DnsmasqConfDir, "conf-dir", "/etc/dnsmasq.d", "Dnsmasq config directory for write configuration to.")
	flag.StringVar(&config.LogLevel, "log-level", "info", "The log level used by the operator.")
	flag.BoolVar(&config.Development, "development", false, "Run the controller in development mode.")
	flag.BoolVar(&config.CleanupDir, "cleanup", false, "Cleanup dnsmasq config directory before start.")

	flag.Parse()
	config.DnsmasqOptions = flag.Args()

	ctrl.SetLogger(logrzap.New(func(o *logrzap.Options) {
		o.Development = config.Development

		if o.Development == false {
			lev := zap.NewAtomicLevel()
			(&lev).UnmarshalText([]byte(config.LogLevel))
			o.Level = &lev
		}
	}))

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	var err error
	config.MyNamespace, _, err = kubeconfig.Namespace()
	if err != nil {
		setupLog.Error(err, "Failed to get watch namespace")
		os.Exit(1)
	}

	if config.ControllerName == "" {
		config.LeaderElectionID = "dnsmasq-controller-leader"
	} else {
		config.LeaderElectionID = config.ControllerName + "-dnsmasq-controller-leader"
	}

	server.Start()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		Namespace:               config.WatchNamespace,
		MetricsBindAddress:      config.MetricsAddr,
		Port:                    9443,
		LeaderElection:          config.EnableLeaderElection,
		LeaderElectionID:        config.LeaderElectionID,
		LeaderElectionNamespace: config.MyNamespace,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.DnsmasqOptionSetReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("DnsmasqOptionSet"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "DnsmasqOptionSet")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
