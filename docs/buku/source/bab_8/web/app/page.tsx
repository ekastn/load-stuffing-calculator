"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Box, Truck, Map, LayoutDashboard, Package } from "lucide-react";
import { dashboardApi, DashboardStats } from "@/lib/api/dashboard";
import { Alert, AlertDescription } from "@/components/ui/alert";

export default function DashboardPage() {
    const [stats, setStats] = useState<DashboardStats | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        const fetchStats = async () => {
            try {
                const data = await dashboardApi.getStats();
                setStats(data);
            } catch (err: any) {
                console.error("Failed to fetch dashboard stats", err);
                setError("Failed to load dashboard data. Please try again.");
            } finally {
                setLoading(false);
            }
        };

        fetchStats();
    }, []);

    if (loading) {
        return (
            <main className="container mx-auto px-4 py-8">
                <div className="flex items-center justify-center p-12">
                     <p className="text-zinc-500 animate-pulse">Loading dashboard...</p>
                </div>
            </main>
        );
    }

    return (
        <main className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold tracking-tight mb-8">Dashboard Overview</h1>
            
            {error && (
                <Alert variant="destructive" className="mb-6">
                    <AlertDescription>{error}</AlertDescription>
                </Alert>
            )}

            {/* Stats Overview */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Plans</CardTitle>
                        <Map className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats?.total_plans || 0}</div>
                        <p className="text-xs text-muted-foreground">Projects created</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Items Shipped</CardTitle>
                        <Package className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats?.total_items_shipped || 0}</div>
                        <p className="text-xs text-muted-foreground">Total items in plans</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Product Catalog</CardTitle>
                        <Box className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats?.total_products || 0}</div>
                        <p className="text-xs text-muted-foreground">Active SKUs</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Container Types</CardTitle>
                        <Truck className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats?.total_containers || 0}</div>
                        <p className="text-xs text-muted-foreground">Available sizes</p>
                    </CardContent>
                </Card>
            </div>

            {/* Quick Actions */}
            <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
            <div className="grid gap-6 md:grid-cols-3">
                <Link href="/containers">
                    <Card className="hover:bg-zinc-50 hover:shadow-md transition-all cursor-pointer h-full border-zinc-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium text-zinc-600">Containers</CardTitle>
                            <Truck className="h-4 w-4 text-zinc-500" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-lg font-bold">Manage Fleet</div>
                            <p className="text-xs text-zinc-500 mt-1">Configure container sizes</p>
                        </CardContent>
                    </Card>
                </Link>

                <Link href="/products">
                    <Card className="hover:bg-zinc-50 hover:shadow-md transition-all cursor-pointer h-full border-zinc-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium text-zinc-600">Products</CardTitle>
                            <Box className="h-4 w-4 text-zinc-500" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-lg font-bold">Update Catalog</div>
                            <p className="text-xs text-zinc-500 mt-1">Manage cargo items</p>
                        </CardContent>
                    </Card>
                </Link>

                <Link href="/plans">
                    <Card className="hover:bg-zinc-50 hover:shadow-md transition-all cursor-pointer h-full border-zinc-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium text-zinc-600">Planning</CardTitle>
                            <Map className="h-4 w-4 text-zinc-500" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-lg font-bold">New Plan</div>
                            <p className="text-xs text-zinc-500 mt-1">Start calculation</p>
                        </CardContent>
                    </Card>
                </Link>
            </div>
        </main>
    );
}
