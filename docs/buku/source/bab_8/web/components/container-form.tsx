"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { containerApi } from "@/lib/api/containers";

export function ContainerForm() {
    const [form, setForm] = useState({
        name: "",
        length_mm: 0,
        width_mm: 0,
        height_mm: 0,
        max_weight_kg: 0
    });
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            await containerApi.create(form);
            setForm({ name: "", length_mm: 0, width_mm: 0, height_mm: 0, max_weight_kg: 0 });
            
            // Dispatch event to refresh list
            window.dispatchEvent(new Event("container:refresh"));
        } catch (error) {
            alert("Failed to create container");
            console.error(error);
        } finally {
            setLoading(false);
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-lg">Add New Container</CardTitle>
                <CardDescription>Define container dimensions and weight limit.</CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="text-xs font-medium text-zinc-500 mb-1 block">Name</label>
                        <Input 
                            placeholder="e.g. 20ft Standard"
                            value={form.name} 
                            onChange={e => setForm({...form, name: e.target.value})}
                            required
                        />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="text-xs font-medium text-zinc-500 mb-1 block">Length (mm)</label>
                            <Input 
                                type="number"
                                value={form.length_mm || ""} 
                                onChange={e => setForm({...form, length_mm: Number(e.target.value)})}
                                required
                            />
                        </div>
                        <div>
                            <label className="text-xs font-medium text-zinc-500 mb-1 block">Width (mm)</label>
                            <Input 
                                type="number"
                                value={form.width_mm || ""} 
                                onChange={e => setForm({...form, width_mm: Number(e.target.value)})}
                                required
                            />
                        </div>
                        <div>
                            <label className="text-xs font-medium text-zinc-500 mb-1 block">Height (mm)</label>
                            <Input 
                                type="number"
                                value={form.height_mm || ""} 
                                onChange={e => setForm({...form, height_mm: Number(e.target.value)})}
                                required
                            />
                        </div>
                        <div>
                            <label className="text-xs font-medium text-zinc-500 mb-1 block">Max Weight (kg)</label>
                            <Input 
                                type="number"
                                value={form.max_weight_kg || ""} 
                                onChange={e => setForm({...form, max_weight_kg: Number(e.target.value)})}
                                required
                            />
                        </div>
                    </div>
                    <Button type="submit" disabled={loading} className="w-full">
                        {loading ? "Creating..." : "Create Container"}
                    </Button>
                </form>
            </CardContent>
        </Card>
    )
}
