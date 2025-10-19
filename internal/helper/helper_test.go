package helper

import "testing"

func TestValidateNamespace(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		wantErr   bool
	}{
		{"valid namespace", "default", false},
		{"valid with hyphen", "my-namespace", false},
		{"valid with numbers", "namespace123", false},
		{"empty namespace", "", true},
		{"namespace with space", "my namespace", true},
		{"namespace with multiple spaces", "my  test  namespace", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNamespace(tt.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatNamespace(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		want      string
	}{
		{"default namespace", "default", "k8s-namespace:default"},
		{"custom namespace", "my-namespace", "k8s-namespace:my-namespace"},
		{"empty namespace", "", "k8s-namespace:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatNamespace(tt.namespace)
			if got != tt.want {
				t.Errorf("FormatNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}
