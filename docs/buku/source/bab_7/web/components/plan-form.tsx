"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { planApi } from "@/lib/api/plans";
import { containerApi } from "@/lib/api/containers";
import { useRouter } from "next/navigation";
import type { Container } from "@/lib/api/types";

export function PlanForm() {
    const [containers, setContainers] = useState<Container[]>([]);
    const [selectedContainer, setSelectedContainer] = useState<string>("");
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    useEffect(() => {
        containerApi.list().then(setContainers).catch(console.error);
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedContainer) return;

        setLoading(true);
        try {
            const res = await planApi.create({ container_id: selectedContainer });
            // Redirect to detail page to add items
            router.push(`/plans/${res.id}`);
        } catch (error) {
            alert("Failed to create plan");
            console.error(error);
            setLoading(false); // Only stop loading on error, on success we redirect
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-lg">Create New Plan</CardTitle>
                <CardDescription>Start by selecting a container.</CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="text-xs font-medium text-zinc-500 mb-1 block">Container</label>
                        <Select onValueChange={setSelectedContainer} value={selectedContainer}>
                            <SelectTrigger>
                                <SelectValue placeholder="Select a container" />
                            </SelectTrigger>
                            <SelectContent>
                                {containers.map((c) => (
                                    <SelectItem key={c.id} value={c.id}>
                                        {c.name} ({c.length_mm}x{c.width_mm}x{c.height_mm} mm)
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                    <Button type="submit" disabled={loading || !selectedContainer} className="w-full">
                        {loading ? "Creating..." : "Start Planning"}
                    </Button>
                </form>
            </CardContent>
        </Card>
    )
}
