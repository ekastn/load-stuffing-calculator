import { OrthographicCamera, Scene, WebGLRenderer, Camera } from "three";
import { ContainerMesh } from "./components/container-mesh";
import { ItemMesh } from "./components/item-mesh";
import { CameraManager } from "./core/camera-manager";
import { LightManager } from "./core/light-manager";
import { SceneManager } from "./core/scene-manager";
import { RendererManager } from "./core/renderer-manager";
import { ControlsManager } from "./core/controls-manager";
import type { SceneConfig, StuffingPlanData } from "./types";
import { ReportGenerator, type ReportGeneratorConfig } from "./report/report-generator";
import { AnimationManager } from "./core/animation-manager";
import { InteractionManager, type HoverCallback } from "./core/interaction-manager";

export class StuffingVisualizer {
    private sceneManager: SceneManager;
    private cameraManager: CameraManager;
    private lightManager: LightManager;
    private rendererManager: RendererManager;
    private controlsManager: ControlsManager;
    private animationManager: AnimationManager;
    private interactionManager: InteractionManager;
    
    private containerMesh: ContainerMesh | null = null;
    private itemMeshes: ItemMesh[] = [];
    private data: StuffingPlanData | null = null;
    private camera: OrthographicCamera;
    private config: SceneConfig;

    constructor(config: SceneConfig = {}) {
        this.config = config;
        this.sceneManager = new SceneManager(config.backgroundColor);
        this.cameraManager = new CameraManager(
            1, 
            config.cameraNear ?? 0.1, 
            config.cameraFar ?? 2000
        );
        this.lightManager = new LightManager();
        this.rendererManager = new RendererManager();
        this.controlsManager = new ControlsManager();
        this.animationManager = new AnimationManager(config.stepDuration ?? 500);
        this.interactionManager = new InteractionManager();
        this.camera = this.cameraManager.getCamera();

        // Add lights to scene
        const lights = this.lightManager.createDefaultLights();
        lights.forEach((light) => this.sceneManager.add(light));

        // Connect internal scene update listener
        this.animationManager.onStepChange((step) => {
            this.updateVisibleItems(step);
        });
    }

    /**
     * Attaches the visualizer to a DOM element.
     * This initializes the renderer, controls, and starts the render loop.
     */
    public attach(container: HTMLElement): void {
        this.rendererManager.attach(container, (width, height, aspect) => {
            this.cameraManager.updateAspect(aspect);
        });
        
        const renderer = this.rendererManager.getRenderer();
        
        // Initialize controls
        this.controlsManager.init(this.camera, renderer);
        
        // Setup interaction
        this.interactionManager.setRenderer(renderer);
        this.interactionManager.setCamera(this.camera);
        
        // Start render loop
        this.rendererManager.startLoop(
            this.sceneManager.getScene(), 
            this.camera,
            () => this.controlsManager.update()
        );
    }

    public loadData(data: StuffingPlanData): void {
        this.data = data;
        this.clear();
        this.build();
    }

    private build(): void {
        if (!this.data) return;

        // Create container
        this.containerMesh = new ContainerMesh(this.data.container);
        this.sceneManager.add(this.containerMesh.getGroup());

        // Create items
        const itemMap = new Map(this.data.items.map((item) => [item.item_id, item]));
        let maxStep = 0;

        this.data.calculation.placements.forEach((placement) => {
            const itemData = itemMap.get(placement.item_id);
            if (itemData) {
                const itemMesh = new ItemMesh(
                    itemData,
                    placement,
                    this.data!.container.length_mm,
                    this.data!.container.width_mm,
                    this.data!.container.height_mm
                );
                this.itemMeshes.push(itemMesh);
                this.sceneManager.add(itemMesh.getGroup());
                maxStep = Math.max(maxStep, itemMesh.getStepNumber());
            }
        });

        this.animationManager.setMaxStep(maxStep);
        this.animationManager.setCurrentStep(maxStep);
        
        // Initial visibility
        this.updateVisibleItems(maxStep);
        
        // Update interaction manager
        this.interactionManager.setSceneChildren(this.sceneManager.getScene().children);
    }

    public setStep(step: number): void {
        this.animationManager.setCurrentStep(step);
    }

    public getMaxStep(): number {
        return this.animationManager.getMaxStep();
    }

    public getCurrentStep(): number {
        return this.animationManager.getCurrentStep();
    }

    private updateVisibleItems(targetStep: number): void {
        this.itemMeshes.forEach((itemMesh) => {
            const stepNumber = itemMesh.getStepNumber();
            itemMesh.getGroup().visible = stepNumber <= targetStep;
        });
    }

    public clear(): void {
        // Dispose container
        if (this.containerMesh) {
            this.containerMesh.dispose();
            this.containerMesh = null;
        }

        // Dispose items
        this.itemMeshes.forEach((item) => item.dispose());
        this.itemMeshes = [];

        // Clear scene (except lights)
        const lights = this.lightManager.getLights();
        this.sceneManager.clear();
        lights.forEach((light) => this.sceneManager.add(light));
        
        // Update interaction manager
        this.interactionManager.setSceneChildren(this.sceneManager.getScene().children);
    }

    public getScene(): Scene {
        return this.sceneManager.getScene();
    }

    public getCamera(): OrthographicCamera {
        return this.cameraManager.getCamera();
    }

    public dispose(): void {
        this.clear();
        this.lightManager.dispose();
        this.animationManager.dispose();
        this.interactionManager.dispose();
        this.controlsManager.dispose();
        this.rendererManager.dispose();
    }

    /**
     * Subscribe to item hover events.
     * @returns Unsubscribe function
     */
    public onItemHover(callback: HoverCallback): () => void {
        return this.interactionManager.onHover(callback);
    }

    public captureScreenshot(): string | null {
        const renderer = this.rendererManager.getRenderer();
        if (!renderer) return null;

        try {
            return renderer.domElement.toDataURL("image/png");
        } catch (error) {
            console.error("Failed to capture screenshot:", error);
            return null;
        }
    }

    public captureStepScreenshot(): { screenshot: string; aspectRatio: number } | null {
        const renderer = this.rendererManager.getRenderer();
        if (!renderer || !this.data) return null;

        try {
            const width = 1920;
            const height = 1080;

            // Delegate setup to camera manager
            this.cameraManager.setupForScreenshot(renderer, width, height, this.data.container);

            // Render
            renderer.render(this.sceneManager.getScene(), this.camera);

            // Capture
            const screenshot = renderer.domElement.toDataURL("image/png");
            const aspectRatio = width / height;

            // Restore
            this.cameraManager.restoreState(renderer);

            return { screenshot, aspectRatio };
        } catch (error) {
            console.error("Failed to capture screenshot:", error);
            return null;
        }
    }

    public async generateReport(config: Partial<ReportGeneratorConfig> = {}): Promise<Blob | null> {
        if (!this.data) return null;

        const generator = new ReportGenerator();

        // Capture a single full screenshot for the summary page.
        // Step pages remain vector-only to avoid memory blowups.
        const full = this.captureStepScreenshot();

        return generator.generateStuffingInstructions(this.data, {
            companyName: this.config.companyName ?? "Load Stuffing Visualization",
            fullScreenshot: full?.screenshot,
            fullScreenshotAspectRatio: full?.aspectRatio,
            ...config,
        });
    }

    // Animation Controls
    public play(): void {
        this.animationManager.play();
    }

    public pause(): void {
        this.animationManager.pause();
    }

    public togglePlay(): void {
        this.animationManager.toggle();
    }

    public reset(): void {
        this.animationManager.reset();
    }

    /**
     * Subscribe to step changes.
     * @returns Unsubscribe function
     */
    public onStepChange(callback: (step: number) => void): () => void {
        return this.animationManager.onStepChange(callback);
    }

    /**
     * Subscribe to play state changes.
     * @returns Unsubscribe function
     */
    public onPlayStateChange(callback: (isPlaying: boolean) => void): () => void {
        return this.animationManager.onPlayStateChange(callback);
    }
}
