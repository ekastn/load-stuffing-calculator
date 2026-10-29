"use client";

import { Button } from "@/components/ui/button";
import { Slider } from "@/components/ui/slider";
import { Pause, Play, RotateCcw, Download, Loader2 } from "lucide-react";

interface AnimationControlsProps {
    currentStep: number;
    maxStep: number;
    isPlaying: boolean;
    isDownloading?: boolean;
    onStepChange: (step: number) => void;
    onPlay: () => void;
    onPause: () => void;
    onReset: () => void;
    onDownload: () => void;
}

export function AnimationControls({
    currentStep,
    maxStep,
    isPlaying,
    isDownloading = false,
    onStepChange,
    onPlay,
    onPause,
    onReset,
    onDownload,
}: AnimationControlsProps) {
    return (
        <div className="absolute bottom-6 left-1/2 -translate-x-1/2 bg-card/95 backdrop-blur-sm rounded-lg shadow-md border border-border p-4 flex flex-col sm:flex-row items-center gap-4 w-[calc(100%-2rem)] max-w-[500px] text-foreground">
            <div className="flex gap-2 w-full sm:w-auto justify-center">
                <Button
                    variant="outline"
                    size="icon"
                    onClick={isPlaying ? onPause : onPlay}
                    className="h-9 w-9 bg-card hover:bg-accent text-foreground border-border"
                >
                    {isPlaying ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
                </Button>
                <Button
                    variant="outline"
                    size="icon"
                    onClick={onReset}
                    className="h-9 w-9 bg-card hover:bg-accent text-foreground border-border"
                >
                    <RotateCcw className="h-4 w-4" />
                </Button>
                <Button
                    variant="outline"
                    size="icon"
                    onClick={onDownload}
                    disabled={isDownloading}
                    className="h-9 w-9 bg-card hover:bg-accent text-foreground border-border"
                >
                    {isDownloading ? (
                        <Loader2 className="h-4 w-4 animate-spin text-signal" />
                    ) : (
                        <Download className="h-4 w-4" />
                    )}
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
                <span className="text-sm font-mono tabular-nums text-foreground min-w-[60px] text-right">
                    {currentStep} / {maxStep}
                </span>
            </div>
        </div>
    );
}
