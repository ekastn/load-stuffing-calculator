import { useEffect, useRef, useState, forwardRef, useImperativeHandle } from "react";
import { StuffingVisualizer } from "@/lib/StuffingVisualizer";
import type { StuffingPlanData } from "@/lib/StuffingVisualizer";
import { AnimationControls } from "./animation-controls";
import { Loader2 } from "lucide-react";

interface StuffingViewerProps {
    data: StuffingPlanData;
}

const DownloadHandler = forwardRef(
    (
        props: {
            visualizer: StuffingVisualizer;
            data: StuffingPlanData;
            currentStep: number;
            setIsDownloading: (isDownloading: boolean) => void;
        },
        ref
    ) => {
        useImperativeHandle(ref, () => ({
            handleDownload: async () => {
                if (!props.visualizer) return;
                props.setIsDownloading(true);
                try {
                    // Give UI time to show loading state
                    await new Promise((resolve) => setTimeout(resolve, 100));
                    
                    const pdfBlob = await props.visualizer.generateReport({
                        companyName: "Load Stuffing Visualization",
                    });
                    
                    if (pdfBlob) {
                        const url = URL.createObjectURL(pdfBlob);
                        const link = document.createElement("a");
                        link.href = url;
                        link.download = `stuffing-plan-${props.data.plan_code}.pdf`;
                        document.body.appendChild(link);
                        link.click();
                        document.body.removeChild(link);
                        URL.revokeObjectURL(url);
                    }
                } catch (error) {
                    console.error("Failed to generate PDF:", error);
                    alert("Failed to generate PDF. Please try again.");
                } finally {
                    props.setIsDownloading(false);
                }
            },
        }));

        return null;
    }
);

export function StuffingViewer({ data }: StuffingViewerProps) {
    const containerRef = useRef<HTMLDivElement>(null);
    const visualizerRef = useRef<StuffingVisualizer | null>(null);
    const downloadHandlerRef = useRef<any>(null);
    
    const [isLoaded, setIsLoaded] = useState(false);
    const [currentStep, setCurrentStep] = useState(0);
    const [maxStep, setMaxStep] = useState(0);
    const [isPlaying, setIsPlaying] = useState(false);
    const [isDownloading, setIsDownloading] = useState(false);
    const [hoveredItem, setHoveredItem] = useState<{ label: string; step_number: number } | null>(
        null
    );
    const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

    useEffect(() => {
        if (!containerRef.current) return;

        // Initialize visualizer
        if (!visualizerRef.current) {
            visualizerRef.current = new StuffingVisualizer({
                backgroundColor: "#fff",
                cameraFar: 100000
            });
        }

        const visualizer = visualizerRef.current;

        // Setup listeners
        const unsubStep = visualizer.onStepChange((step) => {
            setCurrentStep(step);
        });

        const unsubPlay = visualizer.onPlayStateChange((playing) => {
            setIsPlaying(playing);
        });

        const unsubHover = visualizer.onItemHover((item, x, y) => {
            setHoveredItem(item);
            setMousePosition({ x, y });
        });

        // Load data
        visualizer.loadData(data);
        const max = visualizer.getMaxStep();
        setMaxStep(max);
        setCurrentStep(max);

        // Attach to DOM
        visualizer.attach(containerRef.current);
        setIsLoaded(true);

        // Cleanup
        return () => {
            unsubStep();
            unsubPlay();
            unsubHover();
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

    const handleDownloadClick = () => {
        if (downloadHandlerRef.current) {
            downloadHandlerRef.current.handleDownload();
        }
    };

    return (
        <div className="w-full h-full relative">
            <div ref={containerRef} className="w-full h-full" />
            
            {!isLoaded && (
                <div className="absolute inset-0 flex items-center justify-center bg-white">
                    <p className="text-zinc-600">Loading visualization...</p>
                </div>
            )}

            {visualizerRef.current && (
                <DownloadHandler
                    ref={downloadHandlerRef}
                    visualizer={visualizerRef.current}
                    data={data}
                    currentStep={currentStep}
                    setIsDownloading={setIsDownloading}
                />
            )}

            {isDownloading && (
                <div className="absolute inset-0 bg-white/60 backdrop-blur-sm z-[2000] flex flex-col items-center justify-center gap-4">
                    <Loader2 className="h-10 w-10 animate-spin text-blue-600" />
                    <p className="text-lg font-medium text-zinc-900">
                        Generating PDF Instructions...
                    </p>
                    <p className="text-sm text-zinc-500 text-center max-w-xs">
                        Capturing 3D steps and calculating layout. This may take a few moments.
                    </p>
                </div>
            )}

            <AnimationControls
                currentStep={currentStep}
                maxStep={maxStep}
                isPlaying={isPlaying}
                isDownloading={isDownloading}
                onStepChange={handleStepChange}
                onPlay={handlePlay}
                onPause={handlePause}
                onReset={handleReset}
                onDownload={handleDownloadClick}
            />
            
            {hoveredItem && (
                <div
                    className="absolute bg-white p-2 rounded shadow-lg"
                    style={{
                        left: mousePosition.x + 10,
                        top: mousePosition.y + 10,
                        pointerEvents: "none",
                    }}
                >
                    <p className="font-bold">{hoveredItem.label}</p>
                    <p>Step: {hoveredItem.step_number}</p>
                </div>
            )}
        </div>
    );
}
