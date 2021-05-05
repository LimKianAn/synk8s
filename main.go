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
	"context"
	"flag"
	"fmt"
	"os"

	api "github.com/LimKianAn/sync-crd/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	cr "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Change Instance to your CRD struct name
type CRD = api.Instance

var (
	flags = struct {
		metricsAddr, source, destination string
		enableLeaderElection             bool
	}{}
	log    = cr.Log.WithName("setup")
	scheme = runtime.NewScheme()
)

func init() {
	_ = api.AddToScheme(scheme)
}

func main() {
	parseflags()
	cr.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := NewMgr()
	if err != nil {
		log.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err := mgr.RegisterCRDController(); err != nil {
		log.Error(err, "unable to create a new controller")
		os.Exit(1)
	}

	if err := mgr.Start(context.Background()); err != nil {
		log.Error(err, "unable to start a manager")
		os.Exit(1)
	}
}

func NewMgr() (*Mgr, error) {
	SourceConfig, err := clientcmd.BuildConfigFromFlags("", flags.source)
	if err != nil {
		log.Error(err, "unable to build config from flags")
		return nil, err
	}

	mgr, err := cr.NewManager(SourceConfig, cr.Options{
		Scheme:             scheme,
		MetricsBindAddress: flags.metricsAddr,
		Port:               9443,
		LeaderElection:     flags.enableLeaderElection,
		LeaderElectionID:   "sync-crd",
	})
	if err != nil {
		log.Error(err, "unable to create a new manager")
		return nil, err
	}

	return &Mgr{Manager: mgr}, nil
}

type Mgr struct {
	manager.Manager
}

func (m *Mgr) RegisterCRDController() error {
	destinationConfig, err := clientcmd.BuildConfigFromFlags("", flags.destination)
	if err != nil {
		return fmt.Errorf("unable to build config from flags: %w", err)
	}

	destK8sClient, err := client.New(destinationConfig, client.Options{Scheme: scheme})
	if err != nil {
		return fmt.Errorf("unable to create a new k8s client: %w", err)
	}

	return cr.NewControllerManagedBy(m).
		For(&CRD{}).
		Complete(&InstanceReconciler{
			Source: m.GetClient(),
			Log:    cr.Log.WithName("controllers").WithName("CRD"),
			Scheme: m.GetScheme(),
			Dest:   destK8sClient,
		})
}

func parseflags() {
	flag.StringVar(&flags.metricsAddr, "metrics-addr", ":9999", "The address the metric endpoint binds to")
	flag.StringVar(&flags.source, "source-cluster-kubeconfig", "", "The path to the kubeconfig of the source cluster")
	flag.StringVar(&flags.destination, "destination-cluster-kubeconfig", "", "The path to the kubeconfig of the destination cluster")
	flag.BoolVar(&flags.enableLeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()
}

type InstanceReconciler struct {
	Source client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Dest   client.Client
}

func (r *InstanceReconciler) Reconcile(ctx context.Context, req cr.Request) (cr.Result, error) {
	log := r.Log.WithValues("instance", req.NamespacedName)

	source := &CRD{}
	if err := r.Source.Get(ctx, req.NamespacedName, source); err != nil {
		return cr.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("source instance fetched")

	dest := &CRD{}
	if err := r.Dest.Get(ctx, req.NamespacedName, dest); err != nil {
		return cr.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("destination instance fetched")

	dest.Spec = source.Spec
	if err := r.Dest.Update(ctx, dest); err != nil {
		return cr.Result{}, fmt.Errorf("failed to update destination instance: %w", err)
	}
	log.Info("destination instance updated")

	return cr.Result{}, nil
}
