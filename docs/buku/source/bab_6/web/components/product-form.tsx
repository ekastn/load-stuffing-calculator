"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { productApi } from "@/lib/api/products";

export function ProductForm() {
    const [form, setForm] = useState({
        label: "",
        sku: "",
        length_mm: 0,
        width_mm: 0,
        height_mm: 0,
        weight_kg: 0
    });
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            await productApi.create(form);
            setForm({ label: "", sku: "", length_mm: 0, width_mm: 0, height_mm: 0, weight_kg: 0 });
            window.dispatchEvent(new Event("product:refresh"));
        } catch (error) {
            alert("Failed to create product");
            console.error(error);
        } finally {
            setLoading(false);
        }
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle className="text-lg">Add New Product</CardTitle>
                <CardDescription>Register a new product SKU and dimensions.</CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="text-xs font-medium text-zinc-500 mb-1 block">Label</label>
                            <Input 
                                placeholder="e.g. Cardboard Box A"
                                value={form.label} 
                                onChange={e => setForm({...form, label: e.target.value})}
                                required
                            />
                        </div>
                        <div>
                             <label className="text-xs font-medium text-zinc-500 mb-1 block">SKU</label>
                            <Input 
                                placeholder="BOX-001"
                                value={form.sku} 
                                onChange={e => setForm({...form, sku: e.target.value})}
                                required
                            />
                        </div>
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
                            <label className="text-xs font-medium text-zinc-500 mb-1 block">Weight (kg)</label>
                            <Input 
                                type="number"
                                value={form.weight_kg || ""} 
                                onChange={e => setForm({...form, weight_kg: Number(e.target.value)})}
                                required
                            />
                        </div>
                    </div>
                    <Button type="submit" disabled={loading} className="w-full">
                        {loading ? "Creating..." : "Create Product"}
                    </Button>
                </form>
            </CardContent>
        </Card>
    )
}
