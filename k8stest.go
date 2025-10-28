package k8stest

import (
	"context"
	"fmt"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sinternal "github.com/tom1299/k8stest/internal"
)

type Resources struct {
	deployments  []appsv1.Deployment
	statefulSets []appsv1.StatefulSet
	configMaps   []corev1.ConfigMap
	secrets      []corev1.Secret
	Options      []ResourceOption
	TestClients  *TestClients
	Ctx          context.Context
	Timeout      time.Duration
}

// New creates a new Resources object with the given TestClients and context.
// It initializes the Timeout to the default value of 30 seconds.
func New(t *testing.T, ctx context.Context) *Resources {
	return &Resources{
		TestClients: SetupTestClients(t),
		Ctx:         ctx,
		Timeout:     30 * time.Second,
	}
}

type Deployment struct {
	Resources
}

type StatefulSet struct {
	Resources
}

type ResourceOption func(runtime.Object)

func boolPtr(b bool) *bool {
	return &b
}

func createPodTemplateSpec(name string) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
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
	}
}

func attachSecretVolume(podSpec *corev1.PodSpec, secretName string) {
	podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
		Name: "secret-" + secretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	})

	podSpec.Containers[0].VolumeMounts = append(podSpec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      "secret-" + secretName,
		MountPath: "/etc/secret",
	})
}

func attachConfigMapVolume(podSpec *corev1.PodSpec, configMapName string) {
	podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
		Name: "config-map-" + configMapName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMapName,
				},
				Optional: boolPtr(false),
			},
		},
	})

	podSpec.Containers[0].VolumeMounts = append(podSpec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      "config-map-" + configMapName,
		MountPath: "/etc/config",
	})
}

func (r *Resources) Create() (*Resources, error) {
	for _, configMap := range r.configMaps {
		_, err := r.TestClients.ClientSet.CoreV1().ConfigMaps("default").Create(r.Ctx, &configMap, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create configmap: %w", err)
		}
	}
	for _, secret := range r.secrets {
		_, err := r.TestClients.ClientSet.CoreV1().Secrets("default").Create(r.Ctx, &secret, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create secret: %w", err)
		}
	}

	for _, deployment := range r.deployments {
		_, err := r.TestClients.ClientSet.AppsV1().Deployments("default").Create(r.Ctx, &deployment, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create deployment: %w", err)
		}
	}

	for _, statefulSet := range r.statefulSets {
		_, err := r.TestClients.ClientSet.AppsV1().StatefulSets("default").Create(r.Ctx, &statefulSet, metav1.CreateOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create statefulset: %w", err)
		}
	}

	return r, nil
}

func (r *Resources) Wait() error {
	for _, deployment := range r.deployments {
		err := wait.PollUntilContextTimeout(r.Ctx, 100*time.Millisecond, r.Timeout, true, func(ctx context.Context) (bool, error) {
			dep, err := r.TestClients.ClientSet.AppsV1().Deployments("default").Get(ctx, deployment.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}

			return dep.Status.AvailableReplicas == *dep.Spec.Replicas, nil
		})
		if err != nil {
			return fmt.Errorf("failed to wait for deployment %s: %w", deployment.Name, err)
		}
	}

	for _, statefulSet := range r.statefulSets {
		err := wait.PollUntilContextTimeout(r.Ctx, 100*time.Millisecond, r.Timeout, true, func(ctx context.Context) (bool, error) {
			sts, err := r.TestClients.ClientSet.AppsV1().StatefulSets("default").Get(ctx, statefulSet.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}

			return sts.Spec.Replicas != nil && sts.Status.ReadyReplicas == *sts.Spec.Replicas, nil
		})
		if err != nil {
			return fmt.Errorf("failed to wait for statefulset %s: %w", statefulSet.Name, err)
		}
	}

	return nil
}

type deleteFunc func(ctx context.Context, name string, opts metav1.DeleteOptions) error

func deleteResource(ctx context.Context, name, resourceType string, deleter deleteFunc) error {
	err := deleter(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete %s: %w", resourceType, err)
	}

	return nil
}

func (r *Resources) Delete() error {
	for _, statefulSet := range r.statefulSets {
		if err := deleteResource(r.Ctx, statefulSet.Name, "statefulset",
			r.TestClients.ClientSet.AppsV1().StatefulSets("default").Delete); err != nil {
			return err
		}
	}

	for _, deployment := range r.deployments {
		if err := deleteResource(r.Ctx, deployment.Name, "deployment",
			r.TestClients.ClientSet.AppsV1().Deployments("default").Delete); err != nil {
			return err
		}
	}

	for _, secret := range r.secrets {
		if err := deleteResource(r.Ctx, secret.Name, "secret",
			r.TestClients.ClientSet.CoreV1().Secrets("default").Delete); err != nil {
			return err
		}
	}

	for _, configMap := range r.configMaps {
		if err := deleteResource(r.Ctx, configMap.Name, "configmap",
			r.TestClients.ClientSet.CoreV1().ConfigMaps("default").Delete); err != nil {
			return err
		}
	}

	return nil
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

	r.ApplyOptions(&r.secrets[len(r.secrets)-1])

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

	r.ApplyOptions(&r.configMaps[len(r.configMaps)-1])

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
			Template: createPodTemplateSpec(name),
		},
	})

	r.ApplyOptions(&r.deployments[len(r.deployments)-1])

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
			Template: createPodTemplateSpec(name),
		},
	})

	r.ApplyOptions(&r.statefulSets[len(r.statefulSets)-1])

	return &StatefulSet{*r}
}

func (r *Resources) And() *Resources {
	return r
}

// WithTimeout sets the Timeout for Wait operations.
func (r *Resources) WithTimeout(timeout time.Duration) *Resources {
	r.Timeout = timeout
	return r
}

func (d *Deployment) WithSecret(name string) *Deployment {
	resources := &d.Resources
	resources.WithSecret(name)

	deployment := &d.deployments[len(d.deployments)-1]
	attachSecretVolume(&deployment.Spec.Template.Spec, name)

	return d
}

func (d *Deployment) WithConfigMap(name string) *Deployment {
	resources := &d.Resources
	resources.WithConfigMap(name)

	deployment := &d.deployments[len(d.deployments)-1]
	attachConfigMapVolume(&deployment.Spec.Template.Spec, name)

	return d
}

func (d *Deployment) And() *Resources {
	return &d.Resources
}

func (s *StatefulSet) WithSecret(name string) *StatefulSet {
	resources := &s.Resources
	resources.WithSecret(name)

	statefulSet := &s.statefulSets[len(s.statefulSets)-1]
	attachSecretVolume(&statefulSet.Spec.Template.Spec, name)

	return s
}

func (s *StatefulSet) WithConfigMap(name string) *StatefulSet {
	resources := &s.Resources
	resources.WithConfigMap(name)

	statefulSet := &s.statefulSets[len(s.statefulSets)-1]
	attachConfigMapVolume(&statefulSet.Spec.Template.Spec, name)

	return s
}

func (s *StatefulSet) And() *Resources {
	return &s.Resources
}

func (r *Resources) ApplyOptions(object runtime.Object) {
	for _, option := range r.Options {
		option(object)
	}
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
