import { Camera, Raycaster, Vector2, WebGLRenderer, Object3D } from "three";

export type HoverCallback = (
    item: { label: string; step_number: number } | null, 
    x: number, 
    y: number
) => void;

export class InteractionManager {
    private raycaster: Raycaster;
    private mouse: Vector2;
    private camera: Camera | null = null;
    private renderer: WebGLRenderer | null = null;
    private sceneChildren: Object3D[] = [];
    private hoverListeners = new Set<HoverCallback>();
    private boundHandlePointerMove: (event: PointerEvent) => void;

    constructor() {
        this.raycaster = new Raycaster();
        this.mouse = new Vector2();
        this.boundHandlePointerMove = this.handlePointerMove.bind(this);
    }

    public setCamera(camera: Camera): void {
        this.camera = camera;
    }

    public setRenderer(renderer: WebGLRenderer): void {
        if (this.renderer === renderer) return;

        // Cleanup old listener
        if (this.renderer?.domElement) {
            this.renderer.domElement.removeEventListener("pointermove", this.boundHandlePointerMove);
        }

        this.renderer = renderer;

        // Attach new listener if domElement exists
        if (this.renderer?.domElement) {
            this.renderer.domElement.addEventListener("pointermove", this.boundHandlePointerMove);
        }
    }

    public setSceneChildren(children: Object3D[]): void {
        this.sceneChildren = children;
    }

    /**
     * Subscribe to hover events.
     * @returns Unsubscribe function
     */
    public onHover(callback: HoverCallback): () => void {
        this.hoverListeners.add(callback);
        return () => this.hoverListeners.delete(callback);
    }

    public dispose(): void {
        if (this.renderer?.domElement) {
            this.renderer.domElement.removeEventListener("pointermove", this.boundHandlePointerMove);
        }
        this.hoverListeners.clear();
        this.renderer = null;
        this.camera = null;
        this.sceneChildren = [];
    }

    private handlePointerMove(event: PointerEvent): void {
        if (!this.renderer?.domElement || !this.camera || this.hoverListeners.size === 0) return;

        // Calculate mouse position in normalized device coordinates (-1 to +1)
        const rect = this.renderer.domElement.getBoundingClientRect();
        
        // Guard against zero-sized element
        if (rect.width === 0 || rect.height === 0) return;

        const x = event.clientX - rect.left;
        const y = event.clientY - rect.top;

        this.mouse.x = (x / rect.width) * 2 - 1;
        this.mouse.y = -(y / rect.height) * 2 + 1;

        this.raycaster.setFromCamera(this.mouse, this.camera);

        const intersects = this.raycaster.intersectObjects(this.sceneChildren, true);

        let foundItem = null;
        for (const intersect of intersects) {
            let obj: Object3D | null = intersect.object;
            while (obj) {
                if (obj.userData && obj.userData.label) {
                    foundItem = obj.userData as { label: string; step_number: number };
                    break;
                }
                obj = obj.parent;
            }
            if (foundItem) break;
        }

        this.hoverListeners.forEach(callback => callback(foundItem, event.clientX, event.clientY));
    }
}