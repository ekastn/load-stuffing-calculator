"use client";

import { useEffect, useRef, useState, forwardRef, useImperativeHandle } from "react";
import { StuffingVisualizer } from "@/lib/StuffingVisualizer";
import type { StuffingPlanData } from "@/lib/StuffingVisualizer";

interface LoadingViewerProps {
    data: StuffingPlanData;
    step: number;
}

export interface LoadingViewerRef {
    setStep: (step: number) => void;
}

export const LoadingViewer = forwardRef<LoadingViewerRef, LoadingViewerProps>(
    function LoadingViewer({ data, step }, ref) {
        const containerRef = useRef<HTMLDivElement>(null);
        const visualizerRef = useRef<StuffingVisualizer | null>(null);
        const [isLoaded, setIsLoaded] = useState(false);

        // Expose setStep method to parent via ref
        useImperativeHandle(ref, () => ({
            setStep: (newStep: number) => {
                if (visualizerRef.current && isLoaded) {
                    visualizerRef.current.setStep(newStep);
                }
            }
        }), [isLoaded]);

        useEffect(() => {
            if (!containerRef.current) return;

            // Initialize visualizer with minimalist settings
            if (!visualizerRef.current) {
                visualizerRef.current = new StuffingVisualizer({
                    backgroundColor: "#f9fafb", // slate-50
                    cameraFar: 100000
                });
                // Enable loading mode for step-by-step focus
                visualizerRef.current.setLoadingMode(true);
            }

            const visualizer = visualizerRef.current;

            // Load data
            visualizer.loadData(data);
            
            // Initial step
            visualizer.setStep(step);

            // Attach to DOM
            visualizer.attach(containerRef.current);
            setIsLoaded(true);

            // Cleanup
            return () => {
                visualizer.dispose();
                visualizerRef.current = null;
            };
        }, [data]);

        // Update step when it changes from props
        useEffect(() => {
            if (visualizerRef.current && isLoaded) {
                visualizerRef.current.setStep(step);
            }
        }, [step, isLoaded]);

        return (
            <div className="w-full h-full relative overflow-hidden bg-slate-50">
                <div ref={containerRef} className="w-full h-full" />
                
                {!isLoaded && (
                    <div className="absolute inset-0 flex items-center justify-center bg-slate-50">
                        <div className="flex flex-col items-center gap-3">
                            <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-primary" />
                            <p className="text-slate-500 text-sm font-medium">Initializing viewer...</p>
                        </div>
                    </div>
                )}
            </div>
        );
    }
);
