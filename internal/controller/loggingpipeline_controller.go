/*
Copyright 2025.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	jjsiviov1alpha1 "github.com/jjsiv/logging-operator/api/v1alpha1"
)

// LoggingPipelineReconciler reconciles a LoggingPipeline object
type LoggingPipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jjsiv.io,resources=loggingpipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jjsiv.io,resources=loggingpipelines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=jjsiv.io,resources=loggingpipelines/finalizers,verbs=update

func (r *LoggingPipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)
	// Create FluentBitConfig struct from LoggingPipeline
	// Create/update ConfigMap on the cluster from FluentBitConfig
	// To avoid rebuilding entire config everytime, we should build a main config file and then one for each pipeline...
	// So the ConfigMap structure could have structure like this:
	// data:
	//   main.yaml: |
	//     includes:
	//     - *.yaml
	//   my-namespace-my-logging-pipeline.yaml: |
	//     pipeline:
	//       inputs: ...

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LoggingPipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jjsiviov1alpha1.LoggingPipeline{}).
		Named("loggingpipeline").
		Complete(r)
}
