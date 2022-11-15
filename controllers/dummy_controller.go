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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

func (r *DummyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	dummy := &souvikhaldarinv1alpha1.Dummy{}

	// Create a deployment for the nginx pod
	existingDummyDeployment := &appsv1.Deployment{}

	// Fetch the dummy instance
	err := r.Get(ctx, req.NamespacedName, dummy)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// The CR is not found. Either it is deleted or not applied
			// to the cluster, hence should stop reconciliation.
			log.Info("Dummy resource not found, stopping reconciliation")

			// but before stoping let's check:
			// if the deployment for dummy exists. If it does
			// it needs to be deleted
			err = r.Get(
				ctx,
				types.NamespacedName{
					Name:      dummy.Name,
					Namespace: dummy.Namespace,
				},
				existingDummyDeployment,
			)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Info("Deployment not found, stoping reconciliation")
					// no need to do anything further
					return ctrl.Result{}, nil
				}

				log.Error(err, "Failed to get the deployment")
				return ctrl.Result{}, err

			}

			log.Info("Deployment exists, need to delete it")
			if err := r.Delete(ctx, existingDummyDeployment); err != nil {
				log.Error(err, "Error in deleting the deployment")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		log.Info("Failed to get dummy instance", err)
		return ctrl.Result{}, err
	}
	// Step 2 of the task:
	// The custom controller must process each Dummy API object simply by logging its name, namespace and
	// the value of spec.message.
	log.Info(
		"STEP 2.",
		"Name:",
		dummy.Name,
		"Namespace:",
		dummy.Namespace,
		"spec.message:",
		dummy.Spec.Message,
	)

	// Task 1.
	// copy the value of `message` of `spec` into `specEcho` of `status`
	dummy.Status.SpecEcho = dummy.Spec.Message

	// Check if the deployment already exists, if not create a new one
	err = r.Get(
		ctx,
		types.NamespacedName{
			Name:      dummy.Name,
			Namespace: dummy.Namespace,
		},
		existingDummyDeployment,
	)
	if err != nil && apierrors.IsNotFound(err) {
		// define a new deployment
		newDummyDeployment, err := r.deploymentForDummy(dummy)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for dummy")

			return ctrl.Result{}, err
		}
		log.Info(
			"Creating a new Deployment",
			"Deployment.Namespace",
			newDummyDeployment.Namespace,
			"Deployment.Name",
			newDummyDeployment.Name,
		)
		if err = r.Create(ctx, newDummyDeployment); err != nil {
			log.Error(
				err,
				"Failed to create new Deployment",
				"Deployment.Namespace",
				newDummyDeployment.Namespace,
				"Deployment.Name",
				newDummyDeployment.Name,
			)
			return ctrl.Result{}, err
		}

	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	} else if err == nil {
		// Deployment exists, check if it needs to be updated
		if *existingDummyDeployment.Spec.Replicas != dummy.Spec.ReplicaCount {
			existingDummyDeployment.Spec.Replicas = &dummy.Spec.ReplicaCount
			if err := r.Update(
				ctx,
				existingDummyDeployment,
			); err != nil {
				log.Error(
					err,
					"Failed to updated deployment",
					"Deployment.Namespace",
					existingDummyDeployment.Namespace,
					"Deployment.Name",
					existingDummyDeployment.Name,
				)
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *DummyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&souvikhaldarinv1alpha1.Dummy{}).
		Complete(r)
}

// deploymentForDummy returns a dummy Deployment object
func (r *DummyReconciler) deploymentForDummy(
	dummy *souvikhaldarinv1alpha1.Dummy,
) (*appsv1.Deployment, error) {

	// TODO: Remove hardcoding and use env var
	image := "nginx:alpine"
	//image := "ovhplatform/hello:1.0"

	newDummyDeployment := &appsv1.Deployment{
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
							ContainerPort: 80,
							Name:          "http",
							Protocol:      "TCP",
						}},
					}},
				},
			},
		},
	}
	return newDummyDeployment, nil
}
