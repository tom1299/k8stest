package k8stest

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestFluent(t *testing.T) {
	testData := Resources{}
	_, err := testData.WithDeployment("deployment-1").
		WithConfigMap("config-map-1").
		WithSecret("secret-1").
		Create(SetupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	testData := Resources{}
	resources, err := testData.WithDeployment("deployment-delete-1").
		WithConfigMap("config-map-delete-1").
		WithSecret("secret-delete-1").
		Create(SetupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Delete(SetupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestFluentStatefulSet(t *testing.T) {
	testData := Resources{}
	_, err := testData.WithStatefulSet("statefulset-1").
		WithConfigMap("config-map-2").
		WithSecret("secret-2").
		Create(SetupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteStatefulSet(t *testing.T) {
	testData := Resources{}
	resources, err := testData.WithStatefulSet("statefulset-delete-1").
		WithConfigMap("config-map-delete-1").
		WithSecret("secret-delete-1").
		Create(SetupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Delete(SetupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}
}
func TestDeleteNonExistent(t *testing.T) {
	testData := Resources{}
	resources := testData.WithDeployment("non-existent-deployment").
		WithConfigMap("non-existent-configmap").
		WithSecret("non-existent-secret")

	_, err := resources.Delete(SetupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestDeploymentWithMutator(t *testing.T) {
	addAnnotationOption := func(obj runtime.Object) {
		d, ok := obj.(*appsv1.Deployment)
		if !ok {
			return
		}
		if d.Annotations == nil {
			d.Annotations = make(map[string]string)
		}
		d.Annotations["test-annotation"] = "test-value"
	}

	testData := Resources{
		mutators: []ResourceOption{addAnnotationOption},
	}

	clients := SetupTestClients(t)
	resources, err := testData.
		WithDeployment("deployment-with-annotation").
		Create(clients, context.Background())
	if err != nil {
		t.Error(err)
	}

	deployment, err := clients.ClientSet.AppsV1().Deployments("default").Get(
		context.Background(),
		"deployment-with-annotation",
		metav1.GetOptions{},
	)
	if err != nil {
		t.Errorf("Failed to get deployment from cluster: %v", err)
	}

	if deployment.Annotations["test-annotation"] != "test-value" {
		t.Errorf("Expected annotation 'test-annotation' with value 'test-value', got %v",
			deployment.Annotations)
	}

	_, err = resources.Delete(clients, context.Background())
	if err != nil {
		t.Error(err)
	}
}
