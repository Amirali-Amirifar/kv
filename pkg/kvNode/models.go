package kvNode

import (
	"github.com/Amirali-Amirifar/kv/internal"
	"net"
)

type Config struct {
	ID            int
	Status        internal.NodeStatus
	Address       net.TCPAddr
	StoreNodeType internal.StoreNodeType
}
