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

	// Modelgraph event kind
	KindModelgraphModel   = 30304
	KindModelgraphData    = 30305
	KindModelgraphCompute = 30306
	KindModelgraphAlgo    = 30307
	KindModelgraphValid   = 30308
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

	// Business operation types
	OpModel   = "model"   // 30304
	OpData    = "data"    // 30305
	OpCompute = "compute" // 30306
	OpAlgo    = "algo"    // 30307
	OpValid   = "valid"   // 30308
)

const (
	// Default operations string for subspace creation
	DefaultSubspaceOps = "post=30300,propose=30301,vote=30302,invite=30303"
	
	// Modelgraph operations string for subspace creation
	ModelGraphSubspaceOps = "model=30304,data=30305,compute=30306,algo=30307,valid=30308"
)