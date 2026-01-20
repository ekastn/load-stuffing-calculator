import { ProductForm } from "@/components/product-form";
import { ProductList } from "@/components/product-list";

export default function ProductsPage() {
    return (
        <main className="container mx-auto px-4 py-8">
            <div className="flex flex-col items-center justify-center mb-8 text-center bg-white p-6 rounded-lg border shadow-sm">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Products</h1>
                    <p className="text-zinc-500 mt-2">Manage product catalog for packing.</p>
                </div>
            </div>

            <div className="grid gap-8 lg:grid-cols-12">
                <div className="lg:col-span-4">
                    <ProductForm />
                </div>
                <div className="lg:col-span-8">
                    <ProductList />
                </div>
            </div>
        </main>
    );
}
