package recovery

import (
	"encoding/json"

	"github.com/spf13/cobra"
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/atlas/chain/x/recovery/keeper"
	"github.com/atlas/chain/x/recovery/types"
)

type AppModuleBasic struct{ cdc codec.Codec }
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func (AppModuleBasic) Name() string { return types.ModuleName }
func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}
func (AppModuleBasic) RegisterInterfaces(types.InterfaceRegistry) {}
func (AppModuleBasic) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	return json.RawMessage("{}")
}
func (AppModuleBasic) ValidateGenesis(codec.JSONCodec, client.TxEncodingConfig, json.RawMessage) error {
	return nil
}
func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }
func (AppModuleBasic) GetQueryCmd() *cobra.Command { return nil }

func NewAppModule(cdc codec.Codec, k keeper.Keeper) AppModule {
	return AppModule{AppModuleBasic{cdc}, k}
}
func (am AppModule) Name() string { return am.AppModuleBasic.Name() }
func (am AppModule) RegisterServices(module.Configurator) {}
func (am AppModule) InitGenesis(sdk.Context, codec.JSONCodec, json.RawMessage) {}
func (am AppModule) ExportGenesis(sdk.Context, codec.JSONCodec) json.RawMessage {
	return json.RawMessage("{}")
}
func (AppModule) ConsensusVersion() uint64 { return 1 }
func (AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {}
func (AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

