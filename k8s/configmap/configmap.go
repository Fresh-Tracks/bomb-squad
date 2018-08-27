package configmap

import (
	"context"
	"log"
	"time"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	promcfg "github.com/prometheus/prometheus/config"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	kcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// ConfigMap Struct to hold relevant details of a ConfigMap
type ConfigMap struct {
	Client      *k8s.Client
	Name        string
	CM          *corev1.ConfigMap
	Key         string
	LastUpdated time.Duration
	Ctx         context.Context
}

// By most accounts, this type of error can be ignored. Let's hope that's true.
const optimisticMergeConflictError = "kubernetes api: Failure 409 Operation cannot be fulfilled on configmaps \"prometheus\": the object has been modified; please apply your changes to the latest version and try again"

// Init just tries to get a K8s client created, and if it can't, bail
func (c *ConfigMap) Init(ctx context.Context) {
	cm := corev1.ConfigMap{}
	// NewInClusterClient() creates a client that is forced into the same namespace
	// as the entity in which the Client is created. So as long as bomb-squad runs
	// in a container that resides in the same namespace as the Prometheus ConfigMap,
	// we're good
	client, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatal(err)
	}
	c.Client = client

	err = c.Client.Get(ctx, c.Client.Namespace, c.Name, &cm)
	if err != nil {
		log.Fatal(err)
	}
	c.CM = &cm
}

// ReadRawData pulls in value from the `data` key in the ConfigMap as-is
func (c *ConfigMap) ReadRawData(ctx context.Context, key string) []byte {
	cm := corev1.ConfigMap{}
	err := c.Client.Get(ctx, c.Client.Namespace, c.Name, &cm)
	if err != nil {
		log.Fatal(err)
	}
	c.CM = &cm

	if res, ok := c.CM.Data[key]; !ok {
		return []byte{}
	} else {
		return []byte(res)
	}
}

func (c *ConfigMap) Update(ctx context.Context, cfg promcfg.Config) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	c.CM.Data[c.Key] = string(b)
	err = c.UpdateWithRetries(5)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c *ConfigMap) UpdateWithRetries(retries int) error {
	var err error

	for tries := 1; tries <= retries; tries++ {
		if err = c.Client.Update(c.Ctx, c.CM); err != nil && err.Error() != optimisticMergeConflictError {
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}
	return err
}

// Begin proper k8s bits
// ConfigMapWrapper is a struct with public fields, which implements github.com/Fresh-Tracks/bomb-squad/config.Configurator
type ConfigMapWrapper struct {
	// ConfigMapInterface is a client, not the ConfigMap itself
	Client kcorev1.ConfigMapInterface
	Name   string
}

// NewConfigMapWrapper returns a ConfigMapWrapper
func NewConfigMapWrapper(client kubernetes.Interface, namespace string, configMapName string) *ConfigMapWrapper {
	return &ConfigMapWrapper{
		Client: client.CoreV1().ConfigMaps(namespace),
		Name:   configMapName,
	}
}

// Read implements github.com/Fresh-Tracks/bomb-squad/config.Configurator
func (c *ConfigMapWrapper) Read() promcfg.Config {
	return promcfg.Config{}
}

// Write implements github.com/Fresh-Tracks/bomb-squad/config.Configurator
func (c *ConfigMapWrapper) Write([]byte) error {
	return nil
}
