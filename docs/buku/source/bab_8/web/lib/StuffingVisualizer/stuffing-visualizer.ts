import { OrthographicCamera, Scene } from "three";
import { ContainerMesh } from "./components/container-mesh";
import { ItemMesh } from "./components/item-mesh";
import { CameraManager } from "./core/camera-manager";
import { LightManager } from "./core/light-manager";
import { SceneManager } from "./core/scene-manager";
import { RendererManager } from "./core/renderer-manager";
import { ControlsManager } from "./core/controls-manager";
import { AnimationManager } from "./core/animation-manager";
import type { SceneConfig, StuffingPlanData, ItemData } from "./types";

/**
 * Main orchestrator class for 3D stuffing visualization
 */
export class StuffingVisualizer {
    private sceneManager: SceneManager;
    private cameraManager: CameraManager;
    private lightManager: LightManager;
    private rendererManager: RendererManager;
    private controlsManager: ControlsManager;
    private animationManager: AnimationManager;
    
    private containerMesh: ContainerMesh | null = null;
    private itemMeshes: ItemMesh[] = [];
    private data: StuffingPlanData | null = null;
    private camera: OrthographicCamera;

    constructor(config: SceneConfig = {}) {
        this.sceneManager = new SceneManager(config.backgroundColor);
        this.cameraManager = new CameraManager(1, 0.1, 2000);
        this.lightManager = new LightManager();
        this.rendererManager = new RendererManager();
        this.controlsManager = new ControlsManager();
        this.animationManager = new AnimationManager(config.stepDuration ?? 500);
        this.camera = this.cameraManager.getCamera();

        const lights = this.lightManager.createDefaultLights();
        lights.forEach((light) => this.sceneManager.add(light));

        this.animationManager.onStepChange((step) => {
            this.updateVisibleItems(step);
        });
    }

    public attach(container: HTMLElement): void {
        this.rendererManager.attach(container, (_width, _height, aspect) => {
            this.cameraManager.updateAspect(aspect);
        });
        
        const renderer = this.rendererManager.getRenderer();
        this.controlsManager.init(this.camera, renderer);
        
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

        this.containerMesh = new ContainerMesh(this.data.container);
        this.sceneManager.add(this.containerMesh.getGroup());

        const itemMap = new Map<string, ItemData>(
            this.data.items.map((item) => [item.item_id, item])
        );
        let maxStep = 0;

        this.data.placements.forEach((placement) => {
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
        this.updateVisibleItems(maxStep);
    }

    private updateVisibleItems(targetStep: number): void {
        this.itemMeshes.forEach((itemMesh) => {
            const stepNumber = itemMesh.getStepNumber();
            itemMesh.getGroup().visible = stepNumber <= targetStep;
        });
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

    public play(): void {
        this.animationManager.play();
    }

    public pause(): void {
        this.animationManager.pause();
    }

    public reset(): void {
        this.animationManager.reset();
    }

    public onStepChange(callback: (step: number) => void): () => void {
        return this.animationManager.onStepChange(callback);
    }

    public onPlayStateChange(callback: (isPlaying: boolean) => void): () => void {
        return this.animationManager.onPlayStateChange(callback);
    }

    public clear(): void {
        if (this.containerMesh) {
            this.containerMesh.dispose();
            this.containerMesh = null;
        }

        this.itemMeshes.forEach((item) => item.dispose());
        this.itemMeshes = [];

        const lights = this.lightManager.getLights();
        this.sceneManager.clear();
        lights.forEach((light) => this.sceneManager.add(light));
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
        this.controlsManager.dispose();
        this.rendererManager.dispose();
    }
}
