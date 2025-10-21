package k8stest

import (
	"context"
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sinternal "github.com/tom1299/k8stest/internal"
)

type Resources struct {
	deployments  []appsv1.Deployment
	statefulSets []appsv1.StatefulSet
	configMaps   []corev1.ConfigMap
	secrets      []corev1.Secret
}

type Deployment struct {
	Resources
}

type StatefulSet struct {
	Resources
}

func boolPtr(b bool) *bool {
	return &b
}

func (r *Resources) Create(testClients *TestClients, ctx context.Context) (*Resources, error) {
	for _, configMap := range r.configMaps {
		_, err := testClients.ClientSet.CoreV1().ConfigMaps("default").Create(ctx, &configMap, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create configmap: %w", err)
		}
	}
	for _, secret := range r.secrets {
		_, err := testClients.ClientSet.CoreV1().Secrets("default").Create(ctx, &secret, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create secret: %w", err)
		}
	}

	for _, deployment := range r.deployments {
		_, err := testClients.ClientSet.AppsV1().Deployments("default").Create(ctx, &deployment, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create deployment: %w", err)
		}
	}

	for _, statefulSet := range r.statefulSets {
		_, err := testClients.ClientSet.AppsV1().StatefulSets("default").Create(ctx, &statefulSet, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create statefulset: %w", err)
		}
	}

	return r, nil
}

func (r *Resources) Delete(testClients *TestClients, ctx context.Context) (*Resources, error) {
	for _, statefulSet := range r.statefulSets {
		err := testClients.ClientSet.AppsV1().StatefulSets("default").Delete(ctx, statefulSet.Name, metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to delete statefulset: %w", err)
		}
	}

	for _, deployment := range r.deployments {
		err := testClients.ClientSet.AppsV1().Deployments("default").Delete(ctx, deployment.Name, metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to delete deployment: %w", err)
		}
	}

	for _, secret := range r.secrets {
		err := testClients.ClientSet.CoreV1().Secrets("default").Delete(ctx, secret.Name, metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to delete secret: %w", err)
		}
	}

	for _, configMap := range r.configMaps {
		err := testClients.ClientSet.CoreV1().ConfigMaps("default").Delete(ctx, configMap.Name, metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to delete configmap: %w", err)
		}
	}

	return r, nil
}

func (r *Resources) WithSecret(name string) *Resources {
	r.secrets = append(r.secrets, corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "appsv1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string][]byte{
			"key": []byte("value"),
		},
	})

	return r
}

func (r *Resources) WithConfigMap(name string) *Resources {
	r.configMaps = append(r.configMaps, corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Data: map[string]string{
			"key": "value",
		},
	})

	return r
}

func (r *Resources) WithDeployment(name string) *Deployment {
	r.deployments = append(r.deployments, appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/appsv1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container-1",
							Image: "busybox:latest",
							Command: []string{
								"sleep",
								"3600",
							},
						},
					},
				},
			},
		},
	})

	return &Deployment{*r}
}

func (r *Resources) WithStatefulSet(name string) *StatefulSet {
	r.statefulSets = append(r.statefulSets, appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/appsv1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container-1",
							Image: "busybox:latest",
							Command: []string{
								"sleep",
								"3600",
							},
						},
					},
				},
			},
		},
	})

	return &StatefulSet{*r}
}

func (r *Resources) And() *Resources {
	return r
}

func (d *Deployment) WithSecret(name string) *Deployment {
	resources := &d.Resources
	resources.WithSecret(name)

	deployment := &d.deployments[len(d.deployments)-1]
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: "secret-" + name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: name,
			},
		},
	})

	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.
		Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      "secret-" + name,
		MountPath: "/etc/secret",
	})

	return d
}

func (d *Deployment) WithConfigMap(name string) *Deployment {
	resources := &d.Resources
	resources.WithConfigMap(name)

	deployment := &d.deployments[len(d.deployments)-1]
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: "config-map-" + name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
				Optional: boolPtr(false),
			},
		},
	})

	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.
		Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      "config-map-" + name,
		MountPath: "/etc/config",
	})

	return d
}

func (d *Deployment) And() *Resources {
	return &d.Resources
}

func (s *StatefulSet) WithSecret(name string) *StatefulSet {
	resources := &s.Resources
	resources.WithSecret(name)

	statefulSet := &s.statefulSets[len(s.statefulSets)-1]
	statefulSet.Spec.Template.Spec.Volumes = append(statefulSet.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: "secret-" + name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: name,
			},
		},
	})

	statefulSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulSet.Spec.Template.Spec.
		Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      "secret-" + name,
		MountPath: "/etc/secret",
	})

	return s
}

func (s *StatefulSet) WithConfigMap(name string) *StatefulSet {
	resources := &s.Resources
	resources.WithConfigMap(name)

	statefulSet := &s.statefulSets[len(s.statefulSets)-1]
	statefulSet.Spec.Template.Spec.Volumes = append(statefulSet.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: "config-map-" + name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
				Optional: boolPtr(false),
			},
		},
	})

	statefulSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(statefulSet.Spec.Template.Spec.
		Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      "config-map-" + name,
		MountPath: "/etc/config",
	})

	return s
}

func (s *StatefulSet) And() *Resources {
	return &s.Resources
}

type TestClients struct {
	ClientSet *kubernetes.Clientset
	K8sClient client.Client
}

func SetupTestClients(t *testing.T) *TestClients {
	clientset, k8sClient, err := k8sinternal.BuildClients()
	if err != nil {
		t.Fatalf("Failed to set up Kubernetes clients: %v", err)
	}

	return &TestClients{
		ClientSet: clientset,
		K8sClient: k8sClient,
	}
}

func SetupScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	return scheme
}
