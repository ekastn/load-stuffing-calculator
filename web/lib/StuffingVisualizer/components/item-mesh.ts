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

    private getRotatedDims(): { length_mm: number; width_mm: number; height_mm: number } {
        const l = this.itemData.length_mm;
        const w = this.itemData.width_mm;
        const h = this.itemData.height_mm;

        // rotation codes (0..5) come from our backend API and represent a
        // permutation of the item's original (length,width,height) in container
        // coordinates. They are not boxpacker3's native axis order.
        switch (this.placement.rotation) {
            case 0:
                return { length_mm: l, width_mm: w, height_mm: h };
            case 1:
                return { length_mm: w, width_mm: l, height_mm: h };
            case 2:
                return { length_mm: w, width_mm: h, height_mm: l };
            case 3:
                return { length_mm: h, width_mm: w, height_mm: l };
            case 4:
                return { length_mm: h, width_mm: l, height_mm: w };
            case 5:
                return { length_mm: l, width_mm: h, height_mm: w };
            default:
                return { length_mm: l, width_mm: w, height_mm: h };
        }
    }

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
        const rotated = this.getRotatedDims();

        const length = mmToMeters(rotated.length_mm);
        const width = mmToMeters(rotated.width_mm);
        const height = mmToMeters(rotated.height_mm);

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
            rotated.length_mm,
            rotated.width_mm,
            rotated.height_mm
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
