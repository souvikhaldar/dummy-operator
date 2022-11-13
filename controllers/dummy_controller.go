/*
Copyright 2022.

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
	"os"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	souvikhaldarinv1alpha1 "github.com/souvikhaldar/dummy-operator/api/v1alpha1"
)

// DummyReconciler reconciles a Dummy object
type DummyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=souvikhaldar.in,resources=dummies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=souvikhaldar.in,resources=dummies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=souvikhaldar.in,resources=dummies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Dummy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *DummyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO(user): your logic here
	log := log.FromContext(ctx)
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	dummy := &souvikhaldarinv1alpha1.Dummy{}
	err := r.Get(ctx, req.NamespacedName, dummy)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Dummy resource not found")
			return ctrl.Result{}, nil
		}
		log.Info("Failed to get dummy", err)
		return ctrl.Result{}, err
	}

	// copy the value of `message` of `spec` into `specEcho` of `status`
	dummy.Status.SpecEcho = dummy.Spec.Message

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(
		ctx,
		types.NamespacedName{
			Name:      dummy.Name,
			Namespace: dummy.Namespace,
		},
		found,
	)
	if err != nil && apierrors.IsNotFound(err) {
		// define a new deployment
		dep, err := r.deploymentForDummy(dummy)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for dummy")

			return ctrl.Result{}, err
		}
		log.Info("Creating a new Deployment",
			"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create new Deployment",
				"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully
		// We will requeue the reconciliation so that we can ensure the state
		// and move forward for the next operations
		return ctrl.Result{RequeueAfter: time.Minute}, nil

	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DummyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&souvikhaldarinv1alpha1.Dummy{}).
		Complete(r)
}

// deploymentForDummy returns a Memcached Deployment object
func (r *DummyReconciler) deploymentForDummy(
	dummy *souvikhaldarinv1alpha1.Dummy) (*appsv1.Deployment, error) {

	// Get the Operand image
	//image, err := imageForDummy()
	//if err != nil {
	//	return nil, err
	//}

	// TODO: Remove hardcoding
	image := "nginx:alpine"

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummy.Name,
			Namespace: dummy.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &dummy.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "dummy-nginx-server"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "dummy-nginx-server"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  "dummy-nginx",
						Ports: []corev1.ContainerPort{{
							ContainerPort: dummy.Spec.Port,
							Name:          "http",
							Protocol:      "TCP",
						}},
					}},
				},
			},
		},
	}
	// Set the ownerRef for the Deployment
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
	if err := ctrl.SetControllerReference(dummy, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// imageForMemcached gets the Operand image which is managed by this controller
// from the DUMMY_IMAGE environment variable defined in the config/manager/manager.yaml
func imageForDummy() (string, error) {
	var imageEnvVar = "DUMMY_IMAGE"
	image, found := os.LookupEnv(imageEnvVar)
	if !found {
		return "", fmt.Errorf("Unable to find %s environment variable with the image", imageEnvVar)
	}
	return image, nil
}
