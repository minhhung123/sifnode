package keeper_test

import (
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/Sifchain/sifnode/x/clp/test"
	"github.com/Sifchain/sifnode/x/clp/types"
)

func TestKeeper_PolicyRun(t *testing.T) {
	testcases := []struct {
		name                      string
		createBalance             bool
		createPool                bool
		createLPs                 bool
		poolAsset                 string
		address                   string
		nativeBalance             sdk.Int
		externalBalance           sdk.Int
		nativeAssetAmount         sdk.Uint
		externalAssetAmount       sdk.Uint
		poolUnits                 sdk.Uint
		poolAssetDecimals         int64
		poolAssetPermissions      []tokenregistrytypes.Permission
		nativeAssetPermissions    []tokenregistrytypes.Permission
		expectedPool              types.Pool
		expectedSwapPriceNative   sdk.Dec
		expectedSwapPriceExternal sdk.Dec
		err                       error
		errString                 error
	}{
		{
			name:                   "18 decimals asset",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "eth",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUint(1000),
			externalAssetAmount:    sdk.NewUint(1000),
			poolUnits:              sdk.NewUint(1000),
			poolAssetDecimals:      18,
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "eth"},
				NativeAssetBalance:            sdk.NewUint(1000),
				ExternalAssetBalance:          sdk.NewUint(1000),
				PoolUnits:                     sdk.NewUint(1000),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.998002996005000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("0.998002996005000000"),
		},
		{
			name:                   "cusdt with 6 decimals",
			createBalance:          true,
			createPool:             true,
			createLPs:              true,
			poolAsset:              "cusdt",
			address:                "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
			nativeBalance:          sdk.NewInt(10000),
			externalBalance:        sdk.NewInt(10000),
			nativeAssetAmount:      sdk.NewUintFromString("1550459183129248235861408"),
			externalAssetAmount:    sdk.NewUintFromString("174248776094"),
			poolUnits:              sdk.NewUintFromString("1550459183129248235861408"),
			poolAssetDecimals:      6,
			poolAssetPermissions:   []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			nativeAssetPermissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP},
			expectedPool: types.Pool{
				ExternalAsset:                 &types.Asset{Symbol: "cusdt"},
				NativeAssetBalance:            sdk.NewUintFromString("1550459183129248235861408"),
				ExternalAssetBalance:          sdk.NewUintFromString("174248776094"),
				PoolUnits:                     sdk.NewUintFromString("1550459183129248235861408"),
				RewardPeriodNativeDistributed: sdk.ZeroUint(),
			},
			expectedSwapPriceNative:   sdk.MustNewDecFromStr("0.112385271402000000"),
			expectedSwapPriceExternal: sdk.MustNewDecFromStr("8.897963118404021251"),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				trGs := &tokenregistrytypes.GenesisState{
					AdminAccounts: &tokenregistrytypes.AdminAccounts{AdminAccounts: []*tokenregistrytypes.AdminAccount{
						{
							AdminType:    tokenregistrytypes.AdminType_PMTPREWARDS,
							AdminAddress: tc.address,
						},
					}},
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: tc.poolAsset, BaseDenom: tc.poolAsset, Decimals: tc.poolAssetDecimals, Permissions: tc.poolAssetPermissions},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: tc.nativeAssetPermissions},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				if tc.createBalance {
					balances := []banktypes.Balance{
						{
							Address: tc.address,
							Coins: sdk.Coins{
								sdk.NewCoin(tc.poolAsset, tc.externalBalance),
								sdk.NewCoin("rowan", tc.nativeBalance),
							},
						},
					}

					bankGs := banktypes.DefaultGenesisState()
					bankGs.Balances = append(bankGs.Balances, balances...)
					bz, _ = app.AppCodec().MarshalJSON(bankGs)
					genesisState["bank"] = bz
				}

				if tc.createPool {
					pools := []*types.Pool{
						{
							ExternalAsset:        &types.Asset{Symbol: tc.poolAsset},
							NativeAssetBalance:   tc.nativeAssetAmount,
							ExternalAssetBalance: tc.externalAssetAmount,
							PoolUnits:            tc.poolUnits,
						},
					}
					clpGs := types.DefaultGenesisState()
					if tc.createLPs {
						lps := []*types.LiquidityProvider{
							{
								Asset:                    &types.Asset{Symbol: tc.poolAsset},
								LiquidityProviderAddress: tc.address,
								LiquidityProviderUnits:   tc.nativeAssetAmount,
							},
						}
						clpGs.LiquidityProviders = append(clpGs.LiquidityProviders, lps...)
					}
					clpGs.Params = types.Params{
						MinCreatePoolThreshold: 100,
					}
					clpGs.AddressWhitelist = append(clpGs.AddressWhitelist, tc.address)
					clpGs.PoolList = append(clpGs.PoolList, pools...)
					bz, _ = app.AppCodec().MarshalJSON(clpGs)
					genesisState["clp"] = bz
				}

				return genesisState
			})

			app.ClpKeeper.SetPmtpCurrentRunningRate(ctx, sdk.NewDec(0))

			err := app.ClpKeeper.PolicyRun(ctx, sdk.NewDec(0))

			if tc.errString != nil {
				require.EqualError(t, err, tc.errString.Error())
				return
			}
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)

			pool, _ := app.ClpKeeper.GetPool(ctx, tc.poolAsset)

			tc.expectedPool.SwapPriceNative = &tc.expectedSwapPriceNative
			tc.expectedPool.SwapPriceExternal = &tc.expectedSwapPriceExternal

			require.Equal(t, pool, tc.expectedPool)
		})
	}
}
