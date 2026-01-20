import {
    Group,
    BoxGeometry,
    EdgesGeometry,
    LineBasicMaterial,
    LineSegments,
    MeshBasicMaterial,
    Mesh,
} from "three";
import type { ContainerData } from "../types";
import { mmToMeters } from "../utils/conversions";

/**
 * Creates a 3D wireframe representation of the container
 */
export class ContainerMesh {
    private group: Group;
    private data: ContainerData;

    constructor(containerData: ContainerData) {
        this.data = containerData;
        this.group = new Group();
        this.createContainer();
    }

    private createContainer(): void {
        const length = mmToMeters(this.data.length_mm);
        const width = mmToMeters(this.data.width_mm);
        const height = mmToMeters(this.data.height_mm);

        // Wireframe box for the container
        const geometry = new BoxGeometry(length, height, width);
        const edges = new EdgesGeometry(geometry);
        const lineMaterial = new LineBasicMaterial({
            color: 0x888888,
            linewidth: 1,
        });
        const wireframe = new LineSegments(edges, lineMaterial);

        // Floor platform
        const platformGeometry = new BoxGeometry(length, 0.02, width);
        const platformMaterial = new MeshBasicMaterial({
            color: 0xff8800,
        });
        const platform = new Mesh(platformGeometry, platformMaterial);
        platform.position.y = -height / 2 - 0.01;

        this.group.add(wireframe);
        this.group.add(platform);
    }

    public getGroup(): Group {
        return this.group;
    }

    public dispose(): void {
        this.group.traverse((child) => {
            if (child instanceof Mesh || child instanceof LineSegments) {
                child.geometry.dispose();
                if (Array.isArray(child.material)) {
                    child.material.forEach((m) => m.dispose());
                } else {
                    child.material.dispose();
                }
            }
        });
    }
}
