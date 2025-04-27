package cip

// KeyOpMap maps kind values to operation names
var KeyOpMap = map[int]string{
	// common operations
	KindSubspaceCreate: OpSubspaceCreate,
	KindSubspaceJoin:   OpSubspaceJoin,

	// Governance operations
	KindGovernancePost:    OpPost,
	KindGovernancePropose: OpPropose,
	KindGovernanceVote:    OpVote,
	KindGovernanceInvite:  OpInvite,

	// ModelGraph operations
	KindModelgraphModel:   OpModel,
	KindModelgraphData:    OpData,
	KindModelgraphCompute: OpCompute,
	KindModelgraphAlgo:    OpAlgo,
	KindModelgraphValid:   OpValid,
}

// GetOpFromKind returns the operation name for a given kind value
func GetOpFromKind(kind int) (string, bool) {
	op, exists := KeyOpMap[kind]
	return op, exists
}

// GetKindFromOp returns the kind value for a given operation name
func GetKindFromOp(op string) (int, bool) {
	for kind, operation := range KeyOpMap {
		if operation == op {
			return kind, true
		}
	}
	return 0, false
}
