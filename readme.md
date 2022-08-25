# Overview

# Synopsis

Chain naming service (CNS) is a system that allows blockchains (“controller chains”) in Cosmos to update information about themselves on a blockchain (“host chain”).

[![](/assets/map.png)](https://whimsical.com/update-information-McDGXQ1W8gNxdoEAJoKhbT)

Any blockchain can act as a controller or a host chain, but for the purposes of this document, we will assume CNS is deployed on a single host chain and many controller chains.

CNS on the host chain stores the following information about each controller chain:

- Chain name
- Network type
- A list of explorers
- Genesis URL
- API URLs and more.

Currently, the information that will be managed by CNS is stored in a repository on GitHub: [https://github.com/cosmos/chain-registry](https://github.com/cosmos/chain-registry)

# Definitions

- A “host” chain: a chain on which the CNS module is deployed. This is the chain, where the information about other chains is stored.
- A “controller” chain: a chain, information about which is stored on the host chain. For example, if CNS is deployed on the Cosmos Hub (the “host” chain), other chains like Osmosis, Crescent and others will be considered “controller” chains, because they will store information about themselves on the “host” chain.
- Chain’s governance: a mechanism implemented in the standard Cosmos SDK `gov` module that allows holders of the governance token to vote on proposals. For CNS purposes the controller chain’s governance will act as an owner of the corresponding name in the host chain’s CNS.
- A group: a mechanism implemented in the standard Cosmos SDK `group` module that allows a set of accounts to vote on proposals. Group members are selected by the group’s admin. For CNS purposes the group acts on behalf of the controller chain’s governance to make decisions about what information about the controller’s chain goes into CNS. For example, updating the genesis URL or adding an explorer to the list is under the authority of the group.
- Interchain account (ICA): an ICS27 account on the host chain that is allowed to make changes to the values in CNS. The interchain account is registered by the controller chain’s governance and acts on behalf of the controller chain’s group responsible for making decisions. Each controller chain has one interchain account on the host chain.
- An interchain accounts authentication module (ICA AM): a new Cosmos SDK module deployed on controller chains that ensures that only the governance of the controller chain can initiate certain actions (like creating an interchain account on the host chain) and only the particular group can perform certain other actions (like asking the interchain account to change values in CNS).
- A CNS module (CNS): a new Cosmos SDK module deployed on the host chain that controls which interchain accounts can change which values in the CNS store. The CNS module is responsible for storing information about controller chains.

# Motivation

The advantages of managing such information on-chain:

- Name ownership by the governance: the governance registers and owns the name and the associated information on CNS permissionlessly.
- Delegation of authority: to streamline the process of updating the information in the CNS the governance can delegate this responsibility to a group of accounts.
- Data availability: the information is stored on hundreds of nodes of the host chain and is always available for end clients to use through their APIs.

|                      | Github                                                                                                                                                                | Blockchain                                                                                                                                                   |
| -------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Governance           | ❌ Admins can kick each other out of the org. Github admin accounts don’t map well to on-chain accounts and can’t represent the will of the community (token holders). | ✅ On-chain governance delegates responsibility to a group. Chains will own information about themselves. A group on Cosmos Hub will be able to assign names. |
| Censorship           | ❌ The whole system can be shut down unilaterally 🌪️ Any given system is only as decentralized as its least decentralized component.                                    | ✅ It’s a blockchain.                                                                                                                                         |
| On/off-chain data    | ❌ To build blockchains apps that know how to route tokens and make cross-chain requests requires off-chain components (to fetch data) or oracles.                     | ✅ All the data required for advanced cross-chain apps is available on-chain.                                                                                 |
| UI                   | 🆗 Basic editing UI exists. Additional UI is necessary.                                                                                                                | ❌ The UI has to be built.                                                                                                                                    |
| Frontend ease-of-use | ✅ Can be published as a library on npm.                                                                                                                               | ✅ Can be queried dynamically.                                                                                                                                |
| Adds utility to ATOM | ❌                                                                                                                                                                     | ✅                                                                                                                                                            |
| Current adoption     | ❌ Most end-user apps don’t use the chain registry. Osmosis uses a mix of Keplr-style registry and cosmos/chain-registry.                                              | ❌                                                                                                                                                            |
| Data management      | ✅ GH Actions can provide some data consistency. ❌ not possible to do more granular access control (controller chain can edit X, host chain can edit Y).               | ✅ structure can be enforced.                                                                                                                                 |
| Bech32 routing       | ❌ programmatically relying on a centralized service for token sending may lead to loss of funds. Needs an off-chain component to work.                                | ✅ each chain can have its own unique prefix.                                                                                                                 |
| Usernames            | ❌ nobody’s going to want to have a username in a file on Github.                                                                                                      | ✅ Market for second-level names for users: username.chainname. For example, denis.cosmos.                                                                    |

### Chain names assignment

For the initial version, chain names are assigned by a group on the Hub manually. Chain names are not purchased and there is no market for them. The reason is there is very little sense in transferring a name from one chain to another. For example, if the “real” Osmosis was assigned a name, what are the chances that Osmosis would need to sell the name to someone else, or even transfer it?

It makes sense for a group on the Hub to be able to unassign a name from a chain, for example, when a controller chain forks and the Hub needs to choose which one is the “real” one. Giving the ability to a controller chain to sell its name creates more problems.

In the next version of CNS usernames could be introduced. Think of chain names as top-level domains and usernames as regular second-level domains.

So, Osmosis, for example, would have a chain name `osmosis`, and users would be able to purchase names like `alice.osmosis`. These names could be NFTs and there could (and should) be a market for them. How these usernames are resolved to addresses and chains is outside of the scope of this document.

# Overview

A high-level overview of the process:

1. On the controller chain:
    1. The governance decides on which group will be authorized to submit changes to CNS and registers an interchain account (ICA) on the host chain.
    2. The group’s admin can choose group members.
    3. Group members vote on proposals that contain messages that are sent via the ICA AM over IBC to the host chain to be broadcasted by the chain’s ICA account to the CNS module.
    4. If the governance loses confidence in the group’s decisions, they can authorize a different group and invalidate the data currently in the CNS module of the host chain.
2. On the host chain:
    1. CNS module stores the data about chain names.
    2. The governance chooses a group that will be authorized to associate chain info (provided by the controller chain) with a chain name and an IBC client.
    3. The group assigns chain names and IBC clients to chain info.

Interacting with CNS is possible through:

- an interchain account. With an ICA a controller chain’s governance owns the record in CNS. This, however, requires a controller chain to have a ICA controller chain functionality enabled, the group module, and a CNS authentication module installed.
- a regular account. This is an alternative non-IBC way of interacting with CNS. Useful for chains that don’t have the required modules installed or are not IBC-enabled.

### High-level technical overview

Controller networks through governance and a group own (and can modify) information about themselves. The host chain through a group assigns names to networks.

A network can have many chains. Chains can be either testnets or mainnets. A network can have only one mainnet.

Each asset on the interchain has a unique ID. For example, the ATOM token could have an ID of 1. An asset “object” contains descriptive information about the asset.

Each chain maintains a list of mappings between asset IDs and paths. For example, on the Osmosis network, the mainnet chain in their list of assets can have a mapping between asset # 1 (unique global ATOM ID) and a path “transfer/channel-0/uatom”. This path is what Osmosis’ mainnet considers to be the real ATOM.

Each chain also maintains a list of mappings between counterparty chain IDs and IBC clients. For example, on the Osmosis network, the mainnet chain can have a mapping between `07-tendermint-1457` and the chain ID of Cosmos Hub (let’s say 2).

ATOM → transfer/channel-0/uatom → 07-tendermint-1457 → Cosmos Hub ID → endpoints — You need CNS to resolve ATOM to endpoints.

ATOM → transfer/channel-0/uatom → channel-0 (counterparty channel) — You need CNS to resolve ATOM to IBC path.

[![](/assets/high-level.png)](https://whimsical.com/KN18gj3PgsziNvke2vXY1y)

### Updating CNS information

Updating information in the CNS module is important for keeping information about controller chains up to date. The responsibility of updating is delegated by the controller chain’s governance to a specific group. This makes the decision-making process more agile because every change doesn’t have to go through a slow-moving governance process.

[![](/assets/updating.png)](https://whimsical.com/update-cns-object-values-7sAa35dBMxv41NDZxNPcYw)

Process:

1. On the controller chain:
    1. A group member creates a proposal that contains a message that will be routed to ICA AM. The contents of the message describe the changes that need to be committed to the CNS module.
    2. The group votes on the proposal
    3. If the proposal passes, ICA AM processes the message
    4. ICA AM checks that the message from the proposal was created by the group policy with an address that matches the group policy address stored in ICA AM.
    5. If the group ID matches, an IBC packet is sent to the host chain.
2. On the host chain:
    1. ICS27 module processes the packet and instructs the associated with the controller chain interchain account to broadcast a transaction with a message that contains the required changes to the CNS.
    2. CNS checks that the interchain account has the authority to change the values of a particular chain name associated with the controller chain.

### Registering an interchain account

The information about chains in CNS can be updated by interchain accounts. Before any data about a controller chain can be written to CNS, an interchain account has to be created. An interchain account can only be created as a result of a governance vote on the controller chain.

[![](../assets/registering-ica.png)](https://whimsical.com/VFCP16ja6B7kvkTSXhGW89)

Process:

1. On the controller chain:
    1. A governance proposal to create an interchain account is submitted. The proposal includes a message that will be routed to the ICA authentication module (ICA AM).
    2. If a proposal passes, the message is processed by ICA AM. The message is executed if and only if it originates from the governance module (similarly to how the upgrade module works). The message contains an admin address and an amount of tokens that will be sent to the policy address from the community pool.
    3. ICA AM performs the following actions:
        1. Creates a new group with the admin address submitted in the message.
        2. Creates a new group policy.
        3. Sends an IBC packet to the host chain to register a specific interchain account.
        4. Stores an association between a group policy address and a registered interchain account address.
2. On the host chain:
    1. The interchain accounts module receives a packet and registers an interchain account.

Registering an ICA account is only possible through the process above.

## Verify the information in CNS is up-to-date and valid

The controller chain’s governance delegated the responsibility to update information in CNS to a group. However, it is important for the governance of the chain to periodically verify that the information in CNS is indeed correct and up-to-date or invalid. Invalidating might be useful in case the group admin or the group itself becomes malicious and the governance needs to communicate that the info in CNS is no longer valid.

Process:

1. On the controller chain:
    1. A governance proposal to validate CNS info is submitted. The proposal includes a message that will be routed to the ICA authentication module (ICA AM).
    2. If a proposal passes, the message is processed by ICA AM. The message is executed if and only if it originates from the governance module. The message contains a boolean that describes if the information in CNS is valid or not.
    3. ICA AM sends a packet to the host chain.
2. On the host chain:
    1. The interchain accounts module receives a packet and orders the associated interchain account to broadcast a transaction
    2. The interchain account broadcasts a transaction that updates CNS with the information about whether the chain info is valid or not.

User interfaces can show this as an additional bit of information about a controller chain: “Information is valid as of Jul, 25th 2022”.

```protobuf
message Network {
	// ...
	bool verified;
	string verifiedDate;
}
```

## Verify that the information in CNS is associated with a particular chain

The controller chain (the governance and the group) has the authority to change information in CNS associated with that chain. However, a group on the host chain should associate the information provided by the controller chain with a specific name and an IBC client. This is important because two chains can submit information claiming they are the “Foo” chain and it’s up to a group on the host chain to decide which one is the one and only “Foo” chain.

Process:

1. On the host chain:
    1. A group votes on a proposal that includes a message for the CNS module.
    2. If the proposal passes, the CNS module checks that the message was broadcasted by the authorized group. If so, executes the message and update the values in CNS: associate chain info with a chain name and an IBC client.

User interfaces can show this as an additional bit of information about a controller chain: “Osmosis ✅”.

Verification from the host’s side is happening by assigning a name to the network. Only a selected group on the host chain can make the decision to assign a name to a network. If a group wants to revoke the name, they can do so.

```protobuf
message Network {
	// ...
	string name;
}
```

<aside>
❓ Maybe there should be a period after which the name change happens. Let’s say after the host’s group decided to unassign a name, a countdown starts (a week), during which the controller chain can make changes the required changes/updates. If the controller’s group satisfies the requirements of the host’s group, the host’s group can send another tx to assign a name. Name assignment happens instantly.

</aside>

# Architecture

The CNS system is implemented as two Cosmos SDK modules:

- ICA authentication module (ICA AM)
- CNS module

## ICA authentication module (ICA AM)

The ICA authentication module is deployed on controller chains. The purpose of ICA AM is to ensure that only governance can send certain messages (register an ICA account) and that only a particular group can send certain other messages (update CNS values).

## CNS module

The CNS module is deployed on the host chain. The main purpose of the CNS module is to store the association between chain names and chain data of controller chains and to control which interchain account on the host chain is authorized to update the corresponding chain data.

# User archetypes

## Interchain wallet user

“I want to send tokens from Cosmos Hub to Osmosis”.

Cosmos Hub and Osmosis are chain names that are resolved to chain details using CNS.

Type: indirect novice user.

## Interchain wallet developer

Builds a wallet that integrates with CNS.

“When I’m building a token transfer UI, which API endpoints should I use if a user wants to send tokens from Cosmos Hub to Osmosis?” — queries for a list of all verified chains. A response should include mappings between a human-readable chain name and chain details.

“which IBC channel should I use between chains?” — queries for a list of verified channels between chains.

Type: experienced developer.

## Controller governance voter

Type: experienced/novice user.

Actions:

- Vote on which group policy is responsible for making changes to the CNS
- Vote on verifying the data on CNS

Frequency of use: every 4-8 weeks.

## Controller group admin

Type: experienced user.

Actions:

- Select group members
- Update group/policy details

Frequency of use: every 4-8 weeks.

## Controller group member

Type: experienced user.

Actions:

- Proposes changes to the CNS
- Votes on proposals

Frequency of use: every 2-3 weeks.

## Host governance voter

Type: experienced/novice user.

Actions:

- Vote on which group policy is responsible for making changes to the CNS
- Vote on verifying the data on CNS

Frequency of use: every 4-8 weeks.

## Host group admin

Type: experienced user.

Actions:

- Select group members
- Update group/policy details

Frequency of use: every 4-8 weeks.

## Host group member

Type: experienced user.

Actions:

- Vote on assigning a name to a chain

Frequency of use: every 2-3 weeks.