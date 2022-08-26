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

```
Network/[Network ID uint64] -> Network
Chain/[Network ID]-[Chain ID] -> Chain
ChainName/[Chain Name string] -> uint64 Network ID
Asset/[Asset ID] -> Asset
Group -> uint64
```

### Network

```protobuf
message Network { // Network has many chains
	uint64 id = 1;
	string owner = 4; // ICA or regular account address
	bool verified = 6; // Updated by the group on the controller chain
	string verifiedDate = 7; // Updated by the group on the controller chain
  NetworkDetails details = 8;
}
```

A `Network` groups all the testnets and the mainnet (if exists) of a specific blockchain network and defines the common properties.

`id` is an incrementing integer used to identify a network in CNS. `id` is used for internal purposes and not to identify the network on the Interchain. `id` is not the same as the chain name or a Tendermint Core's chain ID.

`owner` is an address of an interchain account or a regular account that created the network and has the authority to modify the values associated with the network. The `owner` is set automatically when the network is first created.

`verified` is a boolean value and `verifiedData` is a string that represents a block timestamp. Verification is a process when the governance of the controller chain verifies that the information in CNS for the blockchain network they represent is correct and up to date. Both values can only be changed by governance of the controller chain through `MsgVerifyNetwork`.

`details` define information common to all blockchain in the network.

### Network Details

```protobuf
message NetworkDetails {
	uint64 mainnet; // ID of the mainnet chain
	uint64 slip44;
	repeated Algorithm algorithms;
	bytes metadata;
}
```

Network details contains information that can be edited by the `owner` address of network.

`mainnet` represents the ID of the current active mainnet. The value is an integer that identifies a `Chain` in CNS. It is not a Tendermint Core's chain ID, but rather an internal to CNS integer ID.

> TODO: What should the value be if a network doesn't yet have a mainnnet? `0` maybe?

`slip44` is an integer that identifies the coin type of the native asset of the chain. For example, for Cosmos Hub's ATOM the value would be `118` as per [SLIP-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md).

`algorithms` are parameters of the elliptic curve used by the network's cryptography.

```protobuf
enum Algorithm {
	SECP256K1 = 0;
	ETHSECP256K1 = 1;
	ED25519 = 2;
	SR25519 = 3;
}
```

### Chain

```protobuf
message Chain { // Chain belongs to a network
	uint64 id;
	ChainDetails details;
}
```

A `Chain` is either a testnet or a mainnet within a `Network`. A network can have any number of testnets and a single mainnet. A chain is assumed to be a testnet, unless its ID is specified in `Network.NetworkDetails.mainnet`.

### Chain Details

```protobuf
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
```

`stakingAssetID` is a CNS ID of the asset used for staking purposes on the chain.

`chainID` is the Tendermint Core's `chain-id` value. This value does not identify the chain on CNS or on the interchain, but rather is used by Tendermint Core and Cosmos SDK to distinguish between chains within the same network.

`description` is a human readable description of the chain. For example, if a testnet was launched specifically for a hackathon, this might be mentioned in the description. Or if a particular testnet is the main one that developers should use to test their applications.

`metadata` may contain any additional data about the chain that doesn't belong in other fields. It is recommended to use JSON encoded data. If enough chains use the same metadata fields, the fields may be added to CNS.

`prefix` is a bech32 human-readable prefix used to denote an address type. Refer to Cosmos SDK documentation to learn more about how different address types are used.

```protobuf
message Prefix {
	string accAddr = 1;
	string accPub = 2;
	string valAddr = 3;
	string valPub = 4;
	string consAddr = 5;
	string consPub = 6;
}
```

`genesis` defines the genesis file of the chain. It is useful for node operators who need a genesis file to be able to start a node. `url` 

```protobuf
message Genesis {
	string url = 1;
	string hash = 2;
}
```
`gas` is the amount of gas that is recommended to use when broadcasting transactions for this particular chain. The gas amount is dependent on the type of messages included in a transaction, so a medium gas amount consumed for transactions with bank's `MsgSend` should be provided.

```protobuf
message Gas {
	uint64 low = 1;
	uint64 average = 2;
	uint64 high = 3;
}
```

`sourcecode` contains the information on how to get the chain's binary executable file.

```protobuf
message Sourcecode {
	string url = 1;
	string hash = 2;
	string daemon = 3;
	string home = 4;
	Version version = 5;
	repeated Executable executables = 6;
	repeated Module modules = 7; // Includes Cosmos SDK, Tendermint Core, CosmWasm and others
  // Prerequisites, maybe? Like system requirements - a nice to have
}
```

`url` is the URL of the repository with the source code of the chain. The URL should be in a format that can be used with `git clone [URL]`.

`hash` is the commit hash of the current version of the software.

`daemon` is the name of the executable file for the blockchain node.

`home` is the default data directory for the blockchain data.

`version` defines both recommended (current) and all the versions, compatible with the current one.

```protobuf
message Version {
	string recommended = 1;
	repeated string compatible = 2;
}
```
`executables` contains a list of `Executable` that define URLs of binaries for different OSes and architectures.

```protobuf
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
```

`modules` contains a list of components that exist on the blockchain. Those include Cosmos SDK modules, Tendermint Core, IBC modules, possibly even CosmWasm contracts. This field is used to guide dapp developers on what is available on the chain (does it have a DEX, CosmWasm, etc.).

```protobuf
message Module {
	string url = 1;
	string version = 2;
	bool enabled = 3;
}
```

The `url` contains the URL to the source code of the module (e.g. https://github.com/tendermint/tendermint) and acts as an identifier.

`version` is the version of the software and `enabled` indicates whether the module is currently enabled on the blockchain.

## Types

```protobuf
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