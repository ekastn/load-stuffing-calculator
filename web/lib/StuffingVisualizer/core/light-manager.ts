import { Light, AmbientLight, DirectionalLight } from "three";

export class LightManager {
    private lights: Light[] = [];

    public createDefaultLights(): Light[] {
        // Ambient light for overall illumination
        const ambientLight = new AmbientLight(0xffffff, 0.8);

        // Directional light for shadows and depth
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
