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
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	Source client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Dest   client.Client
}

func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("instance", req.NamespacedName)

	source := &Instance{}
	if err := r.Source.Get(ctx, req.NamespacedName, source); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	dest := &Instance{}
	if err := r.Dest.Get(ctx, req.NamespacedName, dest); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	dest.Spec = source.Spec
	if err := r.Dest.Update(ctx, dest); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update destination instance: %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&Instance{
			Object:
			Spec:
		}).
		Complete(r)
}

type Object interface {
	metav1.Object
	runtime.Object
}

type Instance struct {
	Object
	Spec interface{}
}
