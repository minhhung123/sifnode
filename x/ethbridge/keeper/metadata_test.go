package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/ethbridge/test"
	"github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/require"
)

const invalidDenom = "A Nonexistant Denom Hash"

var testTokenMetadata = types.TokenMetadata{
	Decimals:          15,
	Name:              "Test Token",
	NetworkDescriptor: oracletypes.NetworkDescriptor_NETWORK_DESCRIPTOR_BINANCE_SMART_CHAIN,
	Symbol:            "TT",
	TokenAddress:      "0x0123456789ABCDEF",
}

func TestGetAddTokenMetadata(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	expected := types.TokenMetadata{}
	result, _ := keeper.GetTokenMetadata(ctx, invalidDenom)
	require.Equal(t, expected, result)
	resultDenom := keeper.AddTokenMetadata(ctx, testTokenMetadata)
	expectedDenom := types.GetDenomHash(
		testTokenMetadata.NetworkDescriptor,
		testTokenMetadata.TokenAddress,
		testTokenMetadata.Decimals,
		testTokenMetadata.Name,
		testTokenMetadata.Symbol,
	)
	require.Equal(t, expectedDenom, resultDenom)
	result, _ = keeper.GetTokenMetadata(ctx, resultDenom)
	require.Equal(t, testTokenMetadata, result)
}

func TestExistsTokenMetadata(t *testing.T) {
	ctx, keeper, _, _, _, _, _, _ := test.CreateTestKeepers(t, 0.7, []int64{3, 3}, "")
	expected := false
	result := keeper.ExistsTokenMetadata(ctx, invalidDenom)
	require.Equal(t, expected, result)
	denom := keeper.AddTokenMetadata(ctx, testTokenMetadata)
	expected = true
	result = keeper.ExistsTokenMetadata(ctx, denom)
	require.Equal(t, expected, result)
}