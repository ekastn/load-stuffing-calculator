"use client";

import { useEffect, useState } from "react";
import { planApi } from "@/lib/api/plans";
import type { Plan } from "@/lib/api/types";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Trash2, Box, ArrowRight } from "lucide-react";
import { PlanForm } from "@/components/plan-form";
import { useRouter } from "next/navigation";
import { Badge } from "@/components/ui/badge";

export default function PlansPage() {
    const [plans, setPlans] = useState<Plan[]>([]);
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    const loadPlans = async () => {
        setLoading(true);
        try {
            const data = await planApi.list();
            setPlans(data);
        } catch (error) {
            console.error("Failed to load plans", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadPlans();
    }, []);

    const handleDelete = async (e: React.MouseEvent, id: string) => {
        e.stopPropagation();
        if (!confirm("Are you sure you want to delete this plan?")) return;
        try {
            await planApi.delete(id);
            loadPlans();
        } catch (error) {
            alert("Failed to delete plan");
        }
    };

    return (
        <main className="container mx-auto px-4 py-8">
            <div className="flex flex-col items-center justify-center mb-8 text-center bg-white p-6 rounded-lg border shadow-sm">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Planning</h1>
                    <p className="text-zinc-500 mt-2">Create and manage stuffing plans.</p>
                </div>
            </div>

            <div className="grid gap-8 lg:grid-cols-12">
                <div className="lg:col-span-4">
                    <PlanForm />
                </div>
                <div className="lg:col-span-8">
                    <div className="border rounded-md">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>Status</TableHead>
                                    <TableHead>Container</TableHead>
                                    <TableHead>Plan ID</TableHead>
                                    <TableHead className="w-[100px]">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {plans.length === 0 && !loading && (
                                     <TableRow>
                                        <TableCell colSpan={4} className="text-center py-8 text-zinc-500">
                                            No plans created yet.
                                        </TableCell>
                                    </TableRow>
                                )}
                                {plans.map((p) => (
                                    <TableRow 
                                        key={p.id} 
                                        className="cursor-pointer hover:bg-zinc-50"
                                        onClick={() => router.push(`/plans/${p.id}`)}
                                    >
                                        <TableCell>
                                            <Badge variant={p.status === "completed" ? "default" : "secondary"}>
                                                {p.status}
                                            </Badge>
                                        </TableCell>
                                        <TableCell>{p.container_name || "Unknown Container"}</TableCell>
                                        <TableCell className="font-mono text-xs text-zinc-500">{p.id.substring(0, 8)}...</TableCell>
                                        <TableCell>
                                            <div className="flex gap-2">
                                                <Button 
                                                    size="sm" 
                                                    variant="ghost" 
                                                    onClick={(e) => handleDelete(e, p.id)}
                                                >
                                                    <Trash2 className="w-4 h-4 text-red-500" />
                                                </Button>
                                                <Button size="sm" variant="ghost">
                                                    <ArrowRight className="w-4 h-4" />
                                                </Button>
                                            </div>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </div>
                </div>
            </div>
        </main>
    );
}
