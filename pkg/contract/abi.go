package contract

import (
	"encoding/json"

	"github.com/tgpxdev/gotron-sdk/pkg/proto/core/contract"
)

// JSONABI data format
type JSONABI struct {
	Anonymous bool `json:"anonymous"`
	Constant  bool `json:"constant"`
	Inputs    []struct {
		Indexed bool   `json:"indexed"`
		Name    string `json:"name"`
		Type    string `json:"type"`
	} `json:"inputs"`
	Name    string `json:"name"`
	Outputs []struct {
		Indexed bool   `json:"indexed"`
		Name    string `json:"name"`
		Type    string `json:"type"`
	} `json:"outputs"`
	Payable         bool   `json:"payable"`
	StateMutability string `json:"stateMutability"`
	Type            string `json:"type"`
}

func getState(str string) contract.SmartContract_ABI_Entry_StateMutabilityType {
	switch str {
	case "pure":
		return contract.SmartContract_ABI_Entry_Pure
	case "view":
		return contract.SmartContract_ABI_Entry_View
	case "nonpayable":
		return contract.SmartContract_ABI_Entry_Nonpayable
	case "payable":
		return contract.SmartContract_ABI_Entry_Payable
	default:
		return contract.SmartContract_ABI_Entry_UnknownMutabilityType
	}
}
func getType(str string) contract.SmartContract_ABI_Entry_EntryType {
	switch str {
	case "constructor":
		return contract.SmartContract_ABI_Entry_Constructor
	case "function":
		return contract.SmartContract_ABI_Entry_Function
	case "event":
		return contract.SmartContract_ABI_Entry_Event
	case "fallback":
		return contract.SmartContract_ABI_Entry_Fallback
	default:
		return contract.SmartContract_ABI_Entry_UnknownEntryType
	}
}

// JSONtoABI converts json string to ABI entry
func JSONtoABI(jsonSTR string) (*contract.SmartContract_ABI, error) {
	jABI := []JSONABI{}
	if err := json.Unmarshal([]byte(jsonSTR), &jABI); err != nil {
		return nil, err
	}
	ABI := &contract.SmartContract_ABI{}

	for _, v := range jABI {
		inputs := []*contract.SmartContract_ABI_Entry_Param{}
		for _, input := range v.Inputs {
			inputs = append(inputs, &contract.SmartContract_ABI_Entry_Param{
				Indexed: input.Indexed,
				Name:    input.Name,
				Type:    input.Type,
			})
		}
		outputs := []*contract.SmartContract_ABI_Entry_Param{}
		for _, output := range v.Outputs {
			outputs = append(outputs, &contract.SmartContract_ABI_Entry_Param{
				Indexed: output.Indexed,
				Name:    output.Name,
				Type:    output.Type,
			})
		}
		ABI.Entrys = append(ABI.Entrys,
			&contract.SmartContract_ABI_Entry{
				Anonymous:       v.Anonymous,
				Constant:        v.Constant,
				Name:            v.Name,
				Payable:         v.Payable,
				Inputs:          inputs,
				Outputs:         outputs,
				Type:            getType(v.Type),
				StateMutability: getState(v.StateMutability),
			})
	}
	return ABI, nil
}
