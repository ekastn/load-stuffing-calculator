"use client";

import Link from "next/link";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Box, Truck, Map, ArrowRight } from "lucide-react";

export default function DashboardPage() {
    return (
        <main className="container mx-auto px-4 py-8">
            <h1 className="text-3xl font-bold tracking-tight mb-8 text-center">Dashboard</h1>
            
            <div className="grid gap-6 md:grid-cols-3">
                <Link href="/containers">
                    <Card className="hover:bg-zinc-50 hover:shadow-md transition-all cursor-pointer h-full border-zinc-200">
                        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                            <CardTitle className="text-sm font-medium text-zinc-600">Containers</CardTitle>
                            <Truck className="h-4 w-4 text-zinc-500" />
                        </CardHeader>
                        <CardContent>
                            <div className="text-2xl font-bold">Manage</div>
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
                            <div className="text-2xl font-bold">Catalog</div>
                            <p className="text-xs text-zinc-500 mt-1">Manage cargo items and SKU</p>
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
                            <div className="text-2xl font-bold">Create Plan</div>
                            <p className="text-xs text-zinc-500 mt-1">Run stuffing calculations and visualize</p>
                        </CardContent>
                    </Card>
                </Link>
            </div>
        </main>
    );
}
