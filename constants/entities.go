package constants

const (
	EntityNode    = "node"
	EntityNodes   = "nodes"
	EntityIndex   = "index"
	EntityIndices = "indices"
	EntityShards  = "shards"
	EntityAliases = "aliases"
	EntityTasks   = "tasks"
	EntityCluster = "cluster"
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
