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
	assert.NotNil(t, flags.Lookup("eth"))
	assert.Equal(t, flags.Lookup("eth").Value.String(), "wss://ropsten.infura.io/ws")
	assert.NotNil(t, flags.Lookup("port"))
	assert.Equal(t, flags.Lookup("port").Value.String(), "5001")
	assert.NotNil(t, flags.ShorthandLookup("p"))
}

