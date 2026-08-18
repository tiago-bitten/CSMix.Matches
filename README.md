# CSMix.Matches

The match itself, from ten accepts to a final score. Event sourced.

## What this service does

Picks up where Matchmaking stops: ten players have accepted, and now there is a
match to run.

- **Creation** — participants and teams, fixed at this point.
- **Lobby** — where the ten wait between accepting and connecting.
- **Map veto** — a state machine with a strict turn order.
- **The match state machine** — the thing that says which transitions are legal
  and which are not.
- **Asking GameServers for a server**, and holding the allocation.
- **Rounds and score.**
- **Abandonment and cancellation** — the paths nobody wants but everybody hits.
- **The result**, and `MatchFinished`.

## Storage

PostgreSQL, in three roles:

| | |
|---|---|
| **Event store** | what happened, append only, the only truth |
| **Transactional outbox** | events leave in the same transaction that recorded them |
| **Read projections** | every query answered from a derived shape |

**Why event sourcing here and nowhere else.** A match is a sequence of things
that happened, in order, that nobody may edit afterwards — which is the one
shape event sourcing is actually for. An account is not; a rating is not. Using
it here and not elsewhere is a decision, not an inconsistency.

**Why an outbox.** Writing to Postgres and publishing to Kafka are two systems,
and there is no transaction across both. The outbox makes it one write, and a
relay does the second half. Without it, every crash between the two is a match
that finished and a rating that never heard about it.

## Events

| | |
|---|---|
| in | `MatchReadyToCreate` |
| out | `MatchFinished`, and the lifecycle events anything else wants to follow |

## What CSMix is

A competitive Counter-Strike platform in the shape of GamersClub or FACEIT:
players sign in with Steam, queue alone or as a party, get matched into two
teams of five, veto maps, play on a server the platform provisions for them, and
come out the other side with a rating that moved.

It is split into services on purpose, and each one owns its own storage:

| Service | Owns |
|---|---|
| [CSMix.Accounts](https://github.com/tiago-bitten/CSMix.Accounts) | identity, Steam sign-in, tokens, bans — .NET |
| [CSMix.BFF](https://github.com/tiago-bitten/CSMix.BFF) | the only service a browser talks to |
| [CSMix.Matchmaking](https://github.com/tiago-bitten/CSMix.Matchmaking) | parties, queues, forming a match, getting it accepted |
| [CSMix.Matches](https://github.com/tiago-bitten/CSMix.Matches) | the match itself, from accepted to finished |
| [CSMix.GameServers](https://github.com/tiago-bitten/CSMix.GameServers) | the machines and the server processes on them |
| [CSMix.Rating](https://github.com/tiago-bitten/CSMix.Rating) | rating, seasons, leaderboards |

The flow they form:

```
Matchmaking  ──MatchReadyToCreate──▶  Matches  ──asks──▶  GameServers
                                         │
                                         └──MatchFinished──▶  Rating
```

## How every CSMix service is put together

- **`internal/<slice>/`** — a feature is a folder, not a layer. Inside one, the
  four folders appear as they earn it: `api` (routes and handlers), `app` (the
  use cases), `domain` (the rules), `infra` (the adapters).
- **`internal/shared/`** — config, logging and the inbound edge. Duplicated
  across services rather than shared as a library: `internal/` could not export
  it anyway, and a common module would couple four deploys together.
- **No `pkg/`, no framework.** `net/http` has had method and path routing since
  Go 1.22, which is all any of these need.
- **`X-Api-Key` on every service-to-service call**, and a bearer token from
  Accounts when a call is on a player behalf. Accounts is the only issuer;
  everyone else verifies with the public key from its JWKS.

## Status

**Skeleton.** What exists is the layout, configuration that fails loudly, a
health endpoint and graceful shutdown. There is no business logic yet, and the
dependencies listed above are not wired — `go.mod` has none.

```bash
cp .env.example .env
make run
```

```bash
make test && make lint
```
