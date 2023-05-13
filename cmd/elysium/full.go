package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	cmdnode "github.com/furyaxyz/elysium-node/cmd"
	"github.com/furyaxyz/elysium-node/nodebuilder/core"
	"github.com/furyaxyz/elysium-node/nodebuilder/gateway"
	"github.com/furyaxyz/elysium-node/nodebuilder/header"
	"github.com/furyaxyz/elysium-node/nodebuilder/node"
	"github.com/furyaxyz/elysium-node/nodebuilder/p2p"
	"github.com/furyaxyz/elysium-node/nodebuilder/rpc"
	"github.com/furyaxyz/elysium-node/nodebuilder/state"
)

// NOTE: We should always ensure that the added Flags below are parsed somewhere, like in the
// PersistentPreRun func on parent command.

func init() {
	flags := []*pflag.FlagSet{
		cmdnode.NodeFlags(),
		p2p.Flags(),
		header.Flags(),
		cmdnode.MiscFlags(),
		// NOTE: for now, state-related queries can only be accessed
		// over an RPC connection with a elysium-core node.
		core.Flags(),
		rpc.Flags(),
		gateway.Flags(),
		state.Flags(),
	}

	fullCmd.AddCommand(
		cmdnode.Init(flags...),
		cmdnode.Start(flags...),
		cmdnode.AuthCmd(flags...),
		cmdnode.ResetStore(flags...),
		cmdnode.RemoveConfigCmd(flags...),
		cmdnode.UpdateConfigCmd(flags...),
	)
}

var fullCmd = &cobra.Command{
	Use:   "full [subcommand]",
	Args:  cobra.NoArgs,
	Short: "Manage your Full node",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return persistentPreRunEnv(cmd, node.Full, args)
	},
}
