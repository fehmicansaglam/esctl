package constants

const (
	EntityNode    = "node"
	EntityIndex   = "index"
	EntityShard   = "shard"
	EntityNodes   = "nodes"
	EntityIndices = "indices"
	EntityShards  = "shards"
)

const (
	ShardStateStarted      = "STARTED"
	ShardStateRelocating   = "RELOCATING"
	ShardStateInitializing = "INITIALIZING"
	ShardStateUnassigned   = "UNASSIGNED"
)

const (
	ShardPrimary = "p"
	ShardReplica = "r"
)
