"use client";

import { useEffect, useState } from "react";
import { productApi } from "@/lib/api/products";
import type { Product } from "@/lib/api/types";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";

export function ProductList() {
    const [products, setProducts] = useState<Product[]>([]);
    const [loading, setLoading] = useState(true);

    const loadProducts = async () => {
        setLoading(true);
        try {
            const data = await productApi.list();
            setProducts(data);
        } catch (error) {
            console.error("Failed to load products", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        loadProducts();
        const handleRefresh = () => loadProducts();
        window.addEventListener("product:refresh", handleRefresh);
        return () => window.removeEventListener("product:refresh", handleRefresh);
    }, []);

    const handleDelete = async (id: string) => {
        if (!confirm("Are you sure you want to delete this product?")) return;
        try {
            await productApi.delete(id);
            loadProducts();
        } catch (error) {
            alert("Failed to delete product");
            console.error(error);
        }
    };

    if (loading && products.length === 0) {
        return <div className="p-4 text-center text-zinc-500">Loading products...</div>;
    }

    if (products.length === 0) {
        return <div className="p-4 text-center text-zinc-500 border rounded-lg border-dashed">No products found. Create one to get started.</div>;
    }

    return (
        <div className="border rounded-md">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>SKU</TableHead>
                        <TableHead>Label</TableHead>
                        <TableHead>Dimensions</TableHead>
                        <TableHead>Weight</TableHead>
                        <TableHead className="w-[100px]">Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {products.map((p) => (
                        <TableRow key={p.id}>
                            <TableCell className="font-medium bg-zinc-50/50 text-xs font-mono">{p.sku}</TableCell>
                            <TableCell>{p.label}</TableCell>
                            <TableCell>
                                {p.length_mm} x {p.width_mm} x {p.height_mm} mm
                            </TableCell>
                            <TableCell>{p.weight_kg} kg</TableCell>
                            <TableCell>
                                <Button 
                                    size="sm" 
                                    variant="destructive" 
                                    onClick={() => handleDelete(p.id)}
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
