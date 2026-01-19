import { Light, AmbientLight, DirectionalLight } from "three";

/**
 * Manages scene lighting
 */
export class LightManager {
    private lights: Light[] = [];

    public createDefaultLights(): Light[] {
        const ambientLight = new AmbientLight(0xffffff, 0.8);
        const directionalLight = new DirectionalLight(0xffffff, 0.5);
        directionalLight.position.set(20, 30, 20);

        this.lights = [ambientLight, directionalLight];
        return this.lights;
    }

    public getLights(): Light[] {
        return this.lights;
    }

    public dispose(): void {
        this.lights.forEach((light) => {
            if (light.parent) {
                light.parent.remove(light);
            }
        });
        this.lights = [];
    }
}
