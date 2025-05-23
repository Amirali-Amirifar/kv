"use client"

import type React from "react"

import { useState } from "react"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Checkbox } from "@/components/ui/checkbox"

interface AddNodeDialogProps {
    open: boolean
    onOpenChange: (open: boolean) => void
}

export function AddNodeDialog({ open, onOpenChange }: AddNodeDialogProps) {
    const [isSubmitting, setIsSubmitting] = useState(false)

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        setIsSubmitting(true)

        // Simulate API call
        setTimeout(() => {
            setIsSubmitting(false)
            onOpenChange(false)
        }, 1500)
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[500px]">
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>Add New Node</DialogTitle>
                        <DialogDescription>Add a new node to your distributed database cluster.</DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="name" className="text-right">
                                Name
                            </Label>
                            <Input id="name" placeholder="db-node-06" className="col-span-3" required />
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="ip" className="text-right">
                                IP Address
                            </Label>
                            <Input id="ip" placeholder="10.0.1.106" className="col-span-3" required />
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="region" className="text-right">
                                Region
                            </Label>
                            <Select defaultValue="us-east">
                                <SelectTrigger id="region" className="col-span-3">
                                    <SelectValue placeholder="Select region" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="us-east">US East</SelectItem>
                                    <SelectItem value="us-west">US West</SelectItem>
                                    <SelectItem value="eu-central">EU Central</SelectItem>
                                    <SelectItem value="ap-south">Asia Pacific</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="cpu" className="text-right">
                                CPU Cores
                            </Label>
                            <Select defaultValue="32">
                                <SelectTrigger id="cpu" className="col-span-3">
                                    <SelectValue placeholder="Select CPU cores" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="8">8 cores</SelectItem>
                                    <SelectItem value="16">16 cores</SelectItem>
                                    <SelectItem value="32">32 cores</SelectItem>
                                    <SelectItem value="64">64 cores</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="memory" className="text-right">
                                Memory
                            </Label>
                            <Select defaultValue="64">
                                <SelectTrigger id="memory" className="col-span-3">
                                    <SelectValue placeholder="Select memory" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="16">16 GB</SelectItem>
                                    <SelectItem value="32">32 GB</SelectItem>
                                    <SelectItem value="64">64 GB</SelectItem>
                                    <SelectItem value="128">128 GB</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Label htmlFor="storage" className="text-right">
                                Storage
                            </Label>
                            <Select defaultValue="2">
                                <SelectTrigger id="storage" className="col-span-3">
                                    <SelectValue placeholder="Select storage" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="1">1 TB</SelectItem>
                                    <SelectItem value="2">2 TB</SelectItem>
                                    <SelectItem value="4">4 TB</SelectItem>
                                    <SelectItem value="8">8 TB</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <div className="text-right">Options</div>
                            <div className="col-span-3 space-y-2">
                                <div className="flex items-center space-x-2">
                                    <Checkbox id="auto-join" defaultChecked />
                                    <Label htmlFor="auto-join">Auto-join cluster</Label>
                                </div>
                                <div className="flex items-center space-x-2">
                                    <Checkbox id="rebalance" defaultChecked />
                                    <Label htmlFor="rebalance">Rebalance partitions</Label>
                                </div>
                            </div>
                        </div>
                    </div>
                    <DialogFooter>
                        <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                            Cancel
                        </Button>
                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? "Adding..." : "Add Node"}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
