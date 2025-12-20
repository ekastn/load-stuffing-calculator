import { Scene, Color, Object3D } from "three";

export class SceneManager {
    public scene: Scene;

    constructor(backgroundColor = "#1a1a1a") {
        this.scene = new Scene();
        this.scene.background = new Color(backgroundColor);
    }

    public add(object: Object3D): void {
        this.scene.add(object);
    }

    public remove(object: Object3D): void {
        this.scene.remove(object);
    }

    public clear(): void {
        while (this.scene.children.length > 0) {
            this.scene.remove(this.scene.children[0]);
        }
    }

    public getScene(): Scene {
        return this.scene;
    }
}
