"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils"; 
import { Box, Truck, Map, LayoutDashboard } from "lucide-react";

export function Navigation() {
    const pathname = usePathname();

    const routes = [
        { href: "/containers", label: "Containers" },
        { href: "/products", label: "Products" },
        { href: "/plans", label: "Plans" },
    ];

    return (
        <header className="bg-white border-b">
            <div className="container mx-auto pt-6 pb-0">
                <div className="flex flex-col items-center gap-6">
                    <Link href="/" className="hover:opacity-80 transition-opacity">
                        <h1 className="text-2xl font-bold tracking-tight text-zinc-900">
                            Load & Stuffing Calculator
                        </h1>
                    </Link>
                    
                    <nav className="flex items-center gap-8 -mb-[1px]">
                        {routes.map(route => {
                             const isActive = pathname.startsWith(route.href);
                             return (
                                <Link
                                    key={route.href}
                                    href={route.href}
                                    className={cn(
                                        "pb-3 text-sm font-medium border-b-2 transition-colors px-2",
                                        isActive 
                                            ? "border-zinc-900 text-zinc-900" 
                                            : "border-transparent text-zinc-500 hover:text-zinc-700 hover:border-zinc-300"
                                    )}
                                >
                                    {route.label}
                                </Link>
                             )
                        })}
                    </nav>
                </div>
            </div>
        </header>
    );
}
