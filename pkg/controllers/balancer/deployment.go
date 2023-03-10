/*
Copyright 2022 netlops.

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

package balancer

import (
	"context"
	exposerv1alpha1 "github.com/netlops/balancer/api/v1alpha1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcilerBalancer) syncDeployment(balancer *exposerv1alpha1.Balancer) error {
	// firstly, we sync configmap
	cm, err := r.syncConfigMap(balancer)
	if err != nil {
		return err
	}

	// now we sync deployment
	dp, err := NewDeployment(balancer)
	if err != nil {
		return err
	}
	annotations := map[string]string{
		exposerv1alpha1.ConfigMapHashKey: ConfigMapHash(cm),
	}

	// always use the newest annotations.
	dp.Spec.Template.ObjectMeta.Annotations = annotations

	// set balancer as the controller owner-reference of dp
	if err = controllerutil.SetControllerReference(balancer, dp, r.scheme); err != nil {
		return err
	}

	foundDp := &appv1.Deployment{}
	err = r.client.Get(context.Background(), types.NamespacedName{Namespace: balancer.Namespace, Name: balancer.Name}, foundDp)
	if err != nil && errors.IsNotFound(err) {
		// corresponding dp not found in the cluster, create it with the newest dp
		if err = r.client.Create(context.Background(), dp); err != nil {
			return err
		}
		log.Info("Sync deployment", dp.Name, "created")
		return nil
	} else if err != nil {
		return err
	}

	//corresponding dp found, update it with the newest dp
	foundDp.Spec.Template = dp.Spec.Template
	if err = r.client.Update(context.Background(), foundDp); err != nil {
		return err
	}
	log.Info("Sync Deployment", foundDp.Name, "updated")
	return nil
}

// NewDeployment creates a new deployment (which controls one nginx pod)  for the Balancer.
func NewDeployment(balancer *exposerv1alpha1.Balancer) (*appv1.Deployment, error) {
	replicas := int32(1)
	labels := NewPodLabels(balancer)
	nginxContainer := corev1.Container{
		Name:  "nginx",
		Image: "nginx:1.15.9",
		Ports: []corev1.ContainerPort{{ContainerPort: 80}},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      ConfigMapName(balancer),
				MountPath: "/etc/nginx",
				ReadOnly:  true,
			},
		},
	}
	nginxVolume := corev1.Volume{
		Name: ConfigMapName(balancer),
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: ConfigMapName(balancer),
				},
			},
		},
	}
	return &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DeploymentName(balancer),
			Namespace: balancer.Namespace,
			Labels:    labels,
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      DeploymentName(balancer),
					Namespace: balancer.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{nginxContainer},
					Volumes:    []corev1.Volume{nginxVolume},
				},
			},
		},
	}, nil
}

func DeploymentName(balancer *exposerv1alpha1.Balancer) string {
	return balancer.Name + "proxy"
}
