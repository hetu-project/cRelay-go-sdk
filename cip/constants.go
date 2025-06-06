package cip

// Event kinds
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

	// CommonGraph event kinds
	KindCommonGraphProject     = 30101
	KindCommonGraphTask        = 30102
	KindCommonGraphEntity      = 30103
	KindCommonGraphRelation    = 30104
	KindCommonGraphObservation = 30105

	// Modelgraph event kind
	KindModelgraphModel        = 30404
	KindModelgraphDataset      = 30405
	KindModelgraphCompute      = 30406
	KindModelgraphAlgo         = 30407
	KindModelgraphValid        = 30408
	KindModelgraphFinetune     = 30409
	KindModelgraphConversation = 30410
	KindModelgraphSession      = 30411

	// OpenResearch event kinds
	KindOpenResearchPaper      = 30501
	KindOpenResearchAnnotation = 30502
	KindOpenResearchReview     = 30503
	KindOpenResearchAIAnalysis = 30504
	KindOpenResearchDiscussion = 30505
	KindOpenResearchReadPaper  = 30506
	KindOpenResearchCoCreate   = 30507

	// Social event kinds
	KindSocialLike     = 30600
	KindSocialCollect  = 30601
	KindSocialShare    = 30602
	KindSocialComment  = 30603
	KindSocialTag      = 30604
	KindSocialFollow   = 30605
	KindSocialUnfollow = 30606
	KindSocialQuestion = 30607
	KindSocialRoom     = 30608
	KindSocialMessage  = 30609

	// Community event kinds
	KindCommunityCreate         = 30700
	KindCommunityInvite         = 30701
	KindCommunityChannelCreate  = 30702
	KindCommunityChannelMessage = 30703
)

// Event operations
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

	// CommonGraph operation types
	OpProject     = "project"     // 30101
	OpTask        = "task"        // 30102
	OpEntity      = "entity"      // 30103
	OpRelation    = "relation"    // 30104
	OpObservation = "observation" // 30105

	// Business operation types
	OpModel        = "model"        // 30404
	OpDataset      = "dataset"      // 30405
	OpCompute      = "compute"      // 30406
	OpAlgo         = "algo"         // 30407
	OpValid        = "valid"        // 30408
	OpFinetune     = "finetune"     // 30409
	OpConversation = "conversation" // 30410
	OpSession      = "session"      // 30411

	// OpenResearch operation types
	OpPaper      = "paper"           // 30501
	OpAnnotation = "annotation"      // 30502
	OpReview     = "review"          // 30503
	OpAIAnalysis = "ai_analysis"     // 30504
	OpDiscussion = "discussion"      // 30505
	OpReadPaper  = "read_paper"      // 30506
	OpCoCreate   = "co_create_paper" // 30507

	// Social operation types
	OpLike     = "like"     // 30600
	OpCollect  = "collect"  // 30601
	OpShare    = "share"    // 30602
	OpComment  = "comment"  // 30603
	OpTag      = "tag"      // 30604
	OpFollow   = "follow"   // 30605
	OpUnfollow = "unfollow" // 30606
	OpQuestion = "question" // 30607
	OpRoom     = "room"     // 30608
	OpMessage  = "message"  // 30609

	// Community operation types
	OpCommunityCreate = "community_create" // 30700
	OpCommunityInvite = "community_invite" // 30701
	OpChannelCreate   = "channel_create"   // 30702
	OpChannelMessage  = "channel_message"  // 30703
)

// Default cip operations
const (
	// Default common project actions
	CommonPrjOps = "project=30101,task=30102"

	// Default common graph actions
	CommonGraphOps = "entity=30103,relation=30104,observation=30105"

	// Default operations string for subspace creation
	DefaultSubspaceOps = "post=30300,propose=30301,vote=30302,invite=30303,mint=30304"

	// Modelgraph operations string for model
	ModelGraphSubspaceOps = "dataset=30405,finetune=30409,conversation=30410,session=30411"

	// OpenResearch operations string
	OpenResearchSubspaceOps = "paper=30501,annotation=30502,review=30503,ai_analysis=30504,discussion=30505,read_paper=30506,co_create_paper=30507"

	// Social operations string
	SocialSubspaceOps = "like=30600,collect=30601,share=30602,comment=30603,tag=30604,follow=30605,unfollow=30606,question=30607,room=30608,message=30609"

	// Community operations string
	CommunitySubspaceOps = "community_create=30700,community_invite=30701,channel_create=30702,channel_message=30703"
)
