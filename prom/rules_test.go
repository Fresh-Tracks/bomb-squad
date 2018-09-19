package prom_test

import (
	"testing"

	"github.com/Fresh-Tracks/bomb-squad/bstesting"
	"github.com/Fresh-Tracks/bomb-squad/config"
	"github.com/stretchr/testify/require"
)

// TODO: The AppendRulesFile function in the prom package is kinda screwy, as it's specific to
// ConfigMaps. Should refactor it to use Configurators, then do this test properly
func TestCanAppendRulesFile(t *testing.T) {
	c := bstesting.NewConfigurator(t)
	promcfg, err := config.ReadPromConfig(c)
	require.NoError(t, err)
	//require.NotEmpty(t, promcfg)
}
