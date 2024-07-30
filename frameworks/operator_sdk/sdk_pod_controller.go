/*
Copyright 2024.

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
        "github.com/go-logr/logr"

        "k8s.io/apimachinery/pkg/runtime"
        "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
        "sigs.k8s.io/controller-runtime/pkg/reconcile"
        "sigs.k8s.io/controller-runtime/pkg/client"
        ctrl "sigs.k8s.io/controller-runtime"
        logf "sigs.k8s.io/controller-runtime/pkg/log"
        "sigs.k8s.io/controller-runtime/pkg/log/zap"

        examplev1 "github.com/careena/pod-creator-operator/api/v1"
        corev1 "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
var logger logr.Logger
// PodReconciler reconciles a Pod object
type PodReconciler struct {
        client.Client
        Scheme *runtime.Scheme
        Log logr.Logger
}

//+kubebuilder:rbac:groups=example.example.com,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.example.com,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.example.com,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PodReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

        pod := &examplev1.Pod{}
        if err := r.Get(ctx, req.NamespacedName, pod); err != nil {
                r.Log.Error(err, "unable to fetch Pod")
                return ctrl.Result{}, client.IgnoreNotFound(err)
        }
        existingPod := &corev1.Pod{}
        err := r.Get(ctx, client.ObjectKey{
                Name: pod.Name + "-created",
                Namespace: pod.Namespace,
        }, existingPod)


        if err == nil {
                r.Log.Info("Pod already exists, no action required")
                return ctrl.Result{}, nil
        }

        newPod := &corev1.Pod{
        ObjectMeta: metav1.ObjectMeta{
                        Name: pod.Name + "-created",
                        Namespace: pod.Namespace,
                },
                Spec: corev1.PodSpec{
                        Containers: []corev1.Container{
                                {
                                        Name: "nginx",
                                        Image: "nginx",
                                },
                        },
                },
        }
        if err := controllerutil.SetControllerReference(pod, newPod, r.Scheme); err != nil {
                r.Log.Error(err, "unable to set owner reference on new Pod")
                return reconcile.Result{}, err
        }


        if err := r.Create(ctx, newPod); err != nil {
                r.Log.Error(err, "unable to create Pod")
                return ctrl.Result{}, err
        }

        r.Log.Info("Pod created", "Name", newPod.Name)

        return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
        logger := logf.Log.WithName("pod-controller")
        r.Log = logger

        logf.SetLogger(zap.New(zap.UseDevMode(true)))
        return ctrl.NewControllerManagedBy(mgr).
                For(&examplev1.Pod{}).
                Complete(r)
}
