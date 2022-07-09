package topup

import (
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/topup/simulation"

	topupCli "github.com/maticnetwork/heimdall/topup/client/cli"
	topupRest "github.com/maticnetwork/heimdall/topup/client/rest"

	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmModule "github.com/maticnetwork/heimdall/types/module"
	simTypes "github.com/maticnetwork/heimdall/types/simulation"
)

var (
	_ module.AppModule             = AppModule{}
	_ module.AppModuleBasic        = AppModuleBasic{}
	_ hmModule.HeimdallModuleBasic = AppModule{}
)

// AppModuleBasic defines the basic application module used by the auth module.
type AppModuleBasic struct{}

// Name returns the auth module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the auth module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the auth
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the auth module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}

	return types.ValidateGenesis(data)
}

// VerifyGenesis performs verification on auth module state.
func (AppModuleBasic) VerifyGenesis(bz map[string]json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the auth module.
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	topupRest.RegisterRoutes(ctx, rtr)
}

// GetTxCmd returns the root tx command for the auth module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return topupCli.GetTxCmd(cdc)
}

// GetQueryCmd returns the root query command for the auth module.
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return topupCli.GetQueryCmd(cdc)
}

//____________________________________________________________________________

// AppModule implements an application module for the auth module.
type AppModule struct {
	AppModuleBasic

	keeper         Keeper
	contractCaller helper.IContractCaller
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, contractCaller helper.IContractCaller) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		contractCaller: contractCaller,
	}
}

// Name returns the auth module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants performs a no-op.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the auth module.
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper, am.contractCaller)
}

// QuerierRoute returns the auth module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the auth module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper, am.contractCaller)
}

// InitGenesis performs genesis initialization for the auth module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, genesisState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the auth
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the auth module.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the auth module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// GenerateGenesisState creates a randomized GenState of the chainManager module
func (AppModule) GenerateGenesisState(simState *hmModule.SimulationState) {
	simulation.RandomizeGenState(simState)
}

// ProposalContents doesn't return any content functions.
func (AppModule) ProposalContents(simState hmModule.SimulationState) []simTypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized param changes for the simulator.
func (AppModule) RandomizedParams(r *rand.Rand) []simTypes.ParamChange {
	return nil
}

// RegisterStoreDecoder registers a decoder for chainmanager module's types
func (AppModule) RegisterStoreDecoder(sdr hmModule.StoreDecoderRegistry) {
}

// WeightedOperations doesn't return any chainmanager module operation.
func (AppModule) WeightedOperations(_ hmModule.SimulationState) []simTypes.WeightedOperation {
	return nil
}

func (am AppModule) NewSideTxHandler() hmTypes.SideTxHandler {
	return NewSideTxHandler(am.keeper, am.contractCaller)
}

// NewPostTxHandler side tx handler
func (am AppModule) NewPostTxHandler() hmTypes.PostTxHandler {
	return NewPostTxHandler(am.keeper, am.contractCaller)
}