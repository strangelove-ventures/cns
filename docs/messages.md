# Messages

## ICA authentication module (on a controller chain)

### Register an interchain account

```protobuf
message MsgRegisterAccount {
	string address; // admin address
	cosmos.base.v1beta1.Coin amount; // to send to the policy address
}
```

Authority: controller’s governance.

Fail conditions:

- The message doesn’t originate from a governance proposal.
- The community pool doesn’t have the `msg.amount` tokens.

State modifications:

- Creates a new group where admin address is `msg.address`.
- Creates a new group policy.
- Set the global `group` value to the group ID.
- Sends an IBC packet to the host chain to register a specific interchain account.
- Stores an association between a group policy address and a registered interchain account address.

### Verify network information

```protobuf
message MsgVerifyNetwork {
	string owner; // owner address (ICA or regular account)
	uint64 networkID;
	bool verified;
}
```

Authority: controller’s governance.

Fail conditions:

- On the controller chain:
    - The message doesn’t originate from a governance proposal.
- On the host chain:
    - A network with Network ID that matches `msg.networkID` doesn't exist

State modifications:

- On the host chain:
    - Get `Network` by `msg.networkID`. Set `verified` to `msg.verified`.
    - Set `verifiedDate` to the current block’s timestamp.

### Change group admin

> TODO: think of a better message name.

```protobuf
message MsgChangeAdmin {
  string admin;
}
```

Fail conditions:

- `msg.admin` is the current group admin address

State modifications:

- Set the admin of the group (with the group ID that matches the value stored in the global `group` in the state) to `msg.admin`.

> TODO: should this message also kick out all the members of the group? Technically, this is the responsibility of the admin, but this message will likely be used when the governance loses the trust in the admin and the members (that the admin selected), so it might makes sense to rmeove all members from the group.

## Chain naming service module (on the host chain)

### Create chain

```protobuf
message MsgCreateChain {
	string owner; // owner address (ICA or regular account)
	uint64 networkID;
	ChainDetails details;
}
```

Authority: controller’s group.

Fail conditions:

- `msg.owner` doesn’t match the owner of the `Network` with `msg.networkID`
- A network with Network ID that matches `msg.networkID` doesn't exist

State modifications:

- Create a new chain entry.

### Update chain

```protobuf
message MsgUpdateChain {
	string owner; // owner address (ICA or regular account)
	uint64 chainID;
	ChainDetails details;
}
```

Fail conditions:

- Get a network that the chain with `msg.chainID` belongs to. Fail if `msg.owner` doesn’t match the owner of the network.
- A chain with ID that matches `msg.chainID` doesn't exist.

State modifications:

- Update the chain with `msg.details`.

### Create network

```protobuf
message MsgCreateNetwork {
	string owner; // owner address (ICA or regular account)
	NetworkDetails details;
}
```

Fail conditions:

- N/A.

State modification:

- Create a new `Network`.
    - Set `[id](http://Network.id)` as an incrementing integer
    - Set `owner` as `msg.owner`
    - Set `name` as an empty string
    - Set `verified` as `false`

### Update network

```protobuf
message MsgUpdateNetwork {
	string owner; // owner address (ICA or regular account)
	uint64 networkID;
	NetworkDetails details;
}
```

Fail conditions:

- `msg.owner` doesn’t match the owner of the `Network` with `msg.networkID`
- A network with Network ID that matches `msg.networkID` doesn't exist

State modifications:

- Set `NetworkDetails`

### Set host’s group


> ⚠️ Think of a better name

```protobuf
message MsgSetGroup {
	string address;
}
```

Authority: host’s governance

Fail conditions:

- N/A

State modifications:

- TODO!

### Assign network name

```protobuf
message MsgAssignNetworkName {
	string address; // Cosmos Hub group policy address
	uint64 networkID;
	string name;
}
```

Authority: host’s group.

Fail conditions:

- A network with Network ID that matches `msg.networkID` doesn't exist