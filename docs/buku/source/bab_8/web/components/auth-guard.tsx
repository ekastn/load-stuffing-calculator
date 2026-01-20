"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";

export function AuthGuard({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const pathname = usePathname();
    const [authorized, setAuthorized] = useState(false);

    useEffect(() => {
        // Public paths that don't need auth
        const publicPaths = ["/login", "/register"];
        
        // Check if current path is public
        const isPublic = publicPaths.some(path => pathname.startsWith(path));
        
        // Get token from storage
        const token = localStorage.getItem("auth_token");

        if (isPublic) {
            // If already logged in and trying to access public page, redirect to home
            if (token) {
                router.push("/");
            } else {
                setAuthorized(true);
            }
        } else {
            // Protected page
            if (!token) {
                router.push("/login");
            } else {
                setAuthorized(true);
            }
        }
    }, [pathname, router]);

    // Simple loading state while checking
    if (!authorized) {
        return null; // Or a loading spinner
    }

    return <>{children}</>;
}
