package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),

		&MsgCreateEthBridgeClaim{},
		&MsgBurn{},
		&MsgLock{},
		&MsgUpdateWhiteListValidator{},
		&MsgUpdateCrossChainFeeReceiverAccount{},
		&MsgRescueCrossChainFee{},
		&MsgSignProphecy{},
		&MsgSetFeeInfo{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

const (
	// Used in debugging logging statements for logs that are interesting for the Peggy test environment
	PeggyTestMarker = "peggytest"
)
