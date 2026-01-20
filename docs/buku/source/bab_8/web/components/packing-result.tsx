import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface PackingResultProps {
    stats: {
        total_items: number;
        fitted_count: number;
        unfitted_count: number;
    }
}

export function PackingResult({ stats }: PackingResultProps) {
    const fitRate = stats.total_items > 0 
        ? ((stats.fitted_count / stats.total_items) * 100).toFixed(1) 
        : "0.0";

    return (
        <Card className="bg-zinc-50 border-dashed">
            <CardContent className="pt-6">
                <div className="grid grid-cols-3 gap-4 text-center">
                    <div>
                        <div className="text-2xl font-bold text-zinc-900">{stats.fitted_count}</div>
                        <div className="text-xs text-zinc-500 uppercase tracking-wider font-semibold">Packed</div>
                    </div>
                    <div>
                        <div className="text-2xl font-bold text-zinc-900">{stats.total_items}</div>
                        <div className="text-xs text-zinc-500 uppercase tracking-wider font-semibold">Total Items</div>
                    </div>
                    <div>
                        <div className="text-2xl font-bold text-green-600">{fitRate}%</div>
                        <div className="text-xs text-zinc-500 uppercase tracking-wider font-semibold">Fit Rate</div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}
