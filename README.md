# cRelay-go-sdk

This project serves as the official Golang SDK for [CausalityGraph](https://github.com/hetu-project/causalitygraph), providing a robust implementation of the Nostr protocol with advanced causality tracking and secure messaging capabilities. It is designed to be the primary development toolkit for building decentralized applications that require strong causality guarantees and secure communication.

It combines three key technologies:

1. Nostr (decentralized messaging protocol) and [CIP](https://github.com/hetu-project/causalitygraph/tree/main/Key) (Causality Implementation Possibilities)
2. Ethereum EIP-191 signatures (secure cryptographic signing)
3. VLC (Verifiable Logical Clock)

## Core Features

### Causality and Graph Operations
- Core causality relation and causality graph functionality
- Event templates for standardized event creation
- Tools for event handling and processing
- Benchmarking utilities for performance evaluation

### Security and Identity
- Ethereum EIP-191 signing for secure message signing
- Identity verification via cryptographic proofs
- Key generation and public/private key conversion
- Causality Key protocol extensions

### Protocol Extensions
- Tools for encoding and decoding common formats
- Event templates for standardized event creation
- Extended Nostr capacity through protocol extensions
- Comprehensive event processing utilities

## CIP Implementation

The `cip` directory implements various CIP (Common Interface Protocol) specifications:

### CIP-01: Basic Subspace Operations
- Subspace creation and management
- Basic event operations
- Authentication and authorization

### CIP-02: Common Graph Operations
- Project and task management
- Entity and relation tracking
- Observation recording
- Structured data organization

### CIP-03: Model Graph Operations
- Model versioning and tracking
- Dataset management
- Compute and algorithm operations
- Validation and fine-tuning
- Conversation and session handling

### CIP-04: OpenResearch Operations
- Research paper submission and indexing
- Paper annotations and reviews
- AI analysis integration
- Research discussions
- Academic collaboration features

Each CIP implementation follows a consistent pattern:
- Event type definitions
- Tag structure specifications
- Event creation and parsing
- Subspace operation management

## Examples

### Basic Event Publishing
```go
// See example/publish/publish.go
```

### Subspace Operations
```go
// See example/subspace/subspace.go
```

### Common Graph Usage
```go
// See example/common_graph/common_graph.go
```

### Model Graph Usage
```go
// See example/model_graph/model_graph.go
```

### OpenResearch Usage
```go
// See example/openresearch/openresearch.go
```

## Security Considerations

- Private Key Safety: Never expose your Ethereum private key
- Relay Trust: Use trusted Nostr relays
- Message Verification: Validate EIP-191 signatures
- Clock Synchronization: Maintain accurate vector clocks