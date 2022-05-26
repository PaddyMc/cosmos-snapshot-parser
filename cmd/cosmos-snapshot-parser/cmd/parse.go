package cmd

import (
	"github.com/cosmos/cosmos-sdk/codec"
	ct "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	"github.com/neilotoole/errgroup"
	"github.com/plural-labs/cosmos-snapshot-parser/parser"
	"github.com/spf13/cobra"
)

func parseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "parse data from the application store and block store",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			errs, _ := errgroup.WithContext(ctx)
			var err error
			errs.Go(func() error {
				interfaceRegistry := ct.NewInterfaceRegistry()
				authtypes.RegisterInterfaces(interfaceRegistry)
				banktypes.RegisterInterfaces(interfaceRegistry)
				distrtypes.RegisterInterfaces(interfaceRegistry)
				evidencetypes.RegisterInterfaces(interfaceRegistry)
				govtypes.RegisterInterfaces(interfaceRegistry)
				slashingtypes.RegisterInterfaces(interfaceRegistry)
				stakingtypes.RegisterInterfaces(interfaceRegistry)
				upgradetypes.RegisterInterfaces(interfaceRegistry)
				ibctransfertypes.RegisterInterfaces(interfaceRegistry)
				cryptocodec.RegisterInterfaces(interfaceRegistry)
				marshaler := codec.NewProtoCodec(interfaceRegistry)
				if err = parser.Parse(
					accountPrefix,
					dataDir,
					connectionString,
					numberOfBlocks,
					marshaler,
				); err != nil {
					return err
				}
				return nil
			})

			return errs.Wait()
		},
	}
	return cmd
}
