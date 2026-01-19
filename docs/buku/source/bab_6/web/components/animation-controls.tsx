"use client";

import { Button } from "@/components/ui/button";
import { Slider } from "@/components/ui/slider";
import { Pause, Play, RotateCcw } from "lucide-react";

interface AnimationControlsProps {
    currentStep: number;
    maxStep: number;
    isPlaying: boolean;
    onStepChange: (step: number) => void;
    onPlay: () => void;
    onPause: () => void;
    onReset: () => void;
}

export function AnimationControls({
    currentStep,
    maxStep,
    isPlaying,
    onStepChange,
    onPlay,
    onPause,
    onReset,
}: AnimationControlsProps) {
    return (
        <div className="absolute bottom-6 left-1/2 -translate-x-1/2 bg-white rounded-lg shadow-lg border border-zinc-200 p-4 flex items-center gap-4 min-w-[400px]">
            <div className="flex gap-2">
                <Button
                    variant="outline"
                    size="icon"
                    onClick={isPlaying ? onPause : onPlay}
                    className="h-9 w-9"
                >
                    {isPlaying ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
                </Button>
                <Button
                    variant="outline"
                    size="icon"
                    onClick={onReset}
                    className="h-9 w-9"
                >
                    <RotateCcw className="h-4 w-4" />
                </Button>
            </div>

            <div className="flex-1 flex items-center gap-3">
                <Slider
                    value={[currentStep]}
                    min={0}
                    max={maxStep}
                    step={1}
                    onValueChange={(value) => onStepChange(value[0])}
                    className="flex-1"
                />
                <span className="text-sm font-medium text-zinc-700 min-w-[60px] text-right">
                    {currentStep} / {maxStep}
                </span>
            </div>
        </div>
    );
}
