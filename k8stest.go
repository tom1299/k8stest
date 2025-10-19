package k8stest

import (
	"context"
	"errors"
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Resources struct {
	deployments []appsv1.Deployment
	configMaps  []corev1.ConfigMap
	secrets     []corev1.Secret
}

type Deployment struct {
	Resources
}

type RestartRule struct {
	Resources
}

func boolPtr(b bool) *bool {
	return &b
}

func (r *Resources) Create(client *TestClients, ctx context.Context) (*Resources, error) {
	for _, configMap := range r.configMaps {
		_, err := client.clientSet.CoreV1().ConfigMaps("default").Create(ctx, &configMap, metav1.CreateOptions{})
		if err != nil {
			return nil, errors.New("failed to create configmap: " + err.Error())
		}
	}
	for _, secret := range r.secrets {
		_, err := client.clientSet.CoreV1().Secrets("default").Create(ctx, &secret, metav1.CreateOptions{})
		if err != nil {
			return nil, errors.New("failed to create secret: " + err.Error())
		}
	}

	for _, deployment := range r.deployments {
		_, err := client.clientSet.AppsV1().Deployments("default").Create(ctx, &deployment, metav1.CreateOptions{})
		if err != nil {
			return nil, errors.New("failed to create deployment: " + err.Error())
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

func (r *Resources) And() *Resources {
	return r
}

func (d *Deployment) WithSecret(name string) *Deployment {
	resources := &d.Resources
	resources.WithSecret(name)

	deployment := &d.deployments[len(d.deployments)-1]
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: fmt.Sprintf("secret-%s", name),
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: name,
			},
		},
	})

	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.
		Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      fmt.Sprintf("secret-%s", name),
		MountPath: "/etc/secret",
	})

	return d
}

func (d *Deployment) WithConfigMap(name string) *Deployment {
	resources := &d.Resources
	resources.WithConfigMap(name)

	deployment := &d.deployments[len(d.deployments)-1]
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: fmt.Sprintf("config-map-%s", name),
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
		Name:      fmt.Sprintf("config-map-%s", name),
		MountPath: "/etc/config",
	})

	return d
}

func (d *Deployment) And() *Resources {
	return &d.Resources
}

type TestClients struct {
	clientSet *kubernetes.Clientset
	k8sClient client.Client
}

func setupTestClients(t *testing.T) *TestClients {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	cfg, err := kubeConfig.ClientConfig()
	if err != nil {
		t.Fatalf("Failed to get kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create kubernetes clientSet: %v", err)
	}

	scheme := setupScheme()
	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		t.Fatalf("Failed to create controller-runtime client: %v", err)
	}

	return &TestClients{
		clientSet: clientset,
		k8sClient: k8sClient,
	}
}

func setupScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	return scheme
}
