import {
    Group,
    BoxGeometry,
    Color,
    MeshStandardMaterial,
    Mesh,
    EdgesGeometry,
    LineBasicMaterial,
    LineSegments,
} from "three";
import type { ItemData, PlacementData } from "../types";
import { mmToMeters, containerToThreeCoords, getItemCenterOffset } from "../utils/conversions";

export class ItemMesh {
    private group: Group;
    private itemData: ItemData;
    private placement: PlacementData;

    constructor(
        itemData: ItemData,
        placement: PlacementData,
        containerLength: number,
        containerWidth: number,
        containerHeight: number
    ) {
        this.itemData = itemData;
        this.placement = placement;
        this.group = new Group();
        this.createItem(containerLength, containerWidth, containerHeight);
    }

    private createItem(
        containerLength: number,
        containerWidth: number,
        containerHeight: number
    ): void {
        const length = mmToMeters(this.itemData.length_mm);
        const width = mmToMeters(this.itemData.width_mm);
        const height = mmToMeters(this.itemData.height_mm);

        // Create box geometry
        const geometry = new BoxGeometry(length, height, width);

        // Parse color
        const color = new Color(this.itemData.color_hex);

        // Create material
        const material = new MeshStandardMaterial({
            color: color,
            metalness: 0.1,
            roughness: 0.6,
            transparent: false,
            opacity: 1,
        });

        const mesh = new Mesh(geometry, material);
        mesh.userData = {
            label: this.itemData.label,
            step_number: this.placement.step_number,
        };


        // Create wireframe edges
        const edges = new EdgesGeometry(geometry);
        const lineMaterial = new LineBasicMaterial({
            color: 0xffffff,
            transparent: false,
            opacity: 1,
        });
        const wireframe = new LineSegments(edges, lineMaterial);

        // Position the item
        const [x, y, z] = containerToThreeCoords(
            this.placement.pos_x,
            this.placement.pos_y,
            this.placement.pos_z,
            containerLength,
            containerWidth,
            containerHeight
        );

        // Add center offset
        const [offsetX, offsetY, offsetZ] = getItemCenterOffset(
            this.itemData.length_mm,
            this.itemData.width_mm,
            this.itemData.height_mm
        );

        this.group.position.set(x + offsetX, y + offsetY, z + offsetZ);

        this.group.add(mesh);
        this.group.add(wireframe);
    }

    public getGroup(): Group {
        return this.group;
    }

    public getStepNumber(): number {
        return this.placement.step_number;
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
