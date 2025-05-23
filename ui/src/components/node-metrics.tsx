"use client"

import { useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Progress } from "@/components/ui/progress"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

// Mock data - in a real app, this would come from an API
const nodeDetails = {
    "node-1": {
        name: "db-node-01",
        cpu: { current: 42, history: generateRandomData(24, 10, 80) },
        memory: { current: 68, history: generateRandomData(24, 40, 90) },
        disk: { current: 72, history: generateRandomData(24, 60, 85) },
        network: { current: 35, history: generateRandomData(24, 5, 60) },
        operations: { current: 1250, history: generateRandomData(24, 800, 2000) },
    },
    "node-2": {
        name: "db-node-02",
        cpu: { current: 38, history: generateRandomData(24, 10, 80) },
        memory: { current: 62, history: generateRandomData(24, 40, 90) },
        disk: { current: 68, history: generateRandomData(24, 60, 85) },
        network: { current: 42, history: generateRandomData(24, 5, 60) },
        operations: { current: 1180, history: generateRandomData(24, 800, 2000) },
    },
    "node-3": {
        name: "db-node-03",
        cpu: { current: 45, history: generateRandomData(24, 10, 80) },
        memory: { current: 70, history: generateRandomData(24, 40, 90) },
        disk: { current: 75, history: generateRandomData(24, 60, 85) },
        network: { current: 38, history: generateRandomData(24, 5, 60) },
        operations: { current: 1320, history: generateRandomData(24, 800, 2000) },
    },
    "node-4": {
        name: "db-node-04",
        cpu: { current: 78, history: generateRandomData(24, 50, 95) },
        memory: { current: 85, history: generateRandomData(24, 60, 95) },
        disk: { current: 72, history: generateRandomData(24, 60, 85) },
        network: { current: 62, history: generateRandomData(24, 30, 80) },
        operations: { current: 1850, history: generateRandomData(24, 1200, 2200) },
    },
    "node-5": {
        name: "db-node-05",
        cpu: { current: 0, history: Array(24).fill(0) },
        memory: { current: 0, history: Array(24).fill(0) },
        disk: { current: 72, history: generateRandomData(24, 60, 85) },
        network: { current: 0, history: Array(24).fill(0) },
        operations: { current: 0, history: Array(24).fill(0) },
    },
}

function generateRandomData(length: number, min: number, max: number) {
    return Array.from({ length }, () => Math.floor(Math.random() * (max - min + 1)) + min)
}

interface NodeMetricsProps {
    nodeId: string
}

export function NodeMetrics({ nodeId }: NodeMetricsProps) {
    const [timeRange, setTimeRange] = useState("24h")
    const node = nodeDetails[nodeId as keyof typeof nodeDetails]

    if (!node) {
        return (
            <Card>
                <CardContent className="pt-6">
                    <p>Node not found</p>
                </CardContent>
            </Card>
        )
    }

    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <div>
                    <CardTitle>Node Metrics: {node.name}</CardTitle>
                    <CardDescription>Performance and resource utilization</CardDescription>
                </div>
                <Select defaultValue={timeRange} onValueChange={setTimeRange}>
                    <SelectTrigger className="w-[120px]">
                        <SelectValue placeholder="Time range" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value="1h">Last hour</SelectItem>
                        <SelectItem value="6h">Last 6 hours</SelectItem>
                        <SelectItem value="24h">Last 24 hours</SelectItem>
                        <SelectItem value="7d">Last 7 days</SelectItem>
                    </SelectContent>
                </Select>
            </CardHeader>
            <CardContent>
                <Tabs defaultValue="resources">
                    <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value="resources">Resources</TabsTrigger>
                        <TabsTrigger value="operations">Operations</TabsTrigger>
                    </TabsList>
                    <TabsContent value="resources" className="space-y-4 pt-4">
                        <div className="space-y-2">
                            <div className="flex justify-between">
                                <div className="text-sm font-medium">CPU Usage</div>
                                <div className="text-sm text-muted-foreground">{node.cpu.current}%</div>
                            </div>
                            <Progress value={node.cpu.current} className="h-2" />
                            <div className="h-24 bg-muted/50 rounded-md flex items-end px-1">
                                {node.cpu.history.map((value, i) => (
                                    <div key={i} className="w-full bg-primary mx-[1px]" style={{ height: `${value}%` }} />
                                ))}
                            </div>
                        </div>

                        <div className="space-y-2">
                            <div className="flex justify-between">
                                <div className="text-sm font-medium">Memory Usage</div>
                                <div className="text-sm text-muted-foreground">{node.memory.current}%</div>
                            </div>
                            <Progress value={node.memory.current} className="h-2" />
                            <div className="h-24 bg-muted/50 rounded-md flex items-end px-1">
                                {node.memory.history.map((value, i) => (
                                    <div key={i} className="w-full bg-primary mx-[1px]" style={{ height: `${value}%` }} />
                                ))}
                            </div>
                        </div>

                        <div className="space-y-2">
                            <div className="flex justify-between">
                                <div className="text-sm font-medium">Disk Usage</div>
                                <div className="text-sm text-muted-foreground">{node.disk.current}%</div>
                            </div>
                            <Progress value={node.disk.current} className="h-2" />
                            <div className="h-24 bg-muted/50 rounded-md flex items-end px-1">
                                {node.disk.history.map((value, i) => (
                                    <div key={i} className="w-full bg-primary mx-[1px]" style={{ height: `${value}%` }} />
                                ))}
                            </div>
                        </div>

                        <div className="space-y-2">
                            <div className="flex justify-between">
                                <div className="text-sm font-medium">Network Usage</div>
                                <div className="text-sm text-muted-foreground">{node.network.current}%</div>
                            </div>
                            <Progress value={node.network.current} className="h-2" />
                            <div className="h-24 bg-muted/50 rounded-md flex items-end px-1">
                                {node.network.history.map((value, i) => (
                                    <div key={i} className="w-full bg-primary mx-[1px]" style={{ height: `${value}%` }} />
                                ))}
                            </div>
                        </div>
                    </TabsContent>

                    <TabsContent value="operations" className="pt-4">
                        <div className="space-y-4">
                            <div className="grid grid-cols-2 gap-4">
                                <Card>
                                    <CardHeader className="pb-2">
                                        <CardTitle className="text-sm">Operations/sec</CardTitle>
                                    </CardHeader>
                                    <CardContent>
                                        <div className="text-2xl font-bold">{node.operations.current}</div>
                                    </CardContent>
                                </Card>

                                <Card>
                                    <CardHeader className="pb-2">
                                        <CardTitle className="text-sm">Latency (avg)</CardTitle>
                                    </CardHeader>
                                    <CardContent>
                                        <div className="text-2xl font-bold">12ms</div>
                                    </CardContent>
                                </Card>

                                <Card>
                                    <CardHeader className="pb-2">
                                        <CardTitle className="text-sm">Read Ops</CardTitle>
                                    </CardHeader>
                                    <CardContent>
                                        <div className="text-2xl font-bold">{Math.floor(node.operations.current * 0.8)}</div>
                                    </CardContent>
                                </Card>

                                <Card>
                                    <CardHeader className="pb-2">
                                        <CardTitle className="text-sm">Write Ops</CardTitle>
                                    </CardHeader>
                                    <CardContent>
                                        <div className="text-2xl font-bold">{Math.floor(node.operations.current * 0.2)}</div>
                                    </CardContent>
                                </Card>
                            </div>

                            <div className="h-48 bg-muted/50 rounded-md flex items-end px-1">
                                {node.operations.history.map((value, i) => (
                                    <div key={i} className="w-full bg-primary mx-[1px]" style={{ height: `${(value / 2200) * 100}%` }} />
                                ))}
                            </div>
                        </div>
                    </TabsContent>
                </Tabs>
            </CardContent>
        </Card>
    )
}
