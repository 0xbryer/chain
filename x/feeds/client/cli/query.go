package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bandprotocol/chain/v2/x/feeds/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetQueryCmdPrices(),
		GetQueryCmdPrice(),
		GetQueryCmdValidatorPrices(),
		GetQueryCmdValidValidator(),
		GetQueryCmdSignalTotalPowers(),
		GetQueryCmdParams(),
		GetQueryCmdReferenceSourceConfig(),
		GetQueryCmdDelegatorSignals(),
		GetQueryCmdCurrentFeeds(),
		GetQueryCmdAllCurrentPrices(),
		GetQueryCmdCurrentPrices(),
		GetQueryCmdIsFeeder(),
	)

	return queryCmd
}

// GetQueryCmdDelegatorSignals implements the query delegator signal command.
func GetQueryCmdDelegatorSignals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-signals [delegator-addr]",
		Short: "Shows delegator's currently active signals",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.DelegatorSignals(
				cmd.Context(),
				&types.QueryDelegatorSignalsRequest{DelegatorAddress: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdPrices implements the query prices command.
func GetQueryCmdPrices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prices",
		Short: "Shows the latest price of all signal ids",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Prices(cmd.Context(), &types.QueryPricesRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "prices")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdPrice implements the query price command.
func GetQueryCmdPrice() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "price [signal_id]",
		Short: "Shows the latest price of a signal id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Price(cmd.Context(), &types.QueryPriceRequest{
				SignalId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdCurrentFeeds implements the query current feeds command.
func GetQueryCmdCurrentFeeds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-feeds",
		Short: "Shows all currently supported feeds",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CurrentFeeds(cmd.Context(), &types.QueryCurrentFeedsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdAllCurrentPrices implements the query all current prices command.
func GetQueryCmdAllCurrentPrices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-current-prices",
		Short: "Shows all current prices of all supported feeds",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AllCurrentPrices(cmd.Context(), &types.QueryAllCurrentPricesRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdCurrentPrices implements the query current prices command.
func GetQueryCmdCurrentPrices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-prices [signal_ids]",
		Short: "Shows current prices of supported feeds",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			signalIdsStr := args[0]
			signalIds := strings.Split(signalIdsStr, ",")

			res, err := queryClient.CurrentPrices(cmd.Context(), &types.QueryCurrentPricesRequest{
				SignalIds: signalIds,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdValidatorPrices implements the query validator prices command.
func GetQueryCmdValidatorPrices() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-prices [validator]",
		Short: "Shows prices of the validator. Optionally filter by signal IDs.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			signalIds, err := cmd.Flags().GetStringSlice("signal-ids")
			if err != nil {
				return err
			}

			res, err := queryClient.ValidatorPrices(cmd.Context(), &types.QueryValidatorPricesRequest{
				ValidatorAddress: args[0],
				SignalIds:        signalIds,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().StringSlice("signal-ids", nil, "Comma-separated list of signal IDs to filter by")

	return cmd
}

// GetQueryCmdValidValidator implements the query valid validator command.
func GetQueryCmdValidValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "valid-validator [validator-address]",
		Short: "Shows if the given address is a valid validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ValidValidator(
				cmd.Context(),
				&types.QueryValidValidatorRequest{ValidatorAddress: args[0]},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdSignalTotalPowers implements the query signal-total-powers command.
func GetQueryCmdSignalTotalPowers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signal-total-powers",
		Short: "Shows all information of all signals and its total power",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.SignalTotalPowers(
				cmd.Context(),
				&types.QuerySignalTotalPowersRequest{Pagination: pageReq},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "signal-total-powers")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdReferenceSourceConfig implements the query reference source config command.
func GetQueryCmdReferenceSourceConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reference-source-config",
		Short: "Shows information of reference source config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ReferenceSourceConfig(
				cmd.Context(),
				&types.QueryReferenceSourceConfigRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdParams implements the query params command.
func GetQueryCmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetQueryCmdIsFeeder implements the query if an address is a feeder command.
func GetQueryCmdIsFeeder() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-feeder [validator-address] [feeder-address]",
		Short: "Checks if the given address is a feeder for the validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.IsFeeder(cmd.Context(), &types.QueryIsFeederRequest{
				ValidatorAddress: args[0],
				FeederAddress:    args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
