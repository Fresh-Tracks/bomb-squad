package prom

import (
	"context"
	"log"

	"github.com/Fresh-Tracks/bomb-squad/configmap"
	promcfg "github.com/Fresh-Tracks/bomb-squad/prom/config"
	yaml "gopkg.in/yaml.v2"
)

// GetPrometheusConfig pulls in the full base Prometheus config
// from the provided ConfigMap. Does not include rules nor AM configs.
func GetPrometheusConfig(ctx context.Context, c configmap.ConfigMap) promcfg.Config {
	raw := c.ReadRawData(ctx, c.Key)
	var cfg promcfg.Config
	err := yaml.Unmarshal(raw, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

// AppendRuleFile Appends a static rule file that Bomb Squad needs into the
// array of rule files that may exist in the current Prometheus config
func AppendRuleFile(ctx context.Context, filename string, c configmap.ConfigMap) {
	cfg := GetPrometheusConfig(ctx, c)
	configRuleFiles := cfg.RuleFiles

	newRuleFiles := append(configRuleFiles, filename)
	cfg.RuleFiles = newRuleFiles
	c.Update(ctx, cfg)
}
