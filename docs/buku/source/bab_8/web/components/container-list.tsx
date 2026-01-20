"use client";

import { useEffect, useState } from "react";
import { containerApi } from "@/lib/api/containers";
import type { Container } from "@/lib/api/types";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";

export function ContainerList() {
    const [containers, setContainers] = useState<Container[]>([]);
    const [loading, setLoading] = useState(true);

    const loadContainers = async () => {
        setLoading(true);
        try {
            const data = await containerApi.list();
            setContainers(data);
        } catch (error) {
            console.error("Failed to load containers", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadContainers();
        
        // Listen for refresh events
        const handleRefresh = () => loadContainers();
        window.addEventListener("container:refresh", handleRefresh);
        return () => window.removeEventListener("container:refresh", handleRefresh);
    }, []);

    const handleDelete = async (id: string) => {
        if (!confirm("Are you sure you want to delete this container?")) return;
        try {
            await containerApi.delete(id);
            loadContainers();
        } catch (error) {
            alert("Failed to delete container");
            console.error(error);
        }
    };

    if (loading && containers.length === 0) {
        return <div className="p-4 text-center text-zinc-500">Loading containers...</div>;
    }

    if (containers.length === 0) {
        return <div className="p-4 text-center text-zinc-500 border rounded-lg border-dashed">No containers found. Create one to get started.</div>;
    }

    return (
        <div className="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>Name</TableHead>
                        <TableHead>Dimensions (Length x Width x Height)</TableHead>
                        <TableHead>Max Weight</TableHead>
                        <TableHead className="w-[100px]">Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {containers.map((c) => (
                        <TableRow key={c.id}>
                            <TableCell className="font-medium">{c.name}</TableCell>
                            <TableCell>
                                {c.length_mm} x {c.width_mm} x {c.height_mm} mm
                            </TableCell>
                            <TableCell>{c.max_weight_kg} kg</TableCell>
                            <TableCell>
                                <Button 
                                    size="sm" 
                                    variant="destructive" 
                                    onClick={() => handleDelete(c.id)}
                                >
                                    <Trash2 className="w-4 h-4" />
                                </Button>
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
}
