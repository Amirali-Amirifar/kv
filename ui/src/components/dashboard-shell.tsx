"use client"

import type React from "react"
import Link from "next/link"
import { Database, LayoutDashboard, Settings, HardDrive, Network, AlertCircle, BarChart3 } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuItem,
    SidebarMenuButton,
    SidebarProvider,
    SidebarTrigger,
} from "@/components/ui/sidebar"

interface DashboardShellProps {
    children: React.ReactNode
}

export function DashboardShell({ children }: DashboardShellProps) {
    return (
        <SidebarProvider>
            <div className="flex min-h-screen">
                <DashboardSidebar />
                <main className="flex-1 bg-muted/40">
                    <div className="flex items-center h-16 px-4 border-b bg-background md:px-6">
                        <SidebarTrigger />
                        <div className="ml-auto flex items-center space-x-4">
                            <Button variant="outline" size="sm">
                                <AlertCircle className="mr-2 h-4 w-4" />
                                System Alerts
                            </Button>
                        </div>
                    </div>
                    <div className="p-4 md:p-8">{children}</div>
                </main>
            </div>
        </SidebarProvider>
    )
}

function DashboardSidebar() {
    return (
        <Sidebar>
            <SidebarHeader>
                <div className="flex items-center gap-2 px-2">
                    <Database className="h-6 w-6" />
                    <span className="font-bold">DB Manager</span>
                </div>
            </SidebarHeader>
            <SidebarContent>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton asChild isActive>
                            <Link href="/dashboard">
                                <LayoutDashboard className="h-4 w-4" />
                                <span>Dashboard</span>
                            </Link>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarContent>
            <SidebarFooter>
                <div className="px-3 py-2">
                    <div className="text-xs text-muted-foreground">Cluster ID: db-cluster-01</div>
                    <div className="text-xs text-muted-foreground">Version: 2.5.3</div>
                </div>
            </SidebarFooter>
        </Sidebar>
    )
}
