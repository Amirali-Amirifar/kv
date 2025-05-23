// lib/api/kv-controller/types.d.ts

 declare namespace ApiTypes {
    // Health check
    export interface HealthCheckResponse {
        // Will contain health status information
    }

    // Node management
    export interface AddNodeRequest {
        // Node details such as host, port, etc.
    }

    export interface AddNodeResponse {
        // Added node details, ID, etc.
    }

    export interface RemoveNodeResponse {
        // Result of node removal
    }

    // Partition management
    export interface IncreasePartitionsRequest {
        // Number of partitions to add, etc.
    }

    export interface IncreasePartitionsResponse {
        // Details of newly added partitions
    }

    export interface DecreasePartitionsRequest {
        // Number of partitions to remove, strategy, etc.
    }

    export interface DecreasePartitionsResponse {
        // Result of partition reduction
    }

    export interface ChangePartitionLeaderRequest {
        // New leader node ID, etc.
    }

    export interface ChangePartitionLeaderResponse {
        // Result of leader change
    }

    export interface MovePartitionRequest {
        // Target node ID, etc.
    }

    export interface MovePartitionResponse {
        // Result of partition move
    }
}