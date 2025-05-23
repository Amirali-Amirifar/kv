// Base API URL - can be configured based on environment
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || '/api';

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