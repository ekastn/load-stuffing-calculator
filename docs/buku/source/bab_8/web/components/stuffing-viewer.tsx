"use client";

import { useEffect, useRef, useState } from "react";
import { StuffingVisualizer } from "@/lib/StuffingVisualizer";
import type { StuffingPlanData } from "@/lib/StuffingVisualizer";
import { AnimationControls } from "./animation-controls";

interface StuffingViewerProps {
    data: StuffingPlanData;
}

/**
 * React component that wraps the StuffingVisualizer
 */
export function StuffingViewer({ data }: StuffingViewerProps) {
    const containerRef = useRef<HTMLDivElement>(null);
    const visualizerRef = useRef<StuffingVisualizer | null>(null);
    
    const [isLoaded, setIsLoaded] = useState(false);
    const [currentStep, setCurrentStep] = useState(0);
    const [maxStep, setMaxStep] = useState(0);
    const [isPlaying, setIsPlaying] = useState(false);

    useEffect(() => {
        if (!containerRef.current) return;

        if (!visualizerRef.current) {
            visualizerRef.current = new StuffingVisualizer({
                backgroundColor: "#fff",
            });
        }

        const visualizer = visualizerRef.current;

        const unsubStep = visualizer.onStepChange((step) => {
            setCurrentStep(step);
        });

        const unsubPlay = visualizer.onPlayStateChange((playing) => {
            setIsPlaying(playing);
        });

        visualizer.loadData(data);
        const max = visualizer.getMaxStep();
        setMaxStep(max);
        setCurrentStep(max);

        visualizer.attach(containerRef.current);
        setIsLoaded(true);

        return () => {
            unsubStep();
            unsubPlay();
            visualizer.dispose();
            visualizerRef.current = null;
        };
    }, [data]);

    const handleStepChange = (step: number) => {
        setCurrentStep(step);
        visualizerRef.current?.setStep(step);
    };

    const handlePlay = () => {
        visualizerRef.current?.play();
    };

    const handlePause = () => {
        visualizerRef.current?.pause();
    };

    const handleReset = () => {
        visualizerRef.current?.reset();
    };

    return (
        <div className="w-full h-full relative">
            <div ref={containerRef} className="w-full h-full" />
            
            {!isLoaded && (
                <div className="absolute inset-0 flex items-center justify-center bg-white">
                    <p className="text-zinc-600">Loading visualization...</p>
                </div>
            )}

            <AnimationControls
                currentStep={currentStep}
                maxStep={maxStep}
                isPlaying={isPlaying}
                onStepChange={handleStepChange}
                onPlay={handlePlay}
                onPause={handlePause}
                onReset={handleReset}
            />
        </div>
    );
}
