package k8stest

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ZeroTerminationGracePeriodOption returns a ResourceOption that sets the
// terminationGracePeriodSeconds to 0 on the PodSpec of Deployments and StatefulSets.
//
//nolint:gocritic // dupBranchBody: identical branches are intentional
func ZeroTerminationGracePeriodOption() ResourceOption {
	return func(obj runtime.Object) {
		switch o := obj.(type) {
		case *appsv1.Deployment:
			zero := int64(0)
			if o.Spec.Template.Spec.TerminationGracePeriodSeconds == nil {
				o.Spec.Template.Spec.TerminationGracePeriodSeconds = &zero
			} else {
				o.Spec.Template.Spec.TerminationGracePeriodSeconds = &zero
			}
		case *appsv1.StatefulSet:
			zero := int64(0)
			if o.Spec.Template.Spec.TerminationGracePeriodSeconds == nil {
				o.Spec.Template.Spec.TerminationGracePeriodSeconds = &zero
			} else {
				o.Spec.Template.Spec.TerminationGracePeriodSeconds = &zero
			}
		}
	}
}

func TestFluent(t *testing.T) {

	resources, err := New(t, context.Background()).
		WithResourceOption(ZeroTerminationGracePeriodOption()).
		WithDeployment("deployment-1").
		WithConfigMap("config-map-1").
		WithSecret("secret-1").
		Create()
	if err != nil {
		t.Error(err)
	}

	err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	resources, err := New(t, context.Background()).WithResourceOption(ZeroTerminationGracePeriodOption()).
		WithDeployment("deployment-delete-1").
		WithConfigMap("config-map-delete-1").
		WithSecret("secret-delete-1").
		Create()
	if err != nil {
		t.Error(err)
	}

	err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestFluentStatefulSet(t *testing.T) {
	resources, err := New(t, context.Background()).WithResourceOption(ZeroTerminationGracePeriodOption()).
		WithStatefulSet("statefulset-1").
		WithConfigMap("config-map-2").
		WithSecret("secret-2").
		Create()
	if err != nil {
		t.Error(err)
	}

	err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteStatefulSet(t *testing.T) {
	resources, err := New(t, context.Background()).WithStatefulSet("statefulset-delete-1").
		WithConfigMap("config-map-delete-1").
		WithSecret("secret-delete-1").
		WithResourceOption(ZeroTerminationGracePeriodOption()).
		Create()
	if err != nil {
		t.Error(err)
	}

	err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}
func TestDeleteNonExistent(t *testing.T) {
	resources := New(t, context.Background()).WithDeployment("non-existent-deployment").
		WithConfigMap("non-existent-configmap").
		WithSecret("non-existent-secret")

	err := resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestDeploymentWithOption(t *testing.T) {
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

	resources, err := New(t, context.Background()).WithResourceOption(ZeroTerminationGracePeriodOption()).
		WithResourceOption(addAnnotationOption).
		WithDeployment("deployment-with-annotation").
		Create()
	if err != nil {
		t.Error(err)
	}

	err = resources.Wait()
	if err != nil {
		t.Error(err)
	}

	deployment, err := resources.TestClients.ClientSet.AppsV1().Deployments("default").Get(
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

	err = resources.Delete()
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

	resources, err := New(t, context.Background()).WithResourceOption(setInvalidImageOption).
		WithDeployment("deployment-with-invalid-image").
		Create()
	if err != nil {
		t.Error(err)
	}

	waitForResourcesCh := make(chan error, 1)
	go func() {
		waitForResourcesCh <- resources.Wait(5 * time.Second)
	}()

	pollResultCh := make(chan bool, 1)

	pollOnce := func() {
		go func() {
			podList, err := resources.TestClients.ClientSet.CoreV1().Pods("default").List(
				context.Background(),
				metav1.ListOptions{
					LabelSelector: "app=deployment-with-invalid-image",
				},
			)

			if err != nil {
				t.Logf("Failed to list pods while polling: %v", err)
				pollResultCh <- false
				return
			}

			cs := podList.Items[0].Status.ContainerStatuses
			if len(cs) == 0 {
				pollResultCh <- false
				return
			}

			pollResultCh <- cs[0].State.Waiting != nil && cs[0].State.Waiting.Reason == "ErrImagePull"
		}()
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			pollOnce()
		case found := <-pollResultCh:
			if found {
				// Pod reached ErrImagePull first => Expected result
				break loop
			}
		case err = <-waitForResourcesCh:
			if err != nil {
				t.Errorf("Unexpected error from Wait: %v", err)
			}
			// Wait should not finish before the Pod has reached ErrImagePull => Error
			t.Error("Expected pod to reach ErrImagePull before Wait finished")
			break loop
		}
	}

	err = resources.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestConfigurableTimeout(t *testing.T) {
	tests := []struct {
		name           string
		deplomentName  string
		timeout        time.Duration
		minExpectedDur time.Duration
		maxExpectedDur time.Duration
	}{
		{
			name:           "Timeout 0 seconds",
			deplomentName:  "deployment-timeout-test-1",
			timeout:        0 * time.Second,
			minExpectedDur: 0 * time.Millisecond,
			maxExpectedDur: 500 * time.Millisecond,
		},
		{
			name:           "Timeout 1 second",
			deplomentName:  "deployment-timeout-test-2",
			timeout:        1 * time.Second,
			minExpectedDur: 900 * time.Millisecond,
			maxExpectedDur: 1500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := New(t, context.Background()).
				WithResourceOption(ZeroTerminationGracePeriodOption()).
				WithDeployment(tt.deplomentName).
				Create()
			if err != nil {
				t.Error(err)
			}

			// Measure time taken for Wait to fail
			start := time.Now()
			err = resources.Wait(tt.timeout)
			duration := time.Since(start)

			// Wait should fail because deployment can't become ready that fast
			if err == nil {
				t.Error("Expected Wait to fail due to timeout, but it succeeded")
			}

			// Verify the Timeout duration is within expected range
			if duration < tt.minExpectedDur || duration > tt.maxExpectedDur {
				t.Errorf("Expected duration between %v and %v, got %v",
					tt.minExpectedDur, tt.maxExpectedDur, duration)
			}

			err = resources.Delete()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
