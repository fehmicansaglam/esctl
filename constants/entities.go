package constants

const (
	EntityNode    = "node"
	EntityNodes   = "nodes"
	EntityIndex   = "index"
	EntityIndices = "indices"
	EntityShard   = "shard"
	EntityShards  = "shards"
	EntityAlias   = "alias"
	EntityAliases = "aliases"
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
