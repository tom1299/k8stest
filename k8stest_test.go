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
