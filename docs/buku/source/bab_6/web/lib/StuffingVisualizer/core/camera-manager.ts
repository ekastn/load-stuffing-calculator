import { OrthographicCamera } from "three";

/**
 * Manages the Three.js Camera
 */
export class CameraManager {
    public camera: OrthographicCamera;

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
}
