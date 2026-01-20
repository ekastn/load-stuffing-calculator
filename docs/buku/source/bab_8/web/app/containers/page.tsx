import { ContainerForm } from "@/components/container-form";
import { ContainerList } from "@/components/container-list";

export default function ContainersPage() {
    return (
        <main className="container mx-auto px-4 py-8">
            <div className="flex flex-col items-center justify-center mb-8 text-center bg-white p-6 rounded-lg border shadow-sm">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Containers</h1>
                    <p className="text-zinc-500 mt-2">Manage available containers for packing.</p>
                </div>
            </div>

            <div className="grid gap-8 lg:grid-cols-12">
                <div className="lg:col-span-4">
                    <ContainerForm />
                </div>
                <div className="lg:col-span-8">
                    <ContainerList />
                </div>
            </div>
        </main>
    );
}
