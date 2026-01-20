"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { FileDown, Loader2 } from "lucide-react";
import { PlanDetail, Container } from "@/lib/api/types";
import { generateStuffingReport } from "@/lib/pdf-report";

interface PDFExportButtonProps {
    plan: PlanDetail;
    container: Container;
    className?: string;
}

export function PDFExportButton({ plan, container, className }: PDFExportButtonProps) {
    const [generating, setGenerating] = useState(false);

    const handleExport = async () => {
        setGenerating(true);
        try {
            await generateStuffingReport(plan, container);
        } catch (error) {
            console.error("Failed to generate PDF", error);
            alert("Failed to generate PDF. See console for details.");
        } finally {
            setGenerating(false);
        }
    };

    return (
        <Button onClick={handleExport} variant="outline" className={className} disabled={generating}>
            {generating ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
                <FileDown className="mr-2 h-4 w-4" />
            )}
            {generating ? "Exporting..." : "Export PDF"}
        </Button>
    )
}
