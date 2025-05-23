"use client";

import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card"
import {Progress} from "@/components/ui/progress"
import {
    AlertTriangle, CheckCircle2, Database, HardDrive, Network, XCircle
} from "lucide-react"
import {useGetCluster, useGetClusterStats} from "@/lib/api/api";

// Mock data - in a real app, this would come from an API


export function ClusterOverview() {
    const { stats: clusterStats} = useGetClusterStats()

    return (<div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
                <CardHeader
                    className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">Node
                        Status</CardTitle>
                    <HardDrive className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div className="text-2xl font-bold">
                        {clusterStats.activeNodes}/{clusterStats.totalNodes}
                    </div>
                    <div
                        className="flex items-center space-x-2 text-xs text-muted-foreground mt-2">
                        <div className="flex items-center">
                            <CheckCircle2
                                className="mr-1 h-3 w-3 text-green-500"/>
                            <span>{clusterStats.activeNodes} Active</span>
                        </div>
                        <div className="flex items-center">
                            <AlertTriangle
                                className="mr-1 h-3 w-3 text-yellow-500"/>
                            <span>{clusterStats.warningNodes} Syncing</span>
                        </div>
                        <div className="flex items-center">
                            <XCircle className="mr-1 h-3 w-3 text-red-500"/>
                            <span>{clusterStats.failedNodes} Failed</span>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader
                    className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle
                        className="text-sm font-medium">Partitions</CardTitle>
                    <Database className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div
                        className="text-2xl font-bold">{clusterStats.totalPartitions}</div>
                    <p className="text-xs text-muted-foreground mt-2">Replication
                        Factor: {clusterStats.replicationFactor}x</p>
                </CardContent>
            </Card>

            <Card>
                <CardHeader
                    className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle
                        className="text-sm font-medium">Storage</CardTitle>
                    <HardDrive className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div
                        className="text-2xl font-bold">{clusterStats.storageUsed}%
                    </div>
                    <Progress value={clusterStats.storageUsed}
                              className="h-2 mt-2"/>
                </CardContent>
            </Card>

            <Card>
                <CardHeader
                    className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle
                        className="text-sm font-medium">Network</CardTitle>
                    <Network className="h-4 w-4 text-muted-foreground"/>
                </CardHeader>
                <CardContent>
                    <div
                        className="text-2xl font-bold">{clusterStats.networkUsage}%
                    </div>
                    <Progress value={clusterStats.networkUsage}
                              className="h-2 mt-2"/>
                </CardContent>
            </Card>
        </div>)
}
