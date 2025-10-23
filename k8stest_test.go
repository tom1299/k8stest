package k8stest

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestFluent(t *testing.T) {
	resources, err := New(t, context.Background()).WithDeployment("deployment-1").
		WithConfigMap("config-map-1").
		WithSecret("secret-1").
		Create()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	resources, err := New(t, context.Background()).WithDeployment("deployment-delete-1").
		WithConfigMap("config-map-delete-1").
		WithSecret("secret-delete-1").
		Create()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestFluentStatefulSet(t *testing.T) {
	resources, err := New(t, context.Background()).WithStatefulSet("statefulset-1").
		WithConfigMap("config-map-2").
		WithSecret("secret-2").
		Create()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteStatefulSet(t *testing.T) {
	resources, err := New(t, context.Background()).WithStatefulSet("statefulset-delete-1").
		WithConfigMap("config-map-delete-1").
		WithSecret("secret-delete-1").
		Create()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}
func TestDeleteNonExistent(t *testing.T) {
	resources := New(t, context.Background()).WithDeployment("non-existent-deployment").
		WithConfigMap("non-existent-configmap").
		WithSecret("non-existent-secret")

	_, err := resources.Delete()
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

	resources := New(t, context.Background())
	resources.mutators = []ResourceOption{addAnnotationOption}

	resources, err := resources.
		WithDeployment("deployment-with-annotation").
		Create()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	deployment, err := resources.testClients.ClientSet.AppsV1().Deployments("default").Get(
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

	_, err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestDeploymentWithInvalidImage(t *testing.T) {
	setInvalidImageOption := func(obj runtime.Object) {
		d, ok := obj.(*appsv1.Deployment)
		if !ok {
			return
		}
		d.Spec.Template.Spec.Containers[0].Image = "invalid-image-name-that-does-not-exist:latest"
	}

	resources := New(t, context.Background())
	resources.mutators = []ResourceOption{setInvalidImageOption}

	resources, err := resources.
		WithDeployment("deployment-with-invalid-image").
		Create()
	if err != nil {
		t.Error(err)
	}

	// Wait should fail because the deployment cannot become available with invalid image
	_, err = resources.Wait()
	if err == nil {
		t.Error("Expected Wait to fail for deployment with invalid image, but it succeeded")
	}

	_, err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestConfigurableTimeout(t *testing.T) {
	// Test that we can configure a custom timeout
	resources, err := New(t, context.Background()).
		WithTimeout(45 * time.Second).
		WithDeployment("deployment-with-custom-timeout").
		Create()
	if err != nil {
		t.Error(err)
	}

	// Verify the timeout is set correctly
	if resources.timeout != 45*time.Second {
		t.Errorf("Expected timeout to be 45s, got %v", resources.timeout)
	}

	_, err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	_, err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}
