"use client"

import {useState} from "react"
import {DashboardShell} from "@/components/dashboard-shell"
import {ClusterOverview} from "@/components/cluster-overview"
import {NodesList} from "@/components/nodes-list"
import {PartitionManager} from "@/components/partition-manager"
import {NodeMetrics} from "@/components/node-metrics"
import {AddNodeDialog} from "@/components/add-node-dialog"
import {Button} from "@/components/ui/button"
import {PlusCircle} from "lucide-react"

export default function DashboardView() {
    const [showAddNodeDialog, setShowAddNodeDialog] = useState(false)
    const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null)

    return (<DashboardShell>
            <div className="flex items-center justify-between">
                <h1 className="text-2xl font-bold tracking-tight">Cluster
                    Dashboard</h1>
                <Button onClick={() => setShowAddNodeDialog(true)}>
                    <PlusCircle className="mr-2 h-4 w-4"/>
                    Add Node
                </Button>
            </div>

            <div className="grid gap-6 mt-6">
                <ClusterOverview/>

                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                    <div className="md:col-span-2">
                        <NodesList
                            onNodeSelect={(id) => setSelectedNodeId(id)}/>
                    </div>
                    <div>
                        <PartitionManager/>
                    </div>
                </div>

                {selectedNodeId && (<div className="mt-6">
                        <NodeMetrics nodeId={selectedNodeId}/>
                    </div>)}
            </div>

            <AddNodeDialog open={showAddNodeDialog}
                           onOpenChange={setShowAddNodeDialog}/>
        </DashboardShell>)
}
