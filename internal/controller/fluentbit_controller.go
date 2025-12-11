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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/goccy/go-yaml"
	"github.com/jjsiv/logging-operator/api/v1alpha1"
	"github.com/jjsiv/logging-operator/api/v1alpha1/input"
	"github.com/jjsiv/logging-operator/api/v1alpha1/output"
	"github.com/jjsiv/logging-operator/internal/fluentbit"
)

// FluentBitReconciler reconciles a FluentBit object
type FluentBitReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=jjsiv.io,resources=fluentbits,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=jjsiv.io,resources=fluentbits/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=jjsiv.io,resources=fluentbits/finalizers,verbs=update
// +kubebuilder:rbac:groups=jjsiv.io,resources=loggingpipelines,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;update;patch

// TODO: perhaps we should not actually handle config creation and daemonset here.
// Using another controller that is able to trigger reconciliation whenever needed might be better
// LoggingPipelines controller would need to be aware where to store the configuration though.
// LP controller could list FluentBits to check which one matches its pipeline and trigger reconciliation as needed.
func (r *FluentBitReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx).WithValues("fluentbit", req.NamespacedName)
	logger.Info("reconciliation started")

	var flb v1alpha1.FluentBit
	if err := r.Get(ctx, req.NamespacedName, &flb); err != nil {
		logger.Error(err, "unable to fetch FluentBit")
		return ctrl.Result{}, err
	}

	configs := make(map[string]*fluentbit.FluentBitConfig)
	configs["main.yaml"] = buildMainFluentBitConfig(&flb)
	logger.Info("built main FluentBit config")

	selectors, err := metav1.LabelSelectorAsSelector(&flb.Spec.PipelineSelector)
	if err != nil {
		logger.Error(err, "failed to select LoggingPipelines based on selector")
		return ctrl.Result{}, err
	}

	var lpList v1alpha1.LoggingPipelineList
	if err := r.List(ctx, &lpList, &client.ListOptions{LabelSelector: selectors}); err != nil {
		logger.Error(err, "failed to list LoggingPipelines")
		return ctrl.Result{}, err
	}

	lps := lpList.Items
	logger.Info("matched LoggingPipelines", "count", len(lps))

	for _, lp := range lps {
		cmKeyName := lp.Namespace + "-" + lp.Name + ".yaml"
		configs[cmKeyName] = buildFluentBitPipelineConfig(&lp)
	}

	cmData := make(map[string]string)
	for key, conf := range configs {
		data, err := yaml.Marshal(conf)
		if err != nil {
			logger.Error(err, "failed to marshal FluentBit config", "key", key)
			continue
		}
		cmData[key] = string(data)
	}

	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      flb.Name + "-config",
			Namespace: flb.Namespace,
		},
	}

	_, err = ctrl.CreateOrUpdate(ctx, r.Client, &cm, func() error {
		if err := ctrl.SetControllerReference(&flb, &cm, r.Scheme); err != nil {
			return err
		}

		cm.Data = cmData
		return nil
	})
	if err != nil {
		logger.Error(err, "failed to create or update FluentBit ConfigMap")
		return ctrl.Result{}, err
	}

	newDs := buildFluentBitDaemonSet(&flb)
	ds := appsv1.DaemonSet{
		ObjectMeta: newDs.ObjectMeta,
	}

	_, err = ctrl.CreateOrUpdate(ctx, r.Client, &ds, func() error {
		if err := ctrl.SetControllerReference(&flb, &ds, r.Scheme); err != nil {
			return err
		}

		ds = *newDs
		return nil
	})
	if err != nil {
		logger.Error(err, "failed to create or update FluentBit DaemonSet")
		return ctrl.Result{}, err
	}

	logger.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

func buildFluentBitDaemonSet(fb *v1alpha1.FluentBit) *appsv1.DaemonSet {
	dsLabels := map[string]string{
		"app": fb.Name,
	}

	ds := appsv1.DaemonSet{
		ObjectMeta: fb.ObjectMeta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: dsLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: dsLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "fluent-bit",
							Image: fb.Spec.Image,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      fb.Name + "-config",
									MountPath: "/fluent-bit/etc/",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: fb.Name + "-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fb.Name + "-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &ds
}

func buildFluentBitPipelineConfig(lp *v1alpha1.LoggingPipeline) *fluentbit.FluentBitConfig {
	tag := lp.GenerateTag()
	// TODO: shared opts for all plugins?
	inputOpts := input.BuildOptions{
		Tag:       tag,
		Namespace: lp.Namespace,
	}
	outputOpts := output.BuildOptions{
		Tag:       tag,
		Namespace: lp.Namespace,
	}

	inputs := lp.Spec.Inputs.ToFluentBitInputs(&inputOpts)
	outputs := lp.Spec.Outputs.ToFluentBitOutputs(&outputOpts)

	pipeline := fluentbit.FluentBitConfigPipelineSpec{}
	for _, in := range inputs {
		pipeline.Inputs = append(pipeline.Inputs, fluentbit.FluentBitPipelineInputSpec{
			Name:  in.InputName(),
			Input: in,
		})
	}

	for _, out := range outputs {
		pipeline.Outputs = append(pipeline.Outputs, fluentbit.FluentBitPipelineOutputSpec{
			Name:   out.OutputName(),
			Output: out,
		})
	}

	return &fluentbit.FluentBitConfig{
		Pipeline: &pipeline,
	}
}

func buildMainFluentBitConfig(flb *v1alpha1.FluentBit) *fluentbit.FluentBitConfig {
	conf := fluentbit.FluentBitConfig{
		Includes: []string{
			"*.yaml",
		},
	}
	if flb.Spec.Config.Flush != nil {
		conf.Service.Flush = *flb.Spec.Config.Flush
	}

	if flb.Spec.Config.Grace != nil {
		conf.Service.Grace = *flb.Spec.Config.Grace
	}

	if flb.Spec.Config.LogLevel != nil {
		conf.Service.LogLevel = *flb.Spec.Config.LogLevel
	}

	return &conf
}

// SetupWithManager sets up the controller with the Manager.
func (r *FluentBitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.FluentBit{}).
		Owns(&corev1.ConfigMap{}).
		Watches(
			&v1alpha1.LoggingPipeline{},
			handler.EnqueueRequestsFromMapFunc(r.findFluentBitsForLoggingPipeline),
		).
		Named("fluentbit").
		Complete(r)
}

// findFluentBitsForLoggingPipeline finds FluentBit resources that match a LoggingPipeline's labels.
func (r *FluentBitReconciler) findFluentBitsForLoggingPipeline(ctx context.Context, obj client.Object) []reconcile.Request {
	pipeline, ok := obj.(*v1alpha1.LoggingPipeline)
	if !ok {
		return nil
	}

	fluentBitList := &v1alpha1.FluentBitList{}
	if err := r.List(ctx, fluentBitList); err != nil {
		return nil
	}

	var requests []reconcile.Request
	for _, fb := range fluentBitList.Items {
		selectors, err := metav1.LabelSelectorAsSelector(&fb.Spec.PipelineSelector)
		// TODO: error handling
		if err != nil {
			continue
		}
		if selectors.Matches(labels.Set(pipeline.Labels)) {
			requests = append(requests, reconcile.Request{
				NamespacedName: client.ObjectKeyFromObject(&fb),
			})
		}
	}

	return requests
}
