package config

type AddressConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ClusterConfig struct {
	Partitions int `mapstructure:"partitions"`
	Replicas   int `mapstructure:"replicas"`
}

type DiscoveryConfig struct {
	HeartbeatIntervalMs int `mapstructure:"heartbeat_interval_ms"`
	FailureTimeoutMs    int `mapstructure:"failure_timeout_ms"`
}

type KvControllerConfig struct {
	Address   AddressConfig   `mapstructure:"address"`
	Cluster   ClusterConfig   `mapstructure:"cluster"`
	Discovery DiscoveryConfig `mapstructure:"discovery"`
}

type KvNodeConfig struct {
	Address       AddressConfig `mapstructure:"address"`
	Controller    AddressConfig `mapstructure:"controller"`
	MasterAddress string        `mapstructure:"master_address"`
	MasterPort    int           `mapstructure:"master_port"`
	HTTPTimeout   int           `mapstructure:"http_timeout_ms"`
	NodeID        int           `mapstructure:"node_id"`
}
