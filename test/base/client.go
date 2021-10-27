package base

import (
	log "github.com/sirupsen/logrus"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset/typed/apiregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func init() {
	apiextensionsv1.AddToScheme(scheme.Scheme)
}

func GetClient() client.Client {
	return GetClientWithConfig(GetConfig())
}

func GetClientWithConfig(config *rest.Config) client.Client {
	c, err := client.New(config, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Fatalf("GetClientWithConfig: %v", err)
	}
	return c
}

func GetConfig() *rest.Config {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("GetConfig: Error getting client config: %v", err)
	}
	return config
}

func GetAPIRegistrationClient() apiregistrationv1.ApiregistrationV1Interface {
	c, err := apiregistrationv1.NewForConfig(GetConfig())
	if err != nil {
		log.Fatalf("Error obtaining API registration client: %v", err)
	}
	return c
}
