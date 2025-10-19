package k8stest

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildClients creates a Kubernetes clientset and a controller-runtime client
// using the current kubeconfig context. It encapsulates the setup logic used by
// tests, returning the constructed clients or an error.
func BuildClients() (*kubernetes.Clientset, ctrclient.Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	cfg, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	scheme := runtime.NewScheme()
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)

	k8sClient, err := ctrclient.New(cfg, ctrclient.Options{Scheme: scheme})
	if err != nil {
		return nil, nil, err
	}

	return clientset, k8sClient, nil
}
