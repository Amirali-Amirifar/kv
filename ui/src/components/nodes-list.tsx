"use client"

import { useState} from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { MoreHorizontal, CheckCircle2, AlertTriangle, XCircle, PowerOff, RefreshCw, Trash2 } from "lucide-react"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import {useGetCluster} from "@/lib/api/api";



interface NodesListProps {
    onNodeSelect: (nodeId: string) => void
}

export function NodesList({ onNodeSelect }: NodesListProps) {
    const [nodeToRemove, setNodeToRemove] = useState<string | null>(null)

    const {data} = useGetCluster()

    // Transform cluster data into flat array of nodes with additional UI properties
    const nodes = data ? Object.values(data.shards).flat().map(node => ({
        ...node,
        name: `node-${node.id}`,
        ip: `${node.address.ip}:${node.address.port}`,
        region: "us-east", // Default region since not in API data
        partitions: node.shard_key || 0,
        cpu: 32, // Default values since not in API data
        memory: "64GB",
        storage: "2TB",
        uptime: node.status === "active" ? "Unknown" : "0d 0h", // Default uptime
    })) : []

    const getStatusIcon = (status: string) => {
        switch (status.toLowerCase()) {
            case "active":
            case "running":
                return <CheckCircle2 className="h-4 w-4 text-green-500" />
            case "warning":
            case "degraded":
                return <AlertTriangle className="h-4 w-4 text-yellow-500" />
            case "failed":
            case "down":
            case "stopped":
                return <XCircle className="h-4 w-4 text-red-500" />
            default:
                return <AlertTriangle className="h-4 w-4 text-gray-500" />
        }
    }

    const getStatusBadge = (status: string) => {
        const normalizedStatus = status.toLowerCase()
        switch (normalizedStatus) {
            case "active":
            case "running":
                return (
                    <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
                        Active
                    </Badge>
                )
            case "warning":
            case "degraded":
                return (
                    <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
                        Warning
                    </Badge>
                )
            case "failed":
            case "down":
            case "stopped":
                return (
                    <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200">
                        Failed
                    </Badge>
                )
            default:
                return (
                    <Badge variant="outline" className="bg-gray-50 text-gray-700 border-gray-200">
                        {status}
                    </Badge>
                )
        }
    }

    const handleRemoveNode = (nodeId: string) => {
        // In a real app, this would call an API to remove the node
        console.log(`Removing node: ${nodeId}`)
        setNodeToRemove(null)
    }

    const getNodeTypeDisplay = (nodeType: string) => {
        return nodeType.charAt(0).toUpperCase() + nodeType.slice(1)
    }

    return (
        <>
            <Card>
                <CardHeader>
                    <CardTitle>Database Nodes</CardTitle>
                    <CardDescription>
                        Manage your distributed database nodes ({nodes.length} nodes total)
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Status</TableHead>
                                <TableHead>Name</TableHead>
                                <TableHead>Address</TableHead>
                                <TableHead>Type</TableHead>
                                <TableHead>Shard Key</TableHead>
                                <TableHead>Leader ID</TableHead>
                                <TableHead className="text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {nodes.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={7} className="text-center text-muted-foreground">
                                        {data ? "No nodes found" : "Loading nodes..."}
                                    </TableCell>
                                </TableRow>
                            ) : (
                                nodes.map((node) => (
                                    <TableRow
                                        key={node.id}
                                        onClick={() => onNodeSelect(node.id.toString())}
                                        className="cursor-pointer"
                                    >
                                        <TableCell>
                                            <div className="flex items-center">
                                                {getStatusIcon(node.status)}
                                                <span className="ml-2 hidden md:inline-block">
                                                    {getStatusBadge(node.status)}
                                                </span>
                                            </div>
                                        </TableCell>
                                        <TableCell className="font-medium">{node.name}</TableCell>
                                        <TableCell>{node.ip}</TableCell>
                                        <TableCell>{getNodeTypeDisplay(node.node_type)}</TableCell>
                                        <TableCell>{node.shard_key}</TableCell>
                                        <TableCell>{node.leader_id}</TableCell>
                                        <TableCell className="text-right">
                                            <DropdownMenu>
                                                <DropdownMenuTrigger asChild>
                                                    <Button variant="ghost" className="h-8 w-8 p-0">
                                                        <span className="sr-only">Open menu</span>
                                                        <MoreHorizontal className="h-4 w-4" />
                                                    </Button>
                                                </DropdownMenuTrigger>
                                                <DropdownMenuContent align="end">
                                                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                                    <DropdownMenuItem onClick={() => onNodeSelect(node.id.toString())}>
                                                        View details
                                                    </DropdownMenuItem>
                                                    <DropdownMenuSeparator />
                                                    <DropdownMenuItem>
                                                        <RefreshCw className="mr-2 h-4 w-4" />
                                                        Restart node
                                                    </DropdownMenuItem>
                                                    <DropdownMenuItem>
                                                        <PowerOff className="mr-2 h-4 w-4" />
                                                        {node.status.toLowerCase() === "failed" || node.status.toLowerCase() === "stopped" ? "Start node" : "Stop node"}
                                                    </DropdownMenuItem>
                                                    <DropdownMenuSeparator />
                                                    <DropdownMenuItem
                                                        className="text-red-600"
                                                        onClick={() => setNodeToRemove(node.id.toString())}
                                                    >
                                                        <Trash2 className="mr-2 h-4 w-4" />
                                                        Remove node
                                                    </DropdownMenuItem>
                                                </DropdownMenuContent>
                                            </DropdownMenu>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>

            <AlertDialog open={!!nodeToRemove} onOpenChange={() => setNodeToRemove(null)}>
                <AlertDialogContent>
                    <AlertDialogHeader>
                        <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                        <AlertDialogDescription>
                            This will remove the node from the cluster. All data on this node will be redistributed according to the
                            replication policy. This action cannot be undone.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                            className="bg-red-600 hover:bg-red-700"
                            onClick={() => nodeToRemove && handleRemoveNode(nodeToRemove)}
                        >
                            Remove Node
                        </AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    )
}