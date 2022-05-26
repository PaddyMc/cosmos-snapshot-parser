package parser

import sdk "github.com/cosmos/cosmos-sdk/types"

// SetConfig initialize the configuration instance for the sdk
func SetConfig(AccountAddressPrefix string) {
	// Various prefixes for accounts and public keys
	var (
		AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
		ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
		ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
		ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
		ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
	)
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}
