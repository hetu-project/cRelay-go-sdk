package cip

const (
	// Subspace common event kinds
	KindSubspaceCreate = 30100
	KindSubspaceJoin   = 30200

	// Governance event kinds
	KindGovernancePost    = 30300
	KindGovernancePropose = 30301
	KindGovernanceVote    = 30302
	KindGovernanceInvite  = 30303
	KindGovernanceMint    = 30304

	// Modelgraph event kind
	KindModelgraphModel   = 30404
	KindModelgraphData    = 30405
	KindModelgraphCompute = 30406
	KindModelgraphAlgo    = 30407
	KindModelgraphValid   = 30408
)

const (
	// General base operation types
	OpSubspaceCreate = "subspace_create" // 30100
	OpSubspaceJoin   = "subspace_join"   // 30200

	// Governance operation types (governance operations)
	OpPost    = "post"    // 30300
	OpPropose = "propose" // 30301
	OpVote    = "vote"    // 30302
	OpInvite  = "invite"  // 30303
	OpMint    = "mint"    // 30304

	// Business operation types
	OpModel   = "model"   // 30404
	OpData    = "data"    // 30405
	OpCompute = "compute" // 30406
	OpAlgo    = "algo"    // 30407
	OpValid   = "valid"   // 30408
)

const (
	// Default operations string for subspace creation
	DefaultSubspaceOps = "post=30300,propose=30301,vote=30302,invite=30303,mint=30304"

	// Modelgraph operations string for subspace creation
	ModelGraphSubspaceOps = "model=30404,data=30405,compute=30406,algo=30407,valid=30408"
)
