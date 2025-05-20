import { Metadata } from "next"
import DashboardView from "@/components/dashboard-view"

export const metadata: Metadata = {
    title: "Database Node Management Dashboard",
    description: "Monitor and manage your distributed database nodes",
}

export default function DashboardPage() {
    return <DashboardView />
}
