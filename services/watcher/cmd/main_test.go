package main

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandConfig(t *testing.T) {
	cmd.Run = func (cmd *cobra.Command, args []string) {

	}
	configureFlags()
	assert.Contains(t, cmd.Short, "Zaidan")
	assert.Contains(t, cmd.Short, "Watcher")

	flags := cmd.PersistentFlags()
	assert.NotNil(t, flags.Lookup("geth"))
	assert.Equal(t, flags.Lookup("geth").Value.String(), "https://ropsten.infura.io")
	assert.NotNil(t, flags.Lookup("port"))
	assert.Equal(t, flags.Lookup("port").Value.String(), "5001")
	assert.NotNil(t, flags.ShorthandLookup("p"))
	// TODO how to test?
}

