package k8stest

import (
	"context"
	"testing"
)

func TestFluent(t *testing.T) {
	testData := Resources{}
	_, err := testData.WithDeployment("deployment-1").
		WithConfigMap("config-map-1").
		WithSecret("secret-1").
		Create(setupTestClients(t), context.Background())
	if err != nil {
		t.Error(err)
	}
}
