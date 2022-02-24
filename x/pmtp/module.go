package pmtp

import (
	"encoding/json"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Sifchain/sifnode/x/pmtp/types"
)

var (
	ModuleName                       = types.ModuleName
	_          module.AppModule      = AppModule{}
	_          module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module.
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) { //nolint
}

// DefaultGenesis returns default genesis state as raw bytes.
func (b AppModuleBasic) DefaultGenesis(marshaler codec.JSONCodec) json.RawMessage {
	return marshaler.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation.
func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, _ sdkclient.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx sdkclient.Context, mux *runtime.ServeMux) {
}

// RegisterRESTRoutes registers the REST routes.
func (b AppModuleBasic) RegisterRESTRoutes(ctx sdkclient.Context, router *mux.Router) {
}

// GetTxCmd returns the root tx command.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns no root query command.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions.
type AppModuleSimulation struct{}

// AppModule implements an application module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation
	Keeper types.Keeper
	Codec  *codec.Codec
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper types.Keeper, cdc *codec.Codec) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		Keeper:              keeper,
		Codec:               cdc,
	}
}

// Name returns the module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers the module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
}

// Route returns the message routing key for the module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, am.NewHandler())
}

// NewHandler returns an sdk.Handler for the module.
func (am AppModule) NewHandler() sdk.Handler {
	return nil
}

// QuerierRoute returns the module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// Deprecated: LegacyQuerierHandler use RegisterServices
func (am AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier { //nolint
	return nil
}

// InitGenesis performs genesis initialization. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, marshaler codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

// ExportGenesis returns the exported genesis state as raw bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, marshaler codec.JSONCodec) json.RawMessage {
	return nil
}

// BeginBlock returns the begin blocker.
func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker. It returns no validator updates.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (AppModule) ConsensusVersion() uint64 { return 1 }