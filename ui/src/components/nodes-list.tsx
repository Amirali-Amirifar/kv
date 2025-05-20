"use client"

import { useState } from "react"
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

// Mock data - in a real app, this would come from an API
const nodes = [
    {
        id: "node-1",
        name: "db-node-01",
        status: "active",
        ip: "10.0.1.101",
        region: "us-east",
        partitions: 12,
        cpu: 32,
        memory: "64GB",
        storage: "2TB",
        uptime: "45d 12h",
    },
    {
        id: "node-2",
        name: "db-node-02",
        status: "active",
        ip: "10.0.1.102",
        region: "us-east",
        partitions: 12,
        cpu: 32,
        memory: "64GB",
        storage: "2TB",
        uptime: "45d 10h",
    },
    {
        id: "node-3",
        name: "db-node-03",
        status: "active",
        ip: "10.0.1.103",
        region: "us-east",
        partitions: 12,
        cpu: 32,
        memory: "64GB",
        storage: "2TB",
        uptime: "44d 22h",
    },
    {
        id: "node-4",
        name: "db-node-04",
        status: "warning",
        ip: "10.0.1.104",
        region: "us-east",
        partitions: 12,
        cpu: 32,
        memory: "64GB",
        storage: "2TB",
        uptime: "15d 6h",
    },
    {
        id: "node-5",
        name: "db-node-05",
        status: "failed",
        ip: "10.0.1.105",
        region: "us-east",
        partitions: 0,
        cpu: 32,
        memory: "64GB",
        storage: "2TB",
        uptime: "0d 0h",
    },
]

interface NodesListProps {
    onNodeSelect: (nodeId: string) => void
}

export function NodesList({ onNodeSelect }: NodesListProps) {
    const [nodeToRemove, setNodeToRemove] = useState<string | null>(null)

    const getStatusIcon = (status: string) => {
        switch (status) {
            case "active":
                return <CheckCircle2 className="h-4 w-4 text-green-500" />
            case "warning":
                return <AlertTriangle className="h-4 w-4 text-yellow-500" />
            case "failed":
                return <XCircle className="h-4 w-4 text-red-500" />
            default:
                return null
        }
    }

    const getStatusBadge = (status: string) => {
        switch (status) {
            case "active":
                return (
                    <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
                        Active
                    </Badge>
                )
            case "warning":
                return (
                    <Badge variant="outline" className="bg-yellow-50 text-yellow-700 border-yellow-200">
                        Warning
                    </Badge>
                )
            case "failed":
                return (
                    <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200">
                        Failed
                    </Badge>
                )
            default:
                return null
        }
    }

    const handleRemoveNode = (nodeId: string) => {
        // In a real app, this would call an API to remove the node
        console.log(`Removing node: ${nodeId}`)
        setNodeToRemove(null)
    }

    return (
        <>
            <Card>
                <CardHeader>
                    <CardTitle>Database Nodes</CardTitle>
                    <CardDescription>Manage your distributed database nodes</CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Status</TableHead>
                                <TableHead>Name</TableHead>
                                <TableHead>IP Address</TableHead>
                                <TableHead>Region</TableHead>
                                <TableHead>Partitions</TableHead>
                                <TableHead>Uptime</TableHead>
                                <TableHead className="text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {nodes.map((node) => (
                                <TableRow key={node.id} onClick={() => onNodeSelect(node.id)} className="cursor-pointer">
                                    <TableCell>
                                        <div className="flex items-center">
                                            {getStatusIcon(node.status)}
                                            <span className="ml-2 hidden md:inline-block">{getStatusBadge(node.status)}</span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="font-medium">{node.name}</TableCell>
                                    <TableCell>{node.ip}</TableCell>
                                    <TableCell>{node.region}</TableCell>
                                    <TableCell>{node.partitions}</TableCell>
                                    <TableCell>{node.uptime}</TableCell>
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
                                                <DropdownMenuItem onClick={() => onNodeSelect(node.id)}>View details</DropdownMenuItem>
                                                <DropdownMenuSeparator />
                                                <DropdownMenuItem>
                                                    <RefreshCw className="mr-2 h-4 w-4" />
                                                    Restart node
                                                </DropdownMenuItem>
                                                <DropdownMenuItem>
                                                    <PowerOff className="mr-2 h-4 w-4" />
                                                    {node.status === "failed" ? "Start node" : "Stop node"}
                                                </DropdownMenuItem>
                                                <DropdownMenuSeparator />
                                                <DropdownMenuItem className="text-red-600" onClick={() => setNodeToRemove(node.id)}>
                                                    <Trash2 className="mr-2 h-4 w-4" />
                                                    Remove node
                                                </DropdownMenuItem>
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    </TableCell>
                                </TableRow>
                            ))}
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
