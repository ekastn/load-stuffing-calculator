export type StepChangeCallback = (step: number) => void;
export type PlayStateChangeCallback = (isPlaying: boolean) => void;

export class AnimationManager {
    private isPlaying = false;
    private currentStep = 0;
    private maxStep = 0;
    private stepDuration = 500;
    private lastStepTime = 0;
    private animationFrameId: number | null = null;
    
    private stepListeners = new Set<StepChangeCallback>();
    private playStateListeners = new Set<PlayStateChangeCallback>();

    constructor(initialStepDuration: number = 500) {
        this.stepDuration = initialStepDuration;
    }

    public setMaxStep(max: number): void {
        this.maxStep = max;
    }

    public getMaxStep(): number {
        return this.maxStep;
    }

    public setCurrentStep(step: number): void {
        const newStep = Math.max(0, Math.min(step, this.maxStep));
        if (this.currentStep === newStep) return;
        
        this.currentStep = newStep;
        this.notifyStepChange();
    }

    public getCurrentStep(): number {
        return this.currentStep;
    }

    /**
     * Subscribe to step changes.
     * @returns Unsubscribe function
     */
    public onStepChange(callback: StepChangeCallback): () => void {
        this.stepListeners.add(callback);
        return () => this.stepListeners.delete(callback);
    }

    /**
     * Subscribe to play/pause state changes.
     * @returns Unsubscribe function
     */
    public onPlayStateChange(callback: PlayStateChangeCallback): () => void {
        this.playStateListeners.add(callback);
        return () => this.playStateListeners.delete(callback);
    }

    public play(): void {
        if (this.isPlaying) return;

        if (this.currentStep >= this.maxStep) {
            this.setCurrentStep(0);
        }

        this.isPlaying = true;
        this.notifyPlayStateChange();
        this.lastStepTime = Date.now();
        this.animate();
    }

    public pause(): void {
        if (!this.isPlaying) return;
        
        this.isPlaying = false;
        this.notifyPlayStateChange();
        if (this.animationFrameId !== null) {
            cancelAnimationFrame(this.animationFrameId);
            this.animationFrameId = null;
        }
    }

    public toggle(): void {
        if (this.isPlaying) {
            this.pause();
        } else {
            this.play();
        }
    }

    public reset(): void {
        this.pause();
        this.setCurrentStep(0);
    }

    public dispose(): void {
        this.pause();
        this.stepListeners.clear();
        this.playStateListeners.clear();
    }

    public isRunning(): boolean {
        return this.isPlaying;
    }

    private notifyStepChange(): void {
        this.stepListeners.forEach(callback => callback(this.currentStep));
    }

    private notifyPlayStateChange(): void {
        this.playStateListeners.forEach(callback => callback(this.isPlaying));
    }

    private animate = (): void => {
        if (!this.isPlaying) return;

        const now = Date.now();
        const elapsed = now - this.lastStepTime;

        if (elapsed >= this.stepDuration) {
            this.lastStepTime = now;
            const nextStep = this.currentStep + 1;

            if (nextStep > this.maxStep) {
                this.pause();
                this.setCurrentStep(this.maxStep);
                return;
            } else {
                this.setCurrentStep(nextStep);
            }
        }

        this.animationFrameId = requestAnimationFrame(this.animate);
    }
}
