"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { planApi } from "@/lib/api/plans";
import { containerApi } from "@/lib/api/containers";
import { productApi } from "@/lib/api/products";
import type { PlanDetail, Container, Product } from "@/lib/api/types";
import { transformToVisualizerData } from "@/lib/transforms";
import { StuffingViewer } from "@/components/stuffing-viewer";
import { PackingResult } from "@/components/packing-result";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Loader2, Box, Truck, Calculator, ArrowLeft } from "lucide-react";
import Link from "next/link";
import { ItemMesh } from "@/lib/StuffingVisualizer/components/item-mesh";
import { PDFExportButton } from "@/components/pdf-export-button";

export default function PlanDetailPage() {
    const params = useParams();
    const id = params.id as string;

    const [plan, setPlan] = useState<PlanDetail | null>(null);
    const [container, setContainer] = useState<Container | null>(null);
    const [products, setProducts] = useState<Product[]>([]);
    
    // Add Item Form State
    const [selectedProduct, setSelectedProduct] = useState<string>("");
    const [quantity, setQuantity] = useState<number>(1);
    
    // Loading States
    const [loadingPlan, setLoadingPlan] = useState(true);
    const [addingItem, setAddingItem] = useState(false);
    const [calculating, setCalculating] = useState(false);

    const loadData = async () => {
        setLoadingPlan(true);
        try {
            const planData = await planApi.get(id);
            setPlan(planData);

            // Load container info
            const containerData = await containerApi.get(planData.container_id);
            setContainer(containerData);

            // Load products for dropdown
            const productsData = await productApi.list();
            setProducts(productsData);
        } catch (error) {
            console.error("Failed to load plan details", error);
        } finally {
            setLoadingPlan(false);
        }
    };

    useEffect(() => {
        if (id) {
            loadData();
        }
    }, [id]);

    const handleAddItem = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedProduct || quantity <= 0) return;

        setAddingItem(true);
        try {
            await planApi.addItem(id, {
                product_id: selectedProduct,
                quantity: quantity
            });
            // Reload plan items
            const updatedPlan = await planApi.get(id);
            setPlan(updatedPlan);
            setQuantity(1);
        } catch (error) {
            alert("Failed to add item");
        } finally {
            setAddingItem(false);
        }
    };

    const handleCalculate = async () => {
        setCalculating(true);
        try {
            await planApi.calculate(id);
            // In a real app we might poll for status, here we wait for response
            // Reload plan to get results
            const updatedPlan = await planApi.get(id);
            setPlan(updatedPlan);
        } catch (error) {
            alert("Calculation failed");
        } finally {
            setCalculating(false);
        }
    };

    if (loadingPlan || !plan || !container) {
        return (
            <div className="flex h-screen items-center justify-center">
                <Loader2 className="w-8 h-8 animate-spin text-zinc-500" />
            </div>
        );
    }

    const visualizerData = transformToVisualizerData(plan, container);
    const hasResults = plan.placements && plan.placements.length > 0;

    return (
        <main className="container mx-auto px-4 py-8 max-w-7xl">
            <div className="mb-8">
                <Link href="/plans" className="text-sm text-zinc-500 hover:text-zinc-900 flex items-center gap-1 mb-4">
                    <ArrowLeft className="w-4 h-4" /> Back to Plans
                </Link>
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight flex items-center gap-3">
                            Plan: {plan.id.substring(0, 8)}
                            <Badge variant={plan.status === "completed" ? "default" : "secondary"}>
                                {plan.status}
                            </Badge>
                        </h1>
                        <p className="text-zinc-500 mt-2 flex items-center gap-2">
                            <Truck className="w-4 h-4" /> 
                            Container: <span className="font-medium text-zinc-900">{container.name}</span>
                            <span className="text-zinc-300">|</span>
                            {container.length_mm} x {container.width_mm} x {container.height_mm} mm
                        </p>
                    </div>
                    <div className="flex items-center gap-2">
                        <PDFExportButton plan={plan} container={container} />
                        {plan.items.length > 0 && (
                            <Button onClick={handleCalculate} disabled={calculating}>
                                {calculating ? (
                                    <>
                                        <Loader2 className="w-4 h-4 mr-2 animate-spin" /> Calculating...
                                    </>
                                ) : (
                                    <>
                                        <Calculator className="w-4 h-4 mr-2" /> Calculate Packing
                                    </>
                                )}
                            </Button>
                        )}
                    </div>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
                {/* Left Column: Configuration */}
                <div className="lg:col-span-4 space-y-6">
                    {/* Items List */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-lg">Cargo Items</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="mb-6">
                                <form onSubmit={handleAddItem} className="space-y-3">
                                    <div className="grid grid-cols-3 gap-3">
                                        <div className="col-span-2">
                                            <Select onValueChange={setSelectedProduct} value={selectedProduct}>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Product" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    {products.map(p => (
                                                        <SelectItem key={p.id} value={p.id}>{p.label}</SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                        </div>
                                        <div>
                                            <Input 
                                                type="number" 
                                                min="1" 
                                                value={quantity} 
                                                onChange={e => setQuantity(parseInt(e.target.value))} 
                                            />
                                        </div>
                                    </div>
                                    <Button size="sm" type="submit" className="w-full" disabled={addingItem || !selectedProduct}>
                                        {addingItem ? "Adding..." : "Add Item"}
                                    </Button>
                                </form>
                            </div>

                            <div className="space-y-2">
                                {plan.items.length === 0 && (
                                    <div className="text-center py-4 text-zinc-500 text-sm border border-dashed rounded">
                                        No items added yet.
                                    </div>
                                )}
                                {plan.items.map((item) => (
                                    <div key={item.id} className="flex items-center justify-between p-3 bg-zinc-50 rounded-lg border">
                                        <div className="flex items-center gap-3">
                                            <div className="w-2 h-8 bg-blue-500 rounded-full" />
                                            <div>
                                                <div className="font-medium text-sm">{item.label}</div>
                                                <div className="text-xs text-zinc-500">
                                                    {item.length_mm}x{item.width_mm}x{item.height_mm}
                                                </div>
                                            </div>
                                        </div>
                                        <div className="font-bold text-lg text-zinc-700">x{item.quantity}</div>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>

                    {/* Stats */}
                    {hasResults && (
                        <PackingResult stats={visualizerData.stats} />
                    )}
                </div>

                {/* Right Column: Visualization */}
                <div className="lg:col-span-8">
                    <Card className="h-[600px] overflow-hidden flex flex-col">
                         <CardHeader className="py-4 border-b">
                            <CardTitle className="text-lg flex items-center justify-between">
                                3D Visualization
                                {!hasResults && <span className="text-xs font-normal text-zinc-500">Calculate to view result</span>}
                            </CardTitle>
                        </CardHeader>
                        <div className="flex-1 bg-zinc-100 relative">
                            {hasResults ? (
                                <StuffingViewer data={visualizerData} />
                            ) : (
                                <div className="absolute inset-0 flex items-center justify-center text-zinc-400">
                                    <div className="text-center">
                                        <Box className="w-16 h-16 mx-auto mb-4 opacity-20" />
                                        <p>Add items and calculate<br/>to generate 3D packing plan</p>
                                    </div>
                                </div>
                            )}
                        </div>
                    </Card>
                </div>
            </div>
        </main>
    );
}
