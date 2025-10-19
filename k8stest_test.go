package k8stest

import "testing"

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config == nil {
		t.Fatal("NewConfig returned nil")
	}
	if config.Namespace != "default" {
		t.Errorf("expected Namespace to be 'default', got '%s'", config.Namespace)
	}
	if config.Timeout != 30 {
		t.Errorf("expected Timeout to be 30, got %d", config.Timeout)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		wantErr   bool
	}{
		{"valid namespace", "default", false},
		{"valid namespace with hyphen", "my-namespace", false},
		{"empty namespace", "", true},
		{"namespace with space", "my namespace", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Namespace: tt.namespace}
			err := config.ValidateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFormattedNamespace(t *testing.T) {
	config := &Config{Namespace: "test"}
	formatted := config.GetFormattedNamespace()
	expected := "k8s-namespace:test"
	if formatted != expected {
		t.Errorf("GetFormattedNamespace() = %v, want %v", formatted, expected)
	}
}
