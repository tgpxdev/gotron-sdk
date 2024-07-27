package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/spf13/cobra"
	"github.com/tgpxdev/gotron-sdk/pkg/address"
	"github.com/tgpxdev/gotron-sdk/pkg/common"
	"github.com/tgpxdev/gotron-sdk/pkg/proto/core"
	"github.com/tgpxdev/gotron-sdk/pkg/proto/core/contract"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ()

func bcSub() []*cobra.Command {
	cmdNode := &cobra.Command{
		Use:   "node",
		Short: "get node metrics",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := conn.GetNodeInfo()
			if err != nil {
				return err
			}

			if noPrettyOutput {
				fmt.Println(info)
				return nil
			}

			asJSON, _ := json.Marshal(info)
			fmt.Println(common.JSONPrettyFormat(string(asJSON)))
			return nil
		},
	}

	cmdMT := &cobra.Command{
		Use:   "mt",
		Short: "get network next maintainance time",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := conn.GetNextMaintenanceTime()
			if err != nil {
				return err
			}

			if noPrettyOutput {
				fmt.Println(info)
				return nil
			}

			t := time.Unix(info.GetNum()/1000, 0)
			result := make(map[string]interface{})
			result["nextTimestamp"] = info.GetNum()
			result["date"] = t.UTC().Format(time.RFC3339)

			asJSON, _ := json.Marshal(result)
			fmt.Println(common.JSONPrettyFormat(string(asJSON)))
			return nil
		},
	}

	cmdTX := &cobra.Command{
		Use:   "tx <HASH>",
		Short: "get tx info by hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tx, err := conn.GetTransactionByID(args[0])
			if err != nil {
				return err
			}
			contracts := tx.GetRawData().GetContract()
			if len(contracts) != 1 {
				return fmt.Errorf("invalid contracts")
			}
			ctt := contracts[0]

			info, err := conn.GetTransactionInfoByID(args[0])
			if err != nil {
				return err
			}

			if noPrettyOutput {
				fmt.Println(tx, info)
				return nil
			}

			result := make(map[string]interface{})
			t := time.Unix(info.GetBlockTimeStamp()/1000, 0)
			result["txID"] = common.BytesToHexString(info.Id)
			result["block"] = info.GetBlockNumber()
			result["timestamp"] = info.GetBlockTimeStamp()
			result["date"] = t.UTC().Format(time.RFC3339)

			result["receipt"] = map[string]interface{}{
				"fee":               info.GetFee(),
				"energyFee":         info.GetReceipt().GetEnergyFee(),
				"energyUsage":       info.GetReceipt().GetEnergyUsage(),
				"originEnergyUsage": info.GetReceipt().GetOriginEnergyUsage(),
				"energyUsageTotal":  info.GetReceipt().GetEnergyUsageTotal(),
				"netFee":            info.GetReceipt().GetNetFee(),
				"netUsage":          info.GetReceipt().GetNetUsage(),
			}

			result["contractName"] = ctt.Type.String()
			//parse contract
			var c interface{}
			switch ctt.Type {
			case core.Transaction_Contract_AccountCreateContract:
				c = &contract.AccountCreateContract{}
			case core.Transaction_Contract_TransferContract:
				c = &contract.TransferContract{}
			case core.Transaction_Contract_TransferAssetContract:
				c = &contract.TransferAssetContract{}
			case core.Transaction_Contract_VoteWitnessContract:
				c = &contract.VoteWitnessContract{}
			case core.Transaction_Contract_WitnessCreateContract:
				c = &contract.WitnessCreateContract{}
			case core.Transaction_Contract_WitnessUpdateContract:
				c = &contract.WitnessUpdateContract{}
			case core.Transaction_Contract_AssetIssueContract:
				c = &contract.AssetIssueContract{}
			case core.Transaction_Contract_ParticipateAssetIssueContract:
				c = &contract.ParticipateAssetIssueContract{}
			case core.Transaction_Contract_AccountUpdateContract:
				c = &contract.AccountUpdateContract{}
			case core.Transaction_Contract_FreezeBalanceContract:
				c = &contract.FreezeBalanceContract{}
			case core.Transaction_Contract_UnfreezeBalanceContract:
				c = &contract.UnfreezeBalanceContract{}
			case core.Transaction_Contract_WithdrawBalanceContract:
				c = &contract.WithdrawBalanceContract{}
			case core.Transaction_Contract_UnfreezeAssetContract:
				c = &contract.UnfreezeAssetContract{}
			case core.Transaction_Contract_UpdateAssetContract:
				c = &contract.UpdateAssetContract{}
			case core.Transaction_Contract_ProposalCreateContract:
				c = &contract.ProposalCreateContract{}
			case core.Transaction_Contract_ProposalApproveContract:
				c = &contract.ProposalApproveContract{}
			case core.Transaction_Contract_ProposalDeleteContract:
				c = &contract.ProposalDeleteContract{}
			case core.Transaction_Contract_SetAccountIdContract:
				c = &contract.SetAccountIdContract{}
			case core.Transaction_Contract_CustomContract:
				return fmt.Errorf("proto unmarshal any: %s", "customContract")
			case core.Transaction_Contract_CreateSmartContract:
				c = &contract.CreateSmartContract{}
			case core.Transaction_Contract_TriggerSmartContract:
				c = &contract.TriggerSmartContract{}
			case core.Transaction_Contract_GetContract:
				return fmt.Errorf("proto unmarshal any: %s", "getContract")
			case core.Transaction_Contract_UpdateSettingContract:
				c = &contract.UpdateSettingContract{}
			case core.Transaction_Contract_ExchangeCreateContract:
				c = &contract.ExchangeCreateContract{}
			case core.Transaction_Contract_ExchangeInjectContract:
				c = &contract.ExchangeInjectContract{}
			case core.Transaction_Contract_ExchangeWithdrawContract:
				c = &contract.ExchangeWithdrawContract{}
			case core.Transaction_Contract_ExchangeTransactionContract:
				c = &contract.ExchangeTransactionContract{}
			case core.Transaction_Contract_UpdateEnergyLimitContract:
				c = &contract.UpdateEnergyLimitContract{}
			case core.Transaction_Contract_AccountPermissionUpdateContract:
				c = &contract.AccountPermissionUpdateContract{}
			case core.Transaction_Contract_ClearABIContract:
				c = &contract.ClearABIContract{}
			case core.Transaction_Contract_UpdateBrokerageContract:
				c = &contract.UpdateBrokerageContract{}
			case core.Transaction_Contract_ShieldedTransferContract:
				c = &contract.ShieldedTransferContract{}
			case core.Transaction_Contract_MarketSellAssetContract:
				c = &contract.MarketSellAssetContract{}
			case core.Transaction_Contract_MarketCancelOrderContract:
				c = &contract.MarketCancelOrderContract{}
			case core.Transaction_Contract_FreezeBalanceV2Contract:
				c = &contract.FreezeBalanceV2Contract{}
			case core.Transaction_Contract_UnfreezeBalanceV2Contract:
				c = &contract.UnfreezeBalanceV2Contract{}
			case core.Transaction_Contract_WithdrawExpireUnfreezeContract:
				c = &contract.WithdrawExpireUnfreezeContract{}
			case core.Transaction_Contract_DelegateResourceContract:
				c = &contract.DelegateResourceContract{}
			case core.Transaction_Contract_UnDelegateResourceContract:
				c = &contract.UnDelegateResourceContract{}
			default:
				return fmt.Errorf("proto unmarshal any: %+w", err)
			}

			if err = ctt.GetParameter().UnmarshalTo(c.(protoreflect.ProtoMessage)); err != nil {
				return fmt.Errorf("proto unmarshal any: %+w", err)
			}
			result["contract"] = parseContractHumanReadable(structs.Map(c))

			asJSON, _ := json.Marshal(result)
			fmt.Println(common.JSONPrettyFormat(string(asJSON)))
			return nil
		},
	}

	return []*cobra.Command{cmdNode, cmdMT, cmdTX}
}

func parseContractHumanReadable(ck map[string]interface{}) map[string]interface{} {
	// Addresses fields
	addresses := map[string]bool{
		"OwnerAddress":    true,
		"ReceiverAddress": true,
		"ToAddress":       true,
		"ContractAddress": true,
	}
	for f, d := range ck {
		if strings.HasPrefix(f, "XXX_") {
			delete(ck, f)
		}

		// convert addresses
		if addresses[f] {
			ck[f] = address.Address(d.([]uint8)).String()
		}
	}

	if v, ok := ck["Votes"]; ok {
		votes := make(map[string]int64)
		for _, d := range v.([]interface{}) {
			dP := d.(map[string]interface{})
			votes[address.Address(dP["VoteAddress"].([]uint8)).String()] = dP["VoteCount"].(int64)
		}
		ck["Votes"] = votes
	}

	return ck
}

func init() {
	cmdBC := &cobra.Command{
		Use:   "bc",
		Short: "Blockchain Actions",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdBC.AddCommand(bcSub()...)
	RootCmd.AddCommand(cmdBC)
}
