package keeper

import (
	"github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) BeginBlock(ctx sdk.Context) {
	params := k.GetParams(ctx)
	for it := k.GetLiquidityProviderIterator(ctx); it.Valid(); it.Next() {
		value := it.Value()
		var lp types.LiquidityProvider
		if len(value) <= 0 {
			continue
		}
		err := k.cdc.Unmarshal(value, &lp)
		if err != nil {
			continue
		}
		k.PruneUnlockRecords(ctx, lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)
	}
}

func EndBlock(ctx sdk.Context, _ abci.RequestEndBlock, keeper Keeper) []abci.ValidatorUpdate {
	params := keeper.GetParams(ctx)
	pools := keeper.GetPools(ctx)
	currentPeriod := keeper.GetCurrentRewardPeriod(ctx, params)
	if currentPeriod != nil {
		err := keeper.DistributeDepthRewards(ctx, currentPeriod, pools)
		if err != nil {
			panic(err)
		}
	}

	keeper.PruneRewardPeriods(ctx, params)

	return []abci.ValidatorUpdate{}
}

func (keeper Keeper) GetCurrentRewardPeriod(ctx sdk.Context, params types.Params) *types.RewardPeriod {
	height := uint64(ctx.BlockHeight())
	for _, period := range params.RewardPeriods {
		if height >= period.StartBlock && height <= period.EndBlock {
			return period
		}
	}
	return nil
}

func (k Keeper) PruneRewardPeriods(ctx sdk.Context, params types.Params) {
	height := uint64(ctx.BlockHeight())
	var write bool
	var periods []*types.RewardPeriod
	for _, period := range params.RewardPeriods {
		if period.EndBlock > height {
			write = true
			continue
		}

		periods = append(periods, period)
	}

	if write {
		params.RewardPeriods = periods
		k.SetParams(ctx, params)
	}
}

func (k Keeper) DistributeDepthRewards(ctx sdk.Context, period *types.RewardPeriod, pools []*types.Pool) error {
	distributed := k.GetRewardsDistributed(ctx)
	remaining := period.Allocation.Sub(distributed)
	periodLength := period.EndBlock - period.StartBlock
	blockDistribution := remaining.QuoUint64(periodLength)

	if remaining.IsZero() || blockDistribution.IsZero() {
		return nil
	}

	totalDepth := sdk.ZeroUint()
	for _, pool := range pools {
		totalDepth = totalDepth.Add(pool.NativeAssetBalance)
	}

	for _, pool := range pools {
		weight := pool.NativeAssetBalance.Quo(totalDepth)
		poolDistribution := blockDistribution.Mul(weight)
		if poolDistribution.GT(remaining) {
			poolDistribution = remaining
		}
		rewardCoins := sdk.NewCoins(sdk.NewCoin(types.GetSettlementAsset().Symbol, sdk.NewIntFromUint64(poolDistribution.Uint64())))
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, rewardCoins)
		if err != nil {
			return err
		}
		pool.NativeAssetBalance = pool.NativeAssetBalance.Add(poolDistribution)
		remaining = remaining.Sub(poolDistribution)
		distributed = distributed.Add(poolDistribution)
		err = k.SetPool(ctx, pool)
		if err != nil {
			return err
		}
	}

	k.SetRewardsDistributed(ctx, distributed)

	return nil
}

func (k Keeper) GetRewardsDistributed(ctx sdk.Context) sdk.Uint {
	var rewardExecution types.RewardExecution
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardExecutionPrefix)
	if bz == nil {
		return sdk.ZeroUint()
	}
	k.cdc.MustUnmarshal(bz, &rewardExecution)
	return rewardExecution.Distributed
}

func (k Keeper) SetRewardsDistributed(ctx sdk.Context, distributed sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	rewardsExecution := types.RewardExecution{
		Distributed: distributed,
	}
	bz := k.cdc.MustMarshal(&rewardsExecution)
	store.Set(types.RewardExecutionPrefix, bz)
}

func (k Keeper) UseUnlockedLiquidity(ctx sdk.Context, lp types.LiquidityProvider, units sdk.Uint) error {
	// Ensure there is enough liquidity requested for unlock, and also passed lock period.
	// Reduce liquidity in one or more unlock records.
	// Remove unlock records with zero units remaining.
	params := k.GetParams(ctx)
	currentHeight := ctx.BlockHeight()
	lockPeriod := params.LiquidityRemovalLockPeriod

	unitsLeftToUse := units
	for _, record := range lp.Unlocks {
		if record.RequestHeight+int64(lockPeriod) <= currentHeight {
			if unitsLeftToUse.GT(record.Units) {
				// use all this record's unit's and continue with remaining
				unitsLeftToUse = unitsLeftToUse.Sub(record.Units)
				record.Units = sdk.ZeroUint()
			} else {
				// use a portion of this record's units and break
				record.Units = record.Units.Sub(unitsLeftToUse)
				unitsLeftToUse = sdk.ZeroUint()
				break
			}
		}
	}

	if !unitsLeftToUse.IsZero() {
		return types.ErrBalanceNotAvailable
	}

	// prune records.
	var records []*types.LiquidityUnlock
	for _, record := range lp.Unlocks {
		/* move to begin blocker
		if currentHeight >= record.RequestHeight + int64(lockPeriod) + cancelPeriod {
			// prune auto cancelled record
			continue
		}*/
		if record.Units.IsZero() {
			// prune used / zero record
			continue
		}
		records = append(records, record)
	}

	lp.Unlocks = records
	k.SetLiquidityProvider(ctx, &lp)

	return nil
}

func (k Keeper) PruneUnlockRecords(ctx sdk.Context, lp types.LiquidityProvider, lockPeriod, cancelPeriod uint64) {
	currentHeight := ctx.BlockHeight()

	var write bool
	var records []*types.LiquidityUnlock
	for _, record := range lp.Unlocks {
		if currentHeight >= record.RequestHeight+int64(lockPeriod)+int64(cancelPeriod) {
			// prune auto cancelled record
			write = true
			continue
		}
		if record.Units.IsZero() {
			// prune used / zero record
			write = true
			continue
		}
		records = append(records, record)
	}

	if write {
		lp.Unlocks = records
		k.SetLiquidityProvider(ctx, &lp)
	}
}