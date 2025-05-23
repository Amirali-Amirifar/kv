// Base API URL - can be configured based on environment
import {config} from "@/lib/config"
import useSWR from "swr";
import {useEffect, useState} from "react";

const API_BASE_URL = config.API_BASE_URL;

/**
 * Check if the API is running
 */
export const checkApiStatus = async (): Promise<string> => {
    const response = await fetch(`${API_BASE_URL}/`);

    if (!response.ok) {
        throw new Error('Failed to check API status');
    }

    return response.text();
};

/**
 * Check the health of the KV system
 */
export const checkHealth = async (): Promise<ApiTypes.HealthCheckResponse> => {
    const response = await fetch(`${API_BASE_URL}/health`);

    if (!response.ok) {
        throw new Error('Failed to check health status');
    }

    return response.json();
};

/**
 * Add a new node to the cluster
 */
export const addNode = async (nodeData: ApiTypes.AddNodeRequest): Promise<ApiTypes.AddNodeResponse> => {
    const response = await fetch(`${API_BASE_URL}/admin/nodes`, {
        method: 'POST', headers: {
            'Content-Type': 'application/json',
        }, body: JSON.stringify(nodeData),
    });

    if (!response.ok) {
        throw new Error('Failed to add node');
    }

    return response.json();
};

/**
 * Remove a node from the cluster
 */
export const removeNode = async (nodeId: string): Promise<ApiTypes.RemoveNodeResponse> => {
    const response = await fetch(`${API_BASE_URL}/admin/nodes/${nodeId}`, {
        method: 'DELETE',
    });

    if (!response.ok) {
        throw new Error('Failed to remove node');
    }

    return response.json();
};

/**
 * Increase the number of partitions
 */
export const increasePartitions = async (data: ApiTypes.IncreasePartitionsRequest): Promise<ApiTypes.IncreasePartitionsResponse> => {
    const response = await fetch(`${API_BASE_URL}/admin/partitions/increase`, {
        method: 'POST', headers: {
            'Content-Type': 'application/json',
        }, body: JSON.stringify(data),
    });

    if (!response.ok) {
        throw new Error('Failed to increase partitions');
    }

    return response.json();
};

/**
 * Decrease the number of partitions
 */
export const decreasePartitions = async (data: ApiTypes.DecreasePartitionsRequest): Promise<ApiTypes.DecreasePartitionsResponse> => {
    const response = await fetch(`${API_BASE_URL}/admin/partitions/decrease`, {
        method: 'POST', headers: {
            'Content-Type': 'application/json',
        }, body: JSON.stringify(data),
    });

    if (!response.ok) {
        throw new Error('Failed to decrease partitions');
    }

    return response.json();
};

/**
 * Change the leader of a partition
 */
export const changePartitionLeader = async (partitionId: string, data: ApiTypes.ChangePartitionLeaderRequest): Promise<ApiTypes.ChangePartitionLeaderResponse> => {
    const response = await fetch(`${API_BASE_URL}/admin/partitions/${partitionId}/leader`, {
        method: 'POST', headers: {
            'Content-Type': 'application/json',
        }, body: JSON.stringify(data),
    });

    if (!response.ok) {
        throw new Error('Failed to change partition leader');
    }

    return response.json();
};

/**
 * Move a partition to another node
 */
export const movePartition = async (partitionId: string, data: ApiTypes.MovePartitionRequest): Promise<ApiTypes.MovePartitionResponse> => {
    const response = await fetch(`${API_BASE_URL}/admin/partitions/${partitionId}/move`, {
        method: 'POST', headers: {
            'Content-Type': 'application/json',
        }, body: JSON.stringify(data),
    });

    if (!response.ok) {
        throw new Error('Failed to move partition');
    }

    return response.json();
};

export interface Address {
    ip: string
    port: number
}

interface Node {
    address: Address
    id: number
    leader_id: number
    node_type: string
    shard_key: number
    status: string
}

type Cluster = {"shards": Record<string, Node[]> }

export const getCluster = async (): Promise<Cluster> => {
    const response = await fetch(`${API_BASE_URL}/admin/cluster`, {
        headers: {
            "Content-Type": 'application/json'
        }
    });
    if (!response.ok) {
        throw new Error("Failed to get data")

    }
    return response.json()
}

export const useGetCluster = () => {
    return useSWR("cluster", () => getCluster())
}

export const useGetClusterStats = () => {
    const { data, error, isLoading } = useGetCluster();
    const [stats, setStats] = useState(() => ({
        totalNodes: 0,
        activeNodes: 0,
        warningNodes: 0,
        failedNodes: 0,
        totalPartitions: 0,
        replicationFactor: 0,
        storageUsed: 0, // percentage
        memoryUsed: 0, // percentage
        networkUsage: 0, // percentage
    }));

    useEffect(() => {
        if (!data || error) return;

        const partitions = (data).shards;

        let totalNodes = 0;
        let activeNodes = 0;
        let warningNodes = 0;
        let failedNodes = 0;
        let totalPartitions = Object.values(partitions).length;
        console.log(totalPartitions)
        let replicationCounts: number[] = [];


        console.log(partitions)
        Object.values(partitions).forEach((nodes) => {
            totalNodes += nodes.length;
            replicationCounts.push(nodes.length);

            Object.values(nodes).forEach(node => {
                if (node.status === "active" ) activeNodes++;
                else if (node.status === "SYNCING") warningNodes++;
                else if (node.status === "failed") failedNodes++;
            });
        });

        const replicationFactor = replicationCounts.length
            ? Math.round(replicationCounts.reduce((a, b) => a + b, 0) / replicationCounts.length)
            : 0;

        // These metrics would typically come from real-time monitoring,
        // here we just mock them randomly for demonstration
        const storageUsed = Math.round(Math.random() * 100);
        const memoryUsed = Math.round(Math.random() * 100);
        const networkUsage = Math.round(Math.random() * 100);

        setStats({
            totalNodes,
            activeNodes,
            warningNodes,
            failedNodes,
            totalPartitions,
            replicationFactor,
            storageUsed,
            memoryUsed,
            networkUsage
        });
    }, [data, error]);

    return {
        stats,
        isLoading,
        data,
        error
    };
};