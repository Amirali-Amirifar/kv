"use client"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Slider } from "@/components/ui/slider"

export function PartitionManager() {
    return (
        <Card>
            <CardHeader>
                <CardTitle>Partition Management</CardTitle>
                <CardDescription>Configure and rebalance partitions</CardDescription>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="replication-factor">Replication Factor</Label>
                        <Select defaultValue="3">
                            <SelectTrigger id="replication-factor">
                                <SelectValue placeholder="Select replication factor" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="1">1 (No redundancy)</SelectItem>
                                <SelectItem value="2">2 (Minimal redundancy)</SelectItem>
                                <SelectItem value="3">3 (Recommended)</SelectItem>
                                <SelectItem value="5">5 (High availability)</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="partition-count">Partition Count</Label>
                        <div className="flex items-center space-x-2">
                            <Input id="partition-count" type="number" defaultValue="48" />
                            <Button variant="outline" size="sm">
                                Update
                            </Button>
                        </div>
                        <p className="text-xs text-muted-foreground">Recommended: 4x the number of nodes</p>
                    </div>

                    <div className="space-y-2">
                        <div className="flex justify-between">
                            <Label>Partition Balance</Label>
                            <span className="text-xs text-muted-foreground">92% balanced</span>
                        </div>
                        <Slider defaultValue={[92]} max={100} step={1} />
                    </div>

                    <div className="pt-4">
                        <Button className="w-full">Rebalance Partitions</Button>
                        <p className="text-xs text-muted-foreground mt-2">
                            Rebalancing will redistribute partitions across all available nodes. This operation may impact
                            performance.
                        </p>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
