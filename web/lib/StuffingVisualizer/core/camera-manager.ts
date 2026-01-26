import { OrthographicCamera, WebGLRenderer, Vector2, Vector3, Quaternion } from "three";
import type { ContainerData } from "../types";

export class CameraManager {
    public camera: OrthographicCamera;
    
    // Saved state for screenshots
    private savedPosition = new Vector3();
    private savedQuaternion = new Quaternion();
    private savedZoom = 1;
    private savedLeft = 0;
    private savedRight = 0;
    private savedTop = 0;
    private savedBottom = 0;
    private savedSize = new Vector2();

    constructor(aspect = 1, near = 0.1, far = 2000) {
        const frustumSize = 20;
        const left = (-frustumSize * aspect) / 2;
        const right = (frustumSize * aspect) / 2;
        const top = frustumSize / 2;
        const bottom = -frustumSize / 2;

        this.camera = new OrthographicCamera(left, right, top, bottom, near, far);
        this.setDefaultPosition();
    }

    public setDefaultPosition(): void {
        this.camera.position.set(15, 8, -10);
        this.camera.zoom = 1;
        this.camera.lookAt(0, 0, 0);
    }

    public setPosition(x: number, y: number, z: number): void {
        this.camera.position.set(x, y, z);
    }

    public updateAspect(aspect: number): void {
        const frustumSize = 20;
        this.camera.left = (-frustumSize * aspect) / 2;
        this.camera.right = (frustumSize * aspect) / 2;
        this.camera.top = frustumSize / 2;
        this.camera.bottom = -frustumSize / 2;
        this.camera.updateProjectionMatrix();
    }

    public getCamera(): OrthographicCamera {
        return this.camera;
    }

    public fitToContainer(container: ContainerData): void {
        // Reset to default angle
        this.camera.position.set(15, 8, -10);
        this.camera.lookAt(0, 0, 0);

        // Convert dimensions to meters
        const length = container.length_mm / 1000;
        const width = container.width_mm / 1000;
        const height = container.height_mm / 1000;

        // Find maximum dimension to fit in frustum
        // Frustum size is hardcoded to 20 in constructor
        const maxDim = Math.max(length, width, height);
        
        // Calculate zoom needed. 
        // Frustum size 20 covers 20 units.
        // We want container (maxDim) to fill ~70% of screen.
        // 20 / zoom = maxDim / 0.7  => zoom = (20 * 0.7) / maxDim = 14 / maxDim
        const targetZoom = 14 / Math.max(maxDim, 1); // Prevent division by zero

        this.camera.zoom = targetZoom;
        this.camera.updateProjectionMatrix();
    }

    public setupForScreenshot(
        renderer: WebGLRenderer, 
        width: number, 
        height: number, 
        container: ContainerData
    ): void {
        // Save current state
        this.savedPosition.copy(this.camera.position);
        this.savedQuaternion.copy(this.camera.quaternion);
        this.savedZoom = this.camera.zoom;
        this.savedLeft = this.camera.left;
        this.savedRight = this.camera.right;
        this.savedTop = this.camera.top;
        this.savedBottom = this.camera.bottom;
        renderer.getSize(this.savedSize);

        // Set screenshot size
        renderer.setSize(width, height);

        // Update frustum for screenshot aspect ratio
        const aspect = width / height;
        const frustumSize = 20;
        this.camera.left = (-frustumSize * aspect) / 2;
        this.camera.right = (frustumSize * aspect) / 2;
        this.camera.top = frustumSize / 2;
        this.camera.bottom = -frustumSize / 2;

        // Set isometric view
        this.camera.position.set(15, 8, -10);
        this.camera.lookAt(0, 0, 0);

        // Calculate zoom to fit container
        const maxDim = Math.max(container.length_mm, container.width_mm, container.height_mm) / 1000;
        this.camera.zoom = 30 / maxDim;

        this.camera.updateProjectionMatrix();
    }

    public restoreState(renderer: WebGLRenderer): void {
        // Restore renderer size
        renderer.setSize(this.savedSize.x, this.savedSize.y);

        // Restore camera state
        this.camera.left = this.savedLeft;
        this.camera.right = this.savedRight;
        this.camera.top = this.savedTop;
        this.camera.bottom = this.savedBottom;
        this.camera.position.copy(this.savedPosition);
        this.camera.quaternion.copy(this.savedQuaternion);
        this.camera.zoom = this.savedZoom;
        this.camera.updateProjectionMatrix();
    }
}
