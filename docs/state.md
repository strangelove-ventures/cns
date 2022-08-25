# State

## ICA authentication module

ICA controller keeper is in charge of storing the association between an interchain account address and the chosen address on the controller chain (in our case, a group policy address).

```go
k.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.ConnectionId, policyAddress, "")
```

```go
Group policy address -> Interchain account address
```

## Chain name service

### Questions to answer

- Which IBC tokens on this chain correspond to real assets? (i.e. “which ibc/HASH is JUNO on Osmosis?”)

```protobuf
message ChainAsset {
	uint64 chainID; // An int that represents, for example, Osmosis
	uint64 assetID; // An int that represents the JUNO token
	string path; // transfer/channel-42/ujuno
}
```

Each controller chain maintains a list of IBC tokens by mapping the IBC paths on their chain with the asset IDs known to CNS.

- If I want to send a token over IBC from chain X to chain Y, which path (channels) should I use?

```protobuf
message ChainClient {
	uint64 chainID; // Chain, on which the client is opened. ID on CNS
	uint64 counterpartyChainID; // Chain, for which the client is opened. ID on CNS
	string client; // Client name
}
```

Let’s say you need to send JUNO to Osmosis.

- On CNS look up a list of assets maintained by Osmosis, find the JUNO token, find the path: `transfer/channel-42/ujuno`
- On the Osmosis chain [look up `channel-42`](https://rest.cosmos.directory/osmosis/ibc/core/channel/v1/channels/channel-42/ports/transfer) and find the counterparty channel ID: `channel-0` and the [client ID](https://rest.cosmos.directory/osmosis/ibc/core/channel/v1/channels/channel-42/ports/transfer/client_state): `07-tendermint-1457`
- On CNS using the client ID identify the chain and get RPC endpoints.
- Send tokens from Juno’s RPC endpoint using `channel-0`.

## State

```protobuf
Chains: [Network ID] -> []uint64 Chain IDs
Network: [Network ID] -> Network
Chain Name: [Chain Name] -> uint64 Network ID
```

## Types

```protobuf
message Network { // Network has many chains
	uint64 id = 1;
	string owner = 4; // ICA or regular account address
	bool verified = 6; // Updated by the group on the controller chain
	string verifiedDate = 7; // Updated by the group on the controller chain
}

message NetworkDetails {
	uint64 mainnet; // ID of the mainnet chain
	uint64 slip44;
	bytes metadata;
	repeated Algorithm algorithms;
}

enum Algorithm {
	SECP256K1 = 0;
	ETHSECP256K1 = 1;
	ED25519 = 2;
	SR25519 = 3;
}

message Chain { // Chain belongs to a network
	uint64 id;
	uint64 networkID;
	ChainDetails details;
}

message ChainDetails {
	unit64 stakingAssetID = 1;
	string chainID = 2;
	string description = 3;
	bytes metadata = 4;
	Prefix prefix = 5;
	Genesis genesis = 6;
	Gas gas = 7; // Gas multiplier
	Sourcecode sourcecode = 8;
	Status status = 9;
	repeated uint64 fees = 10; // Asset IDs of tokens, accepted as fees
	repeated Peer peers = 11;
	repeated Peer seeds = 12;
  repeated API apis = 13;
  repeated Explorer explorers = 14;
	repeated uint64 assetsNative = 15;
	repeated ChainAssets assets = 16;
	repeated ChainClients clients = 17;
}

message Prefix {
	string accAddr = 1;
	string accPub = 2;
	string valAddr = 3;
	string valPub = 4;
	string consAddr = 5;
	string consPub = 6;
}

message Genesis {
	string url = 1;
	string hash = 2;
}

message Gas {
	uint64 low = 1;
	uint64 average = 2;
	uint64 high = 3;
}

message Sourcecode {
	string url = 1;
	string hash = 2; // Commit hash
	string daemon = 3;
	string home = 4;
	Version version = 5;
	repeated Executable executables = 6;
	repeated Module modules = 7; // Includes Cosmos SDK, Tendermint Core, CosmWasm and others
  // Prerequisites, maybe? Like system requirements - a nice to have
}

message Version {
	string recommended = 1;
	repeated compatible = 2;
}

message Executable {
	Arch arch = 1;
	OS os = 2;
	string url = 3;
}

enum Arch {
	AMD64 = 0;
	ARM64 = 1;
	ARMV6 = 2;
}

enum OS {
	LINUX = 0;
	DARWIN = 1;
	WINDOWS = 2;
}

message Module {
	string url = 1; // This is a source code URL. Acts as an identifier
	string version = 2;
	bool enabled = 3;
}

enum Status {
	LIVE = 0;
	UPCOMING = 1;
	KILLED = 2;
}

message Explorer {
  string name = 1;
  string url = 2;
}

message Peer {
  string id = 1;
  string address = 2;
  string provider = 3;
}

message API {
	APIType type = 1;
  string address = 1;
  string provider = 2;
  bool archive = 3; // default: false
}

enum APIType {
	RPC = 0;
	API = 1;
	GPRC = 2;
}

message ChainAsset {
	uint64 assetID = 1;
	string path = 2;
}

message ChainClient {
	uint64 counterpartyChainID = 1; // Chain, for which the client is opened
	string client = 2; // Client name
}

message Asset {
	uint64 id = 1;
	uint64 chainID = 1;
	string description = 1;
	string symbol = 2;
	string address = 3;
	repeated Denom aliases = 4;
	Denom base = 5;
	AssetType type = 6;
	Logo logo = 7;
}

message Denom {
	string name = 1;
	uint64 decimal = 2;
}

enum AssetType {
  SDK = 0;
  CW20 = 1;
  SNIP20 = 2;
  ERC20 = 3;
}

message Logo {
  string png = 1;
  string svg = 2;
}
```