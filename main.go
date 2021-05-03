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

	"github.com/LimKianAn/sync-crd/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	cr "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	schm     = runtime.NewScheme()
	setupLog = cr.Log.WithName("setup")
)

func main() {
	var metricsAddr, source, destination, group, version string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&source, "source-cluster-kubeconfig", "", "The path to the kubeconfig of the source cluster")
	flag.StringVar(&destination, "destination-cluster-kubeconfig", "", "The path to the kubeconfig of the destination cluster")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&group, "api-group", "", "")
	flag.StringVar(&version, "api-version", "", "")

	flag.Parse()

	cr.SetLogger(zap.New(zap.UseDevMode(true)))

	_ = (&scheme.Builder{GroupVersion: schema.GroupVersion{Group: group, Version: version}}).AddToScheme(schm)

	SourceConfig, err := clientcmd.BuildConfigFromFlags("", source)
	if err != nil {
		setupLog.Error(err, "unable to build config from flags")
		os.Exit(1)
	}

	mgr, err := cr.NewManager(SourceConfig, cr.Options{
		Scheme:             schm,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "sync-crd",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	destinationConfig, err := clientcmd.BuildConfigFromFlags("", destination)
	if err != nil {
		setupLog.Error(err, "unable to build config from flags")
		os.Exit(1)
	}
	destK8sClient, err := client.New(destinationConfig, client.Options{Scheme: schm})

	if err = (&controllers.InstanceReconciler{
		Source: mgr.GetClient(),
		Log:    cr.Log.WithName("controllers").WithName("Instance"),
		Scheme: mgr.GetScheme(),
		Dest:   destK8sClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Instance")
		os.Exit(1)
	}
}
