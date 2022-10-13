# Overview

## Synopsis

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

## Definitions

- A “host” chain: a chain on which the CNS module is deployed. This is the chain, where the information about other chains is stored.
- A “controller” chain: a chain, information about which is stored on the host chain. For example, if CNS is deployed on the Cosmos Hub (the “host” chain), other chains like Osmosis, Crescent and others will be considered “controller” chains, because they will store information about themselves on the “host” chain.
- Chain’s governance: a mechanism implemented in the standard Cosmos SDK `gov` module that allows holders of the governance token to vote on proposals. For CNS purposes the controller chain’s governance will act as an owner of the corresponding name in the host chain’s CNS.
- A group: a mechanism implemented in the standard Cosmos SDK `group` module that allows a set of accounts to vote on proposals. Group members are selected by the group’s admin. For CNS purposes the group acts on behalf of the controller chain’s governance to make decisions about what information about the controller’s chain goes into CNS. For example, updating the genesis URL or adding an explorer to the list is under the authority of the group.
- Interchain account (ICA): an ICS27 account on the host chain that is allowed to make changes to the values in CNS. The interchain account is registered by the controller chain’s governance and acts on behalf of the controller chain’s group responsible for making decisions. Each controller chain has one interchain account on the host chain.
- An interchain accounts authentication module (ICA AM): a new Cosmos SDK module deployed on controller chains that ensures that only the governance of the controller chain can initiate certain actions (like creating an interchain account on the host chain) and only the particular group can perform certain other actions (like asking the interchain account to change values in CNS).
- A CNS module (CNS): a new Cosmos SDK module deployed on the host chain that controls which interchain accounts can change which values in the CNS store. The CNS module is responsible for storing information about controller chains.

## Motivation

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

## Specification

### [Overview](./docs/overview.md)

### [Messages](./docs/messages.md)

### [State](./docs/state.md)