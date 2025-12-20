import { Camera, WebGLRenderer } from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";

export class ControlsManager {
    private controls: OrbitControls | null = null;

    public init(camera: Camera, renderer: WebGLRenderer): void {
        this.controls = new OrbitControls(camera, renderer.domElement);
        
        // Default settings
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.05;
        this.controls.minDistance = 1;
        this.controls.maxPolarAngle = Math.PI / 2;
    }

    public update(): void {
        this.controls?.update();
    }

    public getControls(): OrbitControls | null {
        return this.controls;
    }

    public dispose(): void {
        this.controls?.dispose();
        this.controls = null;
    }
}
