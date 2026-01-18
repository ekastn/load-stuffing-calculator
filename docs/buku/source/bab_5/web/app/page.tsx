"use client";

import { useState, useMemo } from "react";
import dynamic from "next/dynamic";
import { Trash2 } from "lucide-react";
import type { StuffingPlanData } from "@/lib/StuffingVisualizer";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

const StuffingViewer = dynamic(
    () => import("@/components/stuffing-viewer").then((mod) => mod.StuffingViewer),
    { ssr: false, loading: () => <div className="flex items-center justify-center h-full">Loading 3D...</div> }
);

interface ItemInput {
    label: string;
    quantity: number;
    length_mm: number;
    width_mm: number;
    height_mm: number;
    weight_kg: number;
    color_hex: string;
}

// Generate random color
function randomColor(): string {
    const colors = ["#3b82f6", "#22c55e", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899", "#06b6d4"];
    return colors[Math.floor(Math.random() * colors.length)];
}

// Convert items to StuffingPlanData
function buildPlanData(items: ItemInput[], container: { length_mm: number; width_mm: number; height_mm: number; max_weight_kg: number }): StuffingPlanData {
    const allPlacements: StuffingPlanData["placements"] = [];
    const allItems: StuffingPlanData["items"] = [];
    
    let stepNumber = 0;
    let currentX = 0;
    let currentY = 0;
    let currentZ = 0;
    let rowMaxWidth = 0;
    let layerMaxHeight = 0;

    items.forEach((item, itemIndex) => {
        const itemId = `item-${itemIndex}`;
        allItems.push({
            item_id: itemId,
            label: item.label,
            length_mm: item.length_mm,
            width_mm: item.width_mm,
            height_mm: item.height_mm,
            weight_kg: item.weight_kg,
            quantity: item.quantity,
            color_hex: item.color_hex,
        });

        for (let q = 0; q < item.quantity; q++) {
            // Simple left-to-right, front-to-back, bottom-to-top packing
            if (currentX + item.length_mm > container.length_mm) {
                currentX = 0;
                currentY += rowMaxWidth;
                rowMaxWidth = 0;
            }
            if (currentY + item.width_mm > container.width_mm) {
                currentY = 0;
                currentZ += layerMaxHeight;
                layerMaxHeight = 0;
            }

            stepNumber++;
            allPlacements.push({
                placement_id: `p-${stepNumber}`,
                item_id: itemId,
                pos_x: currentX,
                pos_y: currentY,
                pos_z: currentZ,
                rotation: 0,
                step_number: stepNumber,
            });

            currentX += item.length_mm;
            rowMaxWidth = Math.max(rowMaxWidth, item.width_mm);
            layerMaxHeight = Math.max(layerMaxHeight, item.height_mm);
        }
    });

    return {
        plan_id: "demo-plan",
        plan_code: "DEMO-001",
        container: {
            name: "Custom Container",
            length_mm: container.length_mm,
            width_mm: container.width_mm,
            height_mm: container.height_mm,
            max_weight_kg: container.max_weight_kg,
        },
        items: allItems,
        placements: allPlacements,
        stats: {
            total_items: allPlacements.length,
            fitted_count: allPlacements.length,
            unfitted_count: 0,
            volume_utilization_pct: 0,
        },
    };
}

export default function HomePage() {
    const [container, setContainer] = useState({
        length_mm: 12000,
        width_mm: 2400,
        height_mm: 2600,
        max_weight_kg: 28000,
    });

    const [items, setItems] = useState<ItemInput[]>([]);
    const [form, setForm] = useState<ItemInput>({
        label: "",
        quantity: 1,
        length_mm: 0,
        width_mm: 0,
        height_mm: 0,
        weight_kg: 0,
        color_hex: randomColor(),
    });

    const [planData, setPlanData] = useState<StuffingPlanData | null>(null);

    const totalVolume = useMemo(
        () => items.reduce((sum, item) => sum + (item.length_mm * item.width_mm * item.height_mm * item.quantity) / 1_000_000_000, 0),
        [items]
    );

    const totalWeight = useMemo(
        () => items.reduce((sum, item) => sum + item.weight_kg * item.quantity, 0),
        [items]
    );

    const canAdd = form.label.trim() && form.quantity > 0 && form.length_mm > 0 && form.width_mm > 0 && form.height_mm > 0;

    const handleAdd = () => {
        if (!canAdd) return;
        setItems([...items, { ...form }]);
        setForm({ ...form, label: "", quantity: 1, length_mm: 0, width_mm: 0, height_mm: 0, weight_kg: 0, color_hex: randomColor() });
    };

    const handleRemove = (index: number) => {
        setItems(items.filter((_, i) => i !== index));
    };

    const handleCalculate = () => {
        if (items.length === 0) return;
        const data = buildPlanData(items, container);
        setPlanData(data);
    };

    return (
        <main className="min-h-screen bg-zinc-50">
            {/* Header */}
            <header className="border-b bg-white">
                <div className="container mx-auto px-4 py-4">
                    <h1 className="text-2xl font-bold">Load & Stuffing Calculator</h1>
                    <p className="text-sm text-zinc-600">3D Container Packing Visualization</p>
                </div>
            </header>

            <div className="container mx-auto px-4 py-6">
                <div className="grid gap-6 lg:grid-cols-12">
                    {/* Left Panel: Inputs */}
                    <div className="lg:col-span-5 space-y-6">
                        {/* Container Settings */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-base">Container</CardTitle>
                                <CardDescription>Enter container dimensions (mm)</CardDescription>
                            </CardHeader>
                            <CardContent className="grid grid-cols-2 gap-3">
                                <div>
                                    <label className="text-xs font-medium text-zinc-500">Length</label>
                                    <Input
                                        type="number"
                                        value={container.length_mm}
                                        onChange={(e) => setContainer({ ...container, length_mm: Number(e.target.value) })}
                                    />
                                </div>
                                <div>
                                    <label className="text-xs font-medium text-zinc-500">Width</label>
                                    <Input
                                        type="number"
                                        value={container.width_mm}
                                        onChange={(e) => setContainer({ ...container, width_mm: Number(e.target.value) })}
                                    />
                                </div>
                                <div>
                                    <label className="text-xs font-medium text-zinc-500">Height</label>
                                    <Input
                                        type="number"
                                        value={container.height_mm}
                                        onChange={(e) => setContainer({ ...container, height_mm: Number(e.target.value) })}
                                    />
                                </div>
                                <div>
                                    <label className="text-xs font-medium text-zinc-500">Max Weight (kg)</label>
                                    <Input
                                        type="number"
                                        value={container.max_weight_kg}
                                        onChange={(e) => setContainer({ ...container, max_weight_kg: Number(e.target.value) })}
                                    />
                                </div>
                            </CardContent>
                        </Card>

                        {/* Add Item */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="text-base">Add Item</CardTitle>
                                <CardDescription>Enter item dimensions (mm) and weight (kg)</CardDescription>
                            </CardHeader>
                            <CardContent className="space-y-3">
                                <div className="grid grid-cols-2 gap-3">
                                    <div className="col-span-2">
                                        <label className="text-xs font-medium text-zinc-500">Label</label>
                                        <Input
                                            placeholder="e.g. Carton Box"
                                            value={form.label}
                                            onChange={(e) => setForm({ ...form, label: e.target.value })}
                                        />
                                    </div>
                                    <div>
                                        <label className="text-xs font-medium text-zinc-500">Quantity</label>
                                        <Input
                                            type="number"
                                            value={form.quantity}
                                            onChange={(e) => setForm({ ...form, quantity: Math.max(1, Number(e.target.value)) })}
                                        />
                                    </div>
                                    <div>
                                        <label className="text-xs font-medium text-zinc-500">Weight (kg)</label>
                                        <Input
                                            type="number"
                                            value={form.weight_kg || ""}
                                            onChange={(e) => setForm({ ...form, weight_kg: Number(e.target.value) })}
                                        />
                                    </div>
                                    <div>
                                        <label className="text-xs font-medium text-zinc-500">Length (mm)</label>
                                        <Input
                                            type="number"
                                            value={form.length_mm || ""}
                                            onChange={(e) => setForm({ ...form, length_mm: Number(e.target.value) })}
                                        />
                                    </div>
                                    <div>
                                        <label className="text-xs font-medium text-zinc-500">Width (mm)</label>
                                        <Input
                                            type="number"
                                            value={form.width_mm || ""}
                                            onChange={(e) => setForm({ ...form, width_mm: Number(e.target.value) })}
                                        />
                                    </div>
                                    <div>
                                        <label className="text-xs font-medium text-zinc-500">Height (mm)</label>
                                        <Input
                                            type="number"
                                            value={form.height_mm || ""}
                                            onChange={(e) => setForm({ ...form, height_mm: Number(e.target.value) })}
                                        />
                                    </div>
                                    <div>
                                        <label className="text-xs font-medium text-zinc-500">Color</label>
                                        <Input
                                            type="color"
                                            value={form.color_hex}
                                            onChange={(e) => setForm({ ...form, color_hex: e.target.value })}
                                            className="h-10 p-1"
                                        />
                                    </div>
                                </div>
                                <Button onClick={handleAdd} disabled={!canAdd} className="w-full">
                                    Add Item
                                </Button>
                            </CardContent>
                        </Card>

                        {/* Item List */}
                        <Card>
                            <CardHeader className="flex flex-row items-center justify-between">
                                <div>
                                    <CardTitle className="text-base">Items</CardTitle>
                                    <CardDescription>{items.length} items added</CardDescription>
                                </div>
                                <div className="flex gap-2">
                                    <Badge variant="outline">{totalVolume.toFixed(2)} m³</Badge>
                                    <Badge variant="outline">{totalWeight.toFixed(1)} kg</Badge>
                                </div>
                            </CardHeader>
                            <CardContent>
                                {items.length === 0 ? (
                                    <p className="text-sm text-zinc-500">No items added yet.</p>
                                ) : (
                                    <Table>
                                        <TableHeader>
                                            <TableRow>
                                                <TableHead>Item</TableHead>
                                                <TableHead>Qty</TableHead>
                                                <TableHead>Size</TableHead>
                                                <TableHead></TableHead>
                                            </TableRow>
                                        </TableHeader>
                                        <TableBody>
                                            {items.map((item, idx) => (
                                                <TableRow key={idx}>
                                                    <TableCell>
                                                        <div className="flex items-center gap-2">
                                                            <div className="w-3 h-3 rounded" style={{ backgroundColor: item.color_hex }} />
                                                            {item.label}
                                                        </div>
                                                    </TableCell>
                                                    <TableCell>{item.quantity}</TableCell>
                                                    <TableCell className="text-xs">{item.length_mm}×{item.width_mm}×{item.height_mm}</TableCell>
                                                    <TableCell>
                                                        <Button variant="ghost" size="sm" onClick={() => handleRemove(idx)}>
                                                            <Trash2 className="h-4 w-4" />
                                                        </Button>
                                                    </TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>
                                )}
                            </CardContent>
                        </Card>

                        {/* Calculate Button */}
                        <Button onClick={handleCalculate} disabled={items.length === 0} size="lg" className="w-full">
                            Calculate Packing
                        </Button>
                    </div>

                    {/* Right Panel: 3D Visualization */}
                    <div className="lg:col-span-7">
                        <Card className="h-[700px]">
                            <CardHeader>
                                <CardTitle className="text-base">3D Visualization</CardTitle>
                                <CardDescription>
                                    {planData ? `${planData.placements.length} placements` : "Add items and click Calculate"}
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="h-[600px]">
                                {planData ? (
                                    <div className="h-full w-full rounded-lg border bg-white overflow-hidden">
                                        <StuffingViewer data={planData} />
                                    </div>
                                ) : (
                                    <div className="h-full flex items-center justify-center border border-dashed rounded-lg bg-zinc-100">
                                        <p className="text-zinc-500">No calculation result yet</p>
                                    </div>
                                )}
                            </CardContent>
                        </Card>
                    </div>
                </div>
            </div>
        </main>
    );
}
