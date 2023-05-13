package constants

const (
	EntityNode    = "node"
	EntityIndex   = "index"
	EntityShard   = "shard"
	EntityNodes   = "nodes"
	EntityIndices = "indices"
	EntityShards  = "shards"
)

const ( // Shard states
	ShardStateStarted      = "STARTED"
	ShardStateRelocating   = "RELOCATING"
	ShardStateInitializing = "INITIALIZING"
	ShardStateUnassigned   = "UNASSIGNED"
)
