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

	//	gammtypes "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
	//	incentivestypes "github.com/osmosis-labs/osmosis/v7/x/incentives/types"
	//	lockuptypes "github.com/osmosis-labs/osmosis/v7/x/lockup/types"
	//	pooltypes "github.com/osmosis-labs/osmosis/v7/x/pool-incentives/types"
	//	superfluidtypes "github.com/osmosis-labs/osmosis/v7/x/superfluid/types"
	//	txfeestypes "github.com/osmosis-labs/osmosis/v7/x/txfees/types"
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

				// The marshaler is defined here, as each chain has
				// their own custom proto types needed for
				// when we unmarshal transactions from
				// a block.
				// XXX: should we have one 'parse' command per chain?

				// Default cosmos codec
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

				// Default osmo codec
				//gammtypes.RegisterInterfaces(interfaceRegistry)
				//incentivestypes.RegisterInterfaces(interfaceRegistry)
				//lockuptypes.RegisterInterfaces(interfaceRegistry)
				//superfluidtypes.RegisterInterfaces(interfaceRegistry)
				//pooltypes.RegisterInterfaces(interfaceRegistry)
				//txfeestypes.RegisterInterfaces(interfaceRegistry)

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
