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
	tests := []struct {
		name           string
		timeout        time.Duration
		minExpectedDur time.Duration
		maxExpectedDur time.Duration
	}{
		{
			name:           "timeout 0 seconds",
			timeout:        0 * time.Second,
			minExpectedDur: 0 * time.Millisecond,
			maxExpectedDur: 500 * time.Millisecond,
		},
		{
			name:           "timeout 1 second",
			timeout:        1 * time.Second,
			minExpectedDur: 900 * time.Millisecond,
			maxExpectedDur: 1500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := New(t, context.Background()).
				WithTimeout(tt.timeout).
				WithDeployment("deployment-timeout-test").
				Create()
			if err != nil {
				t.Error(err)
			}

			// Measure time taken for Wait to fail
			start := time.Now()
			_, err = resources.Wait()
			duration := time.Since(start)

			// Wait should fail because deployment can't become ready that fast
			if err == nil {
				t.Error("Expected Wait to fail due to timeout, but it succeeded")
			}

			// Verify the timeout duration is within expected range
			if duration < tt.minExpectedDur || duration > tt.maxExpectedDur {
				t.Errorf("Expected duration between %v and %v, got %v",
					tt.minExpectedDur, tt.maxExpectedDur, duration)
			}

			_, err = resources.Delete()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
