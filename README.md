## Overview

This project combines three key technologies:

1.  Nostr (decentralized messaging protocol)
2.  Ethereum EIP-191 signatures (secure cryptographic signing)
3.  VLC (Verifiable Logical Clock)

Together, these components enable:

-   Decentralized authentication
-   Secure message signing using Ethereum standards
-   Distributed event ordering and causality tracking
-   Message consistency in distributed systems

This integration creates a robust system for secure, decentralized communication with verifiable event ordering.

## Features

🔗 Nostr Integration

-   Decentralized communication protocol
-   Publish and subscribe to events securely
-   Relay-based message broadcasting
-   Censorship-resistant communication

🔑 Ethereum EIP-191 Signature

-   Secure message signing using Ethereum wallets
-   Identity verification via cryptographic proofs
-   EIP-191 ensures structured, tamper-proof signatures
-   Compatible with Ethereum ecosystem tools

⏰ Verifiable Logical Clock

-   Track causality between distributed events
-   Maintain consistent ordering across nodes
-   Detect concurrent operations
-   Enable distributed consensus

## How It Works

-   Nostr Events: Messages are exchanged using Nostr's relay-based architecture
-   Ethereum Signing: Messages are signed with Ethereum private keys following the EIP-191 standard
-   Verifiable Clock: Each node maintains a vector timestamp to track event causality
-   Verification: Signatures and timestamps are verified to ensure message authenticity and ordering

## Security Considerations

-   Private Key Safety: Never expose your Ethereum private key; always sign messages in a secure environment
-   Relay Trust: Use trusted Nostr relays to prevent data interception
-   Message Verification: Ensure EIP-191 signatures are correctly validated before processing messages
-   Clock Synchronization: Maintain accurate vector clocks to prevent causality violations

## Future Enhancements

-   Multiple relay support for increased redundancy
-   Advanced vector clock optimization for better scalability
-   Integration with additional signing standards
-   Enhanced privacy features
-   Performance optimizations for large-scale deployments

## Examples

### Generating a key

``` go
package main

import (
    "fmt"

    "github.com/nbd-wtf/go-nostr"
    "github.com/nbd-wtf/go-nostr/nip19"
)

func main() {
    sk := nostr.GeneratePrivateKey()
    pk, _ := nostr.GetPublicKey(sk)
    nsec, _ := nip19.EncodePrivateKey(sk)
    npub, _ := nip19.EncodePublicKey(pk)

    fmt.Println("sk:", sk)
    fmt.Println("pk:", pk)
    fmt.Println(nsec)
    fmt.Println(npub)
}
```

### Subscribing to a single relay

``` go
ctx := context.Background()
relay, err := nostr.RelayConnect(ctx, "wss://relay.stoner.com")
if err != nil {
	panic(err)
}

npub := "npub1422a7ws4yul24p0pf7cacn7cghqkutdnm35z075vy68ggqpqjcyswn8ekc"

var filters nostr.Filters
if _, v, err := nip19.Decode(npub); err == nil {
	pub := v.(string)
	filters = []nostr.Filter{{
		Kinds:   []int{nostr.KindTextNote},
		Authors: []string{pub},
		Limit:   1,
	}}
} else {
	panic(err)
}

ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
defer cancel()

sub, err := relay.Subscribe(ctx, filters)
if err != nil {
	panic(err)
}

for ev := range sub.Events {
	// handle returned event.
	// channel will stay open until the ctx is cancelled (in this case, context timeout)
	fmt.Println(ev.ID)
}
```

### Publishing to two relays

``` go
sk := nostr.GeneratePrivateKey()
pub, _ := nostr.GetPublicKey(sk)

ev := nostr.Event{
	PubKey:    pub,
	CreatedAt: nostr.Now(),
	Kind:      nostr.KindTextNote,
	Tags:      nil,
	Content:   "Hello World!",
}

// calling Sign sets the event ID field and the event Sig field
ev.Sign(sk)

// publish the event to two relays
ctx := context.Background()
for _, url := range []string{"wss://relay.stoner.com", "wss://nostr-pub.wellorder.net"} {
	relay, err := nostr.RelayConnect(ctx, url)
	if err != nil {
		fmt.Println(err)
		continue
	}
	if err := relay.Publish(ctx, ev); err != nil {
		fmt.Println(err)
		continue
	}

	fmt.Printf("published to %s\n", url)
}
```

### Logging

To get more logs from the interaction with relays printed to STDOUT you can compile or run your program with `-tags debug`.

To remove the info logs completely, replace `nostr.InfoLogger` with something that prints nothing, like

``` go
nostr.InfoLogger = log.New(io.Discard, "", 0)
```

### Example script

```
go run example/example.go
```

### Using [`libsecp256k1`](https://github.com/bitcoin-core/secp256k1)

[`libsecp256k1`](https://github.com/bitcoin-core/secp256k1) is very fast:

```
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-2400 CPU @ 3.10GHz
BenchmarkWithoutLibsecp256k1/sign-4          	    2794	    434114 ns/op
BenchmarkWithoutLibsecp256k1/check-4         	    4352	    297416 ns/op
BenchmarkWithLibsecp256k1/sign-4             	   12559	     94607 ns/op
BenchmarkWithLibsecp256k1/check-4            	   13761	     84595 ns/op
PASS
```

But to use it you need the host to have it installed as a shared library and CGO to be supported, so we don't compile against it by default.

To use it, use `-tags=libsecp256k1` whenever you're compiling your program that uses this library.

### Test for Wasm

Install [wasmbrowsertest](https://github.com/agnivade/wasmbrowsertest), then run tests:

```sh
GOOS=js GOARCH=wasm go test -short ./...
```

## Warning: risk of goroutine bloat (if used incorrectly)

Remember to cancel subscriptions, either by calling `.Unsub()` on them or ensuring their `context.Context` will be canceled at some point.
If you don't do that they will keep creating a new goroutine for every new event that arrives and if you have stopped listening on the
`sub.Events` channel that will cause chaos and doom in your program.

## Contributing to this repository

Use NIP-34 to send your patches to `naddr1qqyxwmeddehhxarjqy28wumn8ghj7un9d3shjtnyv9kh2uewd9hsz9nhwden5te0wfjkccte9ehx7um5wghxyctwvsq3vamnwvaz7tmjv4kxz7fwwpexjmtpdshxuet5qgsrhuxx8l9ex335q7he0f09aej04zpazpl0ne2cgukyawd24mayt8grqsqqqaueuwmljc`.
