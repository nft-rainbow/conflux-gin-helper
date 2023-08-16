package utils

import (
	"fmt"
	"net/http"
)

type ChainType uint
type ChainID uint
type ContractType uint

const ONE_CFX_IN_DRIP uint64 = 1000000000000000000
const ONE_GDRIP_IN_DRIP uint64 = 1000000000

// chain types
const (
	CHAIN_TYPE_CFX ChainType = iota + 1
	CHAIN_TYPE_ETH
)

const (
	CHAINID_CFX_MAINNET ChainID = 1029
	CHAINID_CFX_TESTNET ChainID = 1
)

// contract types
const (
	CONTRACT_TYPE_ERC721 ContractType = iota + 1
	CONTRACT_TYPE_ERC1155
)

// contract type names
const (
	ERC721  = "erc721"
	ERC1155 = "erc1155"
)

const (
	CONFLUX_TEST = "conflux_test"
	CONFLUX      = "conflux"
)

var (
	MintPaths = map[string]bool{
		"/v1/mints/":                   true,
		"/v1/mints/customizable":       true,
		"/v1/mints/customizable/batch": true,
		"/v1/mints/easy/files":         true,
		"/v1/mints/easy/urls":          true,
	}
	DeployPaths = map[string]bool{
		"/v1/contracts/":                true,
		"/dashboard/apps/:id/contracts": true,
	}
)

func ChainInfoByName(name string) (ChainType, ChainID, error) {
	switch name {
	case CONFLUX_TEST:
		return CHAIN_TYPE_CFX, 1, nil
	case CONFLUX:
		return CHAIN_TYPE_CFX, 1029, nil
	default:
		return 0, 0, fmt.Errorf("unknown chain name: %s", name)
	}
}

func ChainNameByType(t ChainType, id ChainID) (string, error) {
	switch t {
	case CHAIN_TYPE_CFX:
		switch id {
		case CHAINID_CFX_TESTNET:
			return CONFLUX_TEST, nil
		case CHAINID_CFX_MAINNET:
			return CONFLUX, nil
		default:
			return "", fmt.Errorf("unknown chain id: %d", id)
		}
	default:
		return "", fmt.Errorf("unknown chain type: %d", t)
	}
}

func MustGetChainNameByType(t ChainType, id ChainID) string {
	name, err := ChainNameByType(t, id)
	if err != nil {
		panic(err)
	}
	return name
}

func ContractTypeByName(name string) (ContractType, error) {
	switch name {
	case ERC721:
		return CONTRACT_TYPE_ERC721, nil
	case ERC1155:
		return CONTRACT_TYPE_ERC1155, nil
	default:
		return 0, fmt.Errorf("unknown contract type: %s", name)
	}
}

func MustGetContractTypeByName(name string) ContractType {
	t, err := ContractTypeByName(name)
	if err != nil {
		panic(err)
	}
	return t
}

func ContractNameByType(t ContractType) (string, error) {
	switch t {
	case CONTRACT_TYPE_ERC721:
		return ERC721, nil
	case CONTRACT_TYPE_ERC1155:
		return ERC1155, nil
	default:
		return "", fmt.Errorf("unknown contract type: %d", t)
	}
}

func MustGetContractNameByType(t ContractType) string {
	n, err := ContractNameByType(t)
	if err != nil {
		panic(err)
	}
	return n
}

// contains batch mint
func IsMint(method string, path string) bool {
	return method == http.MethodPost && MintPaths[path]
}

func IsBatchMint(method string, path string) bool {
	return method == http.MethodPost && path == "/v1/mints/customizable/batch"
}

func IsDeploy(method string, path string) bool {
	return method == http.MethodPost && DeployPaths[path]
}

func IsTestnetByName(chain string) bool {
	return chain == CONFLUX_TEST
}

func IsMainnetByName(chain string) bool {
	return chain == CONFLUX
}

func IsTestnet(chainType, chaintID uint) bool {
	chainName, _ := ChainNameByType(ChainType(chainType), ChainID(chaintID))
	isTestnet := IsTestnetByName(chainName)
	return isTestnet
}
