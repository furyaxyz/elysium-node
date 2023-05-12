package main

import (
	"github.com/spf13/cobra"

	cmdnode "github.com/elysiumorg/elysium-node/cmd"
	"github.com/elysiumorg/elysium-node/nodebuilder/core"
	"github.com/elysiumorg/elysium-node/nodebuilder/gateway"
	"github.com/elysiumorg/elysium-node/nodebuilder/header"
	"github.com/elysiumorg/elysium-node/nodebuilder/node"
	"github.com/elysiumorg/elysium-node/nodebuilder/p2p"
	"github.com/elysiumorg/elysium-node/nodebuilder/rpc"
	"github.com/elysiumorg/elysium-node/nodebuilder/state"
)

func persistentPreRunEnv(cmd *cobra.Command, nodeType node.Type, _ []string) error {
	var (
		ctx = cmd.Context()
		err error
	)

	ctx = cmdnode.WithNodeType(ctx, nodeType)

	parsedNetwork, err := p2p.ParseNetwork(cmd)
	if err != nil {
		return err
	}
	ctx = cmdnode.WithNetwork(ctx, parsedNetwork)
	ctx = cmdnode.WithNodeBuildInfo(ctx, &node.BuildInfo{
		LastCommit:      lastCommit,
		SemanticVersion: semanticVersion,
		SystemVersion:   systemVersion,
		GolangVersion:   golangVersion,
	})

	// loads existing config into the environment
	ctx, err = cmdnode.ParseNodeFlags(ctx, cmd, cmdnode.Network(ctx))
	if err != nil {
		return err
	}

	cfg := cmdnode.NodeConfig(ctx)

	err = p2p.ParseFlags(cmd, &cfg.P2P)
	if err != nil {
		return err
	}

	err = core.ParseFlags(cmd, &cfg.Core)
	if err != nil {
		return err
	}

	if nodeType != node.Bridge {
		err = header.ParseFlags(cmd, &cfg.Header)
		if err != nil {
			return err
		}
	}

	ctx, err = cmdnode.ParseMiscFlags(ctx, cmd)
	if err != nil {
		return err
	}

	rpc.ParseFlags(cmd, &cfg.RPC)
	gateway.ParseFlags(cmd, &cfg.Gateway)
	state.ParseFlags(cmd, &cfg.State)

	// set config
	ctx = cmdnode.WithNodeConfig(ctx, &cfg)
	cmd.SetContext(ctx)
	return nil
}
