package chain

import (
	"encoding/json"
	"io"
	"os"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	dbm "github.com/cometbft/cometbft/libs/db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/x/ibc-transfer"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/cosmos/cosmos-sdk/x/vesting"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/circuit"
	
	"github.com/atlas/chain/x/compute"
	"github.com/atlas/chain/x/storage"
	"github.com/atlas/chain/x/training"
	"github.com/atlas/chain/x/reward"
	"github.com/atlas/chain/x/model"
	"github.com/atlas/chain/x/health"
	"github.com/atlas/chain/x/recovery"
	"github.com/atlas/chain/x/sharding"
	"github.com/atlas/chain/x/validation"
)

const appName = "atlas"

var (
	DefaultNodeHome string

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		crisis.AppModuleBasic{},
		gov.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibc-transfer.AppModuleBasic{},
		capability.AppModuleBasic{},
		authz.AppModuleBasic{},
		feegrant.AppModuleBasic{},
		nft.AppModuleBasic{},
		group.AppModuleBasic{},
		vesting.AppModuleBasic{},
		consensus.AppModuleBasic{},
		circuit.AppModuleBasic{},
		
		compute.AppModuleBasic{},
		storage.AppModuleBasic{},
		training.AppModuleBasic{},
		reward.AppModuleBasic{},
		model.AppModuleBasic{},
		health.AppModuleBasic{},
		recovery.AppModuleBasic{},
		sharding.AppModuleBasic{},
		validation.AppModuleBasic{},
	)
)

type AtlasApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	GovKeeper        govkeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	EvidenceKeeper   evidencekeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper
	IBCTransferKeeper ibctransferkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper
	FeegrantKeeper   feegrantkeeper.Keeper
	NFTKeeper        nftkeeper.Keeper
	GroupKeeper      groupkeeper.Keeper
	VestingKeeper    vestingkeeper.Keeper
	ConsensusKeeper  consensuskeeper.Keeper
	CircuitKeeper    circuitkeeper.Keeper
	
	ComputeKeeper    computekeeper.Keeper
	StorageKeeper    storagekeeper.Keeper
	TrainingKeeper   trainingkeeper.Keeper
	RewardKeeper     rewardkeeper.Keeper
	ModelKeeper      modelkeeper.Keeper
	HealthKeeper     healthkeeper.Keeper
	RecoveryKeeper   recoverykeeper.Keeper
	ShardingKeeper   shardingkeeper.Keeper
	ValidationKeeper validationkeeper.Keeper

	mm *module.Manager

	sm *module.SimulationManager
}

func NewAtlasApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig appparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *AtlasApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, upgradetypes.StoreKey, evidencetypes.StoreKey,
		ibcexported.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		authzkeeper.StoreKey, feegrant.StoreKey, nftkeeper.StoreKey,
		group.StoreKey, vestingtypes.StoreKey, consensustypes.StoreKey,
		circuittypes.StoreKey,
		
		computetypes.StoreKey,
		storagetypes.StoreKey,
		trainingtypes.StoreKey,
		rewardtypes.StoreKey,
		modeltypes.StoreKey,
		healthtypes.StoreKey,
		recoverytypes.StoreKey,
		shardingtypes.StoreKey,
		validationtypes.StoreKey,
	)

	app := &AtlasApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
	}

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], authtypes.ProtoBaseAccount, maccPerms,
		sdk.Bech32MainPrefix, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper,
		app.BlockedModuleAccountAddrs(), authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.StakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		app.StakingKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], app.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec, keys[crisistypes.StoreKey], invCheckPeriod,
		app.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.AccountKeeper, app.BankKeeper,
		app.StakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath,
		app.BaseApp, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.EvidenceKeeper = evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], app.StakingKeeper, app.SlashingKeeper,
	)

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec, keys[capabilitytypes.StoreKey], app.GetMemKey(capabilitytypes.MemStoreKey),
	)

	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibcexported.StoreKey], app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper, app.UpgradeKeeper, app.CapabilityKeeper.ScopeToModule("ibc"),
	)

	app.IBCTransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, app.IBCKeeper.ScopedTransferKeeper,
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper,
	)

	app.FeegrantKeeper = feegrantkeeper.NewKeeper(
		appCodec, keys[feegrant.StoreKey], app.AccountKeeper,
	)

	app.NFTKeeper = nftkeeper.NewKeeper(
		keys[nftkeeper.StoreKey], appCodec, app.AccountKeeper, app.BankKeeper,
	)

	app.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper,
		group.DefaultConfig(),
	)

	app.VestingKeeper = vestingkeeper.NewKeeper(
		app.AccountKeeper, app.BankKeeper,
	)

	app.ConsensusKeeper = consensuskeeper.NewKeeper(
		appCodec, keys[consensustypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.CircuitKeeper = circuitkeeper.NewKeeper(
		appCodec, keys[circuittypes.StoreKey], authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.ComputeKeeper = computekeeper.NewKeeper(
		appCodec, keys[computetypes.StoreKey], keys[computetypes.MemStoreKey],
		app.BankKeeper,
	)

	app.StorageKeeper = storagekeeper.NewKeeper(
		appCodec, keys[storagetypes.StoreKey], keys[storagetypes.MemStoreKey],
		app.BankKeeper,
	)

	app.TrainingKeeper = trainingkeeper.NewKeeper(
		appCodec, keys[trainingtypes.StoreKey], keys[trainingtypes.MemStoreKey],
		app.ComputeKeeper, app.StorageKeeper, app.BankKeeper,
	)

	app.RewardKeeper = rewardkeeper.NewKeeper(
		appCodec, keys[rewardtypes.StoreKey], keys[rewardtypes.MemStoreKey],
		app.BankKeeper, app.ComputeKeeper, app.StorageKeeper,
	)

	app.ModelKeeper = modelkeeper.NewKeeper(
		appCodec, keys[modeltypes.StoreKey], keys[modeltypes.MemStoreKey],
		app.StorageKeeper,
	)

	app.HealthKeeper = healthkeeper.NewKeeper(
		appCodec, keys[healthtypes.StoreKey], keys[healthtypes.MemStoreKey],
		app.ComputeKeeper,
	)

	app.RecoveryKeeper = recoverykeeper.NewKeeper(
		appCodec, keys[recoverytypes.StoreKey], keys[recoverytypes.MemStoreKey],
		app.TrainingKeeper, app.ComputeKeeper, app.HealthKeeper,
	)

	app.ShardingKeeper = shardingkeeper.NewKeeper(
		appCodec, keys[shardingtypes.StoreKey], keys[shardingtypes.MemStoreKey],
		app.StorageKeeper, app.TrainingKeeper,
	)

	app.ValidationKeeper = validationkeeper.NewKeeper(
		appCodec, keys[validationtypes.StoreKey], keys[validationtypes.MemStoreKey],
		app.TrainingKeeper, app.ShardingKeeper,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx, encodingConfig.TxConfig),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		distribution.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipUpgradeHeights),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.IBCTransferKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		authz.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper),
		feegrant.NewAppModule(appCodec, app.FeegrantKeeper, app.AccountKeeper),
		nft.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
		group.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusKeeper),
		circuit.NewAppModule(appCodec, app.CircuitKeeper),
		
		compute.NewAppModule(appCodec, app.ComputeKeeper, app.AccountKeeper, app.BankKeeper),
		storage.NewAppModule(appCodec, app.StorageKeeper, app.AccountKeeper, app.BankKeeper),
		training.NewAppModule(appCodec, app.TrainingKeeper, app.AccountKeeper, app.BankKeeper),
		reward.NewAppModule(appCodec, app.RewardKeeper, app.AccountKeeper, app.BankKeeper),
		model.NewAppModule(appCodec, app.ModelKeeper, app.AccountKeeper, app.BankKeeper),
		health.NewAppModule(appCodec, app.HealthKeeper),
		recovery.NewAppModule(appCodec, app.RecoveryKeeper),
		sharding.NewAppModule(appCodec, app.ShardingKeeper),
		validation.NewAppModule(appCodec, app.ValidationKeeper),
	)

	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		authzkeeper.ModuleName,
		feegrant.ModuleName,
		nftkeeper.ModuleName,
		group.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,
		circuittypes.ModuleName,
		
		healthtypes.ModuleName,
		recoverytypes.ModuleName,
		computetypes.ModuleName,
		storagetypes.ModuleName,
		trainingtypes.ModuleName,
		rewardtypes.ModuleName,
		modeltypes.ModuleName,
		shardingtypes.ModuleName,
		validationtypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authzkeeper.ModuleName,
		feegrant.ModuleName,
		nftkeeper.ModuleName,
		group.ModuleName,
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensustypes.ModuleName,
		circuittypes.ModuleName,
		
		healthtypes.ModuleName,
		recoverytypes.ModuleName,
		computetypes.ModuleName,
		storagetypes.ModuleName,
		trainingtypes.ModuleName,
		rewardtypes.ModuleName,
		modeltypes.ModuleName,
		shardingtypes.ModuleName,
		validationtypes.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		authzkeeper.ModuleName,
		feegrant.ModuleName,
		nftkeeper.ModuleName,
		group.ModuleName,
		vestingtypes.ModuleName,
		upgradetypes.ModuleName,
		consensustypes.ModuleName,
		circuittypes.ModuleName,
		
		computetypes.ModuleName,
		storagetypes.ModuleName,
		trainingtypes.ModuleName,
		rewardtypes.ModuleName,
		modeltypes.ModuleName,
		healthtypes.ModuleName,
		recoverytypes.ModuleName,
		shardingtypes.ModuleName,
		validationtypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterServices(app.configurator)

	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		distribution.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipUpgradeHeights),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.IBCTransferKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		authz.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper),
		feegrant.NewAppModule(appCodec, app.FeegrantKeeper, app.AccountKeeper),
		nft.NewAppModule(appCodec, app.NFTKeeper, app.AccountKeeper, app.BankKeeper),
		group.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		consensus.NewAppModule(appCodec, app.ConsensusKeeper),
		circuit.NewAppModule(appCodec, app.CircuitKeeper),
		
		compute.NewAppModule(appCodec, app.ComputeKeeper, app.AccountKeeper, app.BankKeeper),
		storage.NewAppModule(appCodec, app.StorageKeeper, app.AccountKeeper, app.BankKeeper),
		training.NewAppModule(appCodec, app.TrainingKeeper, app.AccountKeeper, app.BankKeeper),
		reward.NewAppModule(appCodec, app.RewardKeeper, app.AccountKeeper, app.BankKeeper),
		model.NewAppModule(appCodec, app.ModelKeeper, app.AccountKeeper, app.BankKeeper),
		health.NewAppModule(appCodec, app.HealthKeeper),
		recovery.NewAppModule(appCodec, app.RecoveryKeeper),
		sharding.NewAppModule(appCodec, app.ShardingKeeper),
		validation.NewAppModule(appCodec, app.ValidationKeeper),
	)

	app.sm.RegisterStoreDecoders()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

func (app *AtlasApp) Name() string { return app.BaseApp.Name() }

func (app *AtlasApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *AtlasApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *AtlasApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *AtlasApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *AtlasApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

func (app *AtlasApp) BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	return modAccAddrs
}

func (app *AtlasApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *AtlasApp) AppCodec() codec.Codec {
	return app.appCodec
}

func (app *AtlasApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

func (app *AtlasApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	keys := app.GetKVStoreKeys()
	return keys[storeKey]
}

func (app *AtlasApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	keys := app.GetTransientStoreKeys()
	return keys[storeKey]
}

func (app *AtlasApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	keys := app.GetMemStoreKeys()
	return keys[storeKey]
}

func (app *AtlasApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *AtlasApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *AtlasApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
}

func (app *AtlasApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *AtlasApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".atlas")
}

