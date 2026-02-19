import type { StuffingPlanData, ItemData, PlacementData } from "../types";

export interface ReportGeneratorConfig {
    brandLogoUrl?: string;
    companyName?: string;

    /**
     * Legacy raster export options (kept for backwards compatibility).
     * The current implementation draws vector schematics instead of embedding images.
     */
    stepImages?: string[];
    stepAspectRatios?: number[];
    fullScreenshot?: string;
    fullScreenshotAspectRatio?: number;
}

export class ReportGenerator {
    private readonly ITEM_TEXT_COLOR = [0, 0, 0];
    private jsPDF: any = null;

    private readonly BRAND_BLUE = [0, 64, 255]; // RGB for blue color
    private readonly HEADER_BLUE = [30, 144, 255]; // Lighter blue for headers

    async initialize() {
        if (this.jsPDF) return;

        // Dynamically import jsPDF
        const { default: jsPDF } = await import("jspdf");
        this.jsPDF = jsPDF;
    }

    async generateStuffingInstructions(
        data: StuffingPlanData,
        config: ReportGeneratorConfig = {}
    ): Promise<Blob> {
        await this.initialize();

        const doc = new this.jsPDF({
            orientation: "landscape",
            unit: "mm",
            format: "a4",
        });

        const pageWidth = doc.internal.pageSize.getWidth();
        const pageHeight = doc.internal.pageSize.getHeight();
        const margin = 20;

        this.createSummaryPage(doc, data, config, pageWidth, pageHeight, margin);

        this.createStepPages(doc, data, config, pageWidth, pageHeight, margin);

        return doc.output("blob");
    }

    private createSummaryPage(
        doc: any,
        data: StuffingPlanData,
        config: ReportGeneratorConfig,
        pageWidth: number,
        pageHeight: number,
        margin: number
    ) {
        let yPosition = margin;

        doc.setTextColor(...this.BRAND_BLUE);
        doc.setFontSize(28);
        doc.setFont(undefined, "bold");
        doc.text("STUFFING RESULT", margin, yPosition);

        // Logo on the right (if provided)
        if (config.brandLogoUrl) {
            try {
                doc.addImage(
                    config.brandLogoUrl,
                    "PNG",
                    pageWidth - margin - 40,
                    yPosition - 8,
                    40,
                    20
                );
            } catch (error) {
                console.error("Failed to add logo:", error);
            }
        }

        yPosition += 15;

        const containerX = margin;
        const containerY = yPosition;
        const containerWidth = pageWidth - 2 * margin;
        const containerHeight = 120;

        // Draw border around entire content area
        doc.setDrawColor(200, 200, 200);
        doc.setLineWidth(0.5);
        doc.rect(containerX, containerY, containerWidth, containerHeight);

        // Left side: Container visualization (35% of width for better image display)
        const leftWidth = containerWidth * 0.35;
        doc.line(
            containerX + leftWidth,
            containerY,
            containerX + leftWidth,
            containerY + containerHeight
        );

        yPosition = containerY + 12;

        // Container name
        doc.setTextColor(0, 0, 0);
        doc.setFontSize(16);
        doc.setFont(undefined, "bold");
        doc.text(data.container.name, containerX + 15, yPosition);

        yPosition += 15;

        if (config.fullScreenshot) {
            try {
                const imgWidth = leftWidth - 30;
                const aspectRatio = config.fullScreenshotAspectRatio || 1;
                const imgHeight = imgWidth / aspectRatio;
                const imgX = containerX + 15;
                doc.addImage(config.fullScreenshot, "PNG", imgX, yPosition, imgWidth, imgHeight);
                yPosition += imgHeight + 8;
            } catch (error) {
                console.error("Failed to add full screenshot:", error);
                yPosition += 50;
            }
        } else {
            yPosition += 50;
        }

        // Container count
        doc.setFontSize(12);
        doc.setFont(undefined, "normal");
        doc.text("1 unit", containerX + 15, yPosition);

        let rightY = containerY + 12;
        const rightX = containerX + leftWidth + 20;

        // Summary stats with better spacing
        doc.setFontSize(11);
        doc.setFont(undefined, "bold");
        doc.text("Total:", rightX, rightY);
        doc.setFont(undefined, "normal");
        doc.text(`${data.calculation.placements.length} packages`, rightX + 50, rightY);
        rightY += 8;

        doc.setFont(undefined, "bold");
        doc.text("Cargo volume:", rightX, rightY);
        doc.setFont(undefined, "normal");
        doc.text(
            `${((data.container.volume_m3 * data.calculation.volume_utilization_pct) / 100).toFixed(2)} m3`,
            rightX + 50,
            rightY
        );
        rightY += 8;

        doc.setFont(undefined, "bold");
        doc.text("Cargo weight:", rightX, rightY);
        doc.setFont(undefined, "normal");
        const totalWeight = data.calculation.placements.reduce((sum, p) => {
            const item = data.items.find((i) => i.item_id === p.item_id);
            return sum + (item?.weight_kg || 0);
        }, 0);
        doc.text(`${totalWeight.toFixed(2)} kg`, rightX + 50, rightY);
        rightY += 12;

        doc.setDrawColor(200, 200, 200);
        doc.setLineWidth(0.3);
        const tableStartY = rightY;
        doc.line(rightX, rightY, containerX + containerWidth - 20, rightY);
        rightY += 6;

        // Calculate column positions based on available width
        const availableWidth = (containerX + containerWidth - 20) - rightX;
        const colSKU = rightX;
        const colName = rightX + (availableWidth * 0.12);
        const colDims = rightX + (availableWidth * 0.38);
        const colPackages = rightX + (availableWidth * 0.52);
        const colVolume = rightX + (availableWidth * 0.62);
        const colWeight = rightX + (availableWidth * 0.75);

        doc.setFontSize(9);
        doc.setFont(undefined, "bold");
        doc.text("SKU", colSKU, rightY);
        doc.text("Name", colName, rightY);
        doc.text("Unit Dims", colDims, rightY);
        doc.text("(mm)", colDims, rightY + 3);
        doc.text("Count", colPackages, rightY);
        doc.text("Volume", colVolume, rightY);
        doc.text("Weight", colWeight, rightY);
        rightY += 6;

        doc.line(rightX, rightY, containerX + containerWidth - 20, rightY);
        rightY += 4;

        // Items data with better formatting
        doc.setFontSize(8);
        doc.setFont(undefined, "normal");

        const itemCounts = new Map<string, { item: ItemData; count: number }>();
        data.calculation.placements.forEach((placement) => {
            const item = data.items.find((i) => i.item_id === placement.item_id);
            if (item) {
                const existing = itemCounts.get(item.item_id);
                if (existing) {
                    existing.count++;
                } else {
                    itemCounts.set(item.item_id, { item, count: 1 });
                }
            }
        });

        itemCounts.forEach(({ item, count }) => {
            const totalVolume =
                (item.length_mm * item.width_mm * item.height_mm * count) / 1_000_000_000;
            const totalWeight = item.weight_kg * count;
            const unitDims = `${item.length_mm}×${item.width_mm}×${item.height_mm}`;
            const sku = item.sku || item.product_sku || "N/A";

            doc.text(sku, colSKU, rightY);
            doc.text(item.label.substring(0, 20), colName, rightY);
            doc.text(unitDims, colDims, rightY);
            doc.text(count.toString(), colPackages + 3, rightY);
            doc.text(`${totalVolume.toFixed(1)} m3`, colVolume, rightY);
            doc.text(`${totalWeight.toFixed(1)} kg`, colWeight, rightY);
            rightY += 5;
        });
    }

    private createStepPages(
        doc: any,
        data: StuffingPlanData,
        config: ReportGeneratorConfig,
        pageWidth: number,
        pageHeight: number,
        margin: number
    ) {
        const sortedPlacements = [...data.calculation.placements].sort(
            (a, b) => a.step_number - b.step_number
        );

        for (let i = 0; i < sortedPlacements.length; i++) {
            doc.addPage();
            this.renderStepLandscape(doc, sortedPlacements[i], data, config, pageWidth, pageHeight, margin);
        }

        // Footer on last page
        doc.setTextColor(120, 120, 120);
        doc.setFontSize(8);
        doc.setFont(undefined, "italic");
        doc.text(`Generated on ${new Date().toLocaleString()}`, margin, pageHeight - 10);
    }

    private renderStepLandscape(
        doc: any,
        placement: PlacementData,
        data: StuffingPlanData,
        _config: ReportGeneratorConfig,
        pageWidth: number,
        pageHeight: number,
        margin: number
    ) {
        const container = data.container;
        const placements = data.calculation.placements;

        const itemMap = new Map<string, ItemData>(data.items.map((i) => [i.item_id, i]));
        const currentItem = itemMap.get(placement.item_id);

        let yPosition = margin + 5;

        // Title
        doc.setTextColor(...this.BRAND_BLUE);
        doc.setFontSize(20);
        doc.setFont(undefined, "bold");
        const stepText = `Step ${placement.step_number}`;
        const textWidth = doc.getTextWidth(stepText);
        doc.text(stepText, (pageWidth - textWidth) / 2, yPosition);
        yPosition += 10;

        // Layout
        const gap = 10;
        const legendHeight = 28;
        const diagramTop = yPosition;
        const diagramHeight = pageHeight - margin - legendHeight - diagramTop;

        const diagramWidth = pageWidth - 2 * margin;
        const diagramBoxWidth = (diagramWidth - gap) / 2;

        const topViewBox = {
            x: margin,
            y: diagramTop,
            w: diagramBoxWidth,
            h: diagramHeight,
        };

        const sideViewBox = {
            x: margin + diagramBoxWidth + gap,
            y: diagramTop,
            w: diagramBoxWidth,
            h: diagramHeight,
        };

        // Draw frames
        doc.setDrawColor(220, 220, 220);
        doc.setLineWidth(0.3);
        doc.rect(topViewBox.x, topViewBox.y, topViewBox.w, topViewBox.h);
        doc.rect(sideViewBox.x, sideViewBox.y, sideViewBox.w, sideViewBox.h);

        // Labels
        doc.setTextColor(60, 60, 60);
        doc.setFontSize(10);
        doc.setFont(undefined, "bold");
        doc.text("Top view (X/Y)", topViewBox.x + 4, topViewBox.y + 6);
        doc.text("Side view (X/Z)", sideViewBox.x + 4, sideViewBox.y + 6);

        // Inner drawing areas
        const innerPadding = 10;
        const topInner = {
            x: topViewBox.x + innerPadding,
            y: topViewBox.y + 10,
            w: topViewBox.w - 2 * innerPadding,
            h: topViewBox.h - 14,
        };

        const sideInner = {
            x: sideViewBox.x + innerPadding,
            y: sideViewBox.y + 10,
            w: sideViewBox.w - 2 * innerPadding,
            h: sideViewBox.h - 14,
        };

        const fitScale = (contentW: number, contentH: number, boxW: number, boxH: number) => {
            const scale = Math.min(boxW / contentW, boxH / contentH);
            const drawW = contentW * scale;
            const drawH = contentH * scale;
            return {
                scale,
                offsetX: (boxW - drawW) / 2,
                offsetY: (boxH - drawH) / 2,
            };
        };

        const topFit = fitScale(container.length_mm, container.width_mm, topInner.w, topInner.h);
        const sideFit = fitScale(container.length_mm, container.height_mm, sideInner.w, sideInner.h);

        const topOriginX = topInner.x + topFit.offsetX;
        const topOriginY = topInner.y + topFit.offsetY;
        const sideOriginX = sideInner.x + sideFit.offsetX;
        const sideOriginY = sideInner.y + sideFit.offsetY;

        // Container outlines
        doc.setDrawColor(40, 40, 40);
        doc.setLineWidth(0.6);
        const topDrawW = container.length_mm * topFit.scale;
        const topDrawH = container.width_mm * topFit.scale;
        const sideDrawW = container.length_mm * sideFit.scale;
        const sideDrawH = container.height_mm * sideFit.scale;

        doc.rect(topOriginX, topOriginY, topDrawW, topDrawH);
        doc.rect(sideOriginX, sideOriginY, sideDrawW, sideDrawH);

        // Axis helpers (origin bottom-left, "math" orientation)
        this.drawAxes(doc, {
            originX: topOriginX,
            originY: topOriginY + topDrawH,
            xLen: topDrawW,
            yLen: topDrawH,
            xLabel: "X (mm)",
            yLabel: "Y (mm)",
            xMaxLabel: `${container.length_mm.toFixed(0)}`,
            yMaxLabel: `${container.width_mm.toFixed(0)}`,
        });

        this.drawAxes(doc, {
            originX: sideOriginX,
            originY: sideOriginY + sideDrawH,
            xLen: sideDrawW,
            yLen: sideDrawH,
            xLabel: "X (mm)",
            yLabel: "Z (mm)",
            xMaxLabel: `${container.length_mm.toFixed(0)}`,
            yMaxLabel: `${container.height_mm.toFixed(0)}`,
        });

        const visiblePlacements = placements
            .filter((p) => p.step_number <= placement.step_number)
            .sort((a, b) => a.step_number - b.step_number);

        for (const p of visiblePlacements) {
            const isCurrent = p.step_number === placement.step_number;
            const item = itemMap.get(p.item_id);
            if (!item) continue;

            const dims = this.getDimsForRotation(item, p.rotation);

            // Coordinate system notes:
            // - PDF Y grows downward.
            // - For a "natural" diagram orientation, we draw higher coordinates nearer the top
            //   by inverting the axis using the container dimension.
            // - We draw with a "math" orientation: origin at bottom-left,
            //   +X to the right, and +Y/+Z upward.

            // Top view (X/Y): +X right, +Y up.
            const topX = topOriginX + p.pos_x * topFit.scale;
            const topY = topOriginY + (container.width_mm - (p.pos_y + dims.width_mm)) * topFit.scale;

            this.drawItemRect(
                doc,
                {
                    x: topX,
                    y: topY,
                    w: dims.length_mm * topFit.scale,
                    h: dims.width_mm * topFit.scale,
                },
                item.color_hex,
                isCurrent
            );

            // Side view (X/Z): +X right, +Z up.
            const sideX = sideOriginX + p.pos_x * sideFit.scale;
            const sideY = sideOriginY + (container.height_mm - (p.pos_z + dims.height_mm)) * sideFit.scale;

            this.drawItemRect(
                doc,
                {
                    x: sideX,
                    y: sideY,
                    w: dims.length_mm * sideFit.scale,
                    h: dims.height_mm * sideFit.scale,
                },
                item.color_hex,
                isCurrent
            );
        }

        // Legend/info
        const infoY = pageHeight - margin - legendHeight + 6;

        doc.setTextColor(0, 0, 0);
        doc.setFontSize(12);
        doc.setFont(undefined, "bold");
        const label = currentItem?.label ?? "Unknown item";
        doc.text(label, margin, infoY);

        doc.setFont(undefined, "normal");
        doc.setFontSize(10);

        if (currentItem) {
            const dims = this.getDimsForRotation(currentItem, placement.rotation);
            const sku = currentItem.sku || currentItem.product_sku || "N/A";
            
            // Line 1: SKU and Unit Weight
            doc.text(`SKU: ${sku}   Unit Weight: ${currentItem.weight_kg} kg`, margin, infoY + 6);
            
            // Line 2: Position, Dimensions, Rotation
            const infoLine = `pos(mm): (${placement.pos_x.toFixed(0)}, ${placement.pos_y.toFixed(0)}, ${placement.pos_z.toFixed(0)})   dims(mm): ${dims.length_mm}×${dims.width_mm}×${dims.height_mm}   rot: ${placement.rotation}`;
            doc.text(infoLine, margin, infoY + 13);
        } else {
            doc.text(
                `pos(mm): (${placement.pos_x.toFixed(0)}, ${placement.pos_y.toFixed(0)}, ${placement.pos_z.toFixed(0)})   rot: ${placement.rotation}`,
                margin,
                infoY + 8
            );
        }

        // Highlight swatch
        if (currentItem?.color_hex) {
            const swatchSize = 6;
            const swatchX = pageWidth - margin - swatchSize;
            doc.setDrawColor(80, 80, 80);
            const rgb = this.hexToRgb(currentItem.color_hex);
            if (rgb) {
                doc.setFillColor(rgb.r, rgb.g, rgb.b);
                doc.rect(swatchX, infoY - 5, swatchSize, swatchSize, "FD");
            } else {
                doc.rect(swatchX, infoY - 5, swatchSize, swatchSize, "S");
            }
        }
    }

    private getDimsForRotation(item: ItemData, rotation: number): { length_mm: number; width_mm: number; height_mm: number } {
        // rotation codes (0..5) come from our backend API and represent a
        // permutation of (length,width,height) in container coordinates.
        const l = item.length_mm;
        const w = item.width_mm;
        const h = item.height_mm;

        switch (rotation) {
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

    private drawItemRect(
        doc: any,
        rect: { x: number; y: number; w: number; h: number },
        colorHex: string,
        isCurrent: boolean
    ) {
        if (rect.w <= 0 || rect.h <= 0) return;

        if (isCurrent) {
            const rgb = this.hexToRgb(colorHex);
            if (rgb) {
                doc.setFillColor(rgb.r, rgb.g, rgb.b);
                doc.setDrawColor(20, 20, 20);
                doc.setLineWidth(0.6);
                doc.rect(rect.x, rect.y, rect.w, rect.h, "FD");
                return;
            }

            doc.setDrawColor(20, 20, 20);
            doc.setLineWidth(0.6);
            doc.rect(rect.x, rect.y, rect.w, rect.h, "S");
            return;
        }

        // Prior placements: outline only
        doc.setDrawColor(190, 190, 190);
        doc.setLineWidth(0.2);
        doc.rect(rect.x, rect.y, rect.w, rect.h, "S");
    }

    private drawAxes(
        doc: any,
        opts: {
            originX: number;
            originY: number;
            xLen: number;
            yLen: number;
            xLabel: string;
            yLabel: string;
            xMaxLabel: string;
            yMaxLabel: string;
        }
    ) {
        const { originX, originY, xLen, yLen, xLabel, yLabel, xMaxLabel, yMaxLabel } = opts;

        const arrow = 2.5;
        doc.setDrawColor(120, 120, 120);
        doc.setLineWidth(0.3);

        // X axis (to the right)
        doc.line(originX, originY, originX + xLen, originY);
        doc.line(originX + xLen, originY, originX + xLen - arrow, originY - arrow / 2);
        doc.line(originX + xLen, originY, originX + xLen - arrow, originY + arrow / 2);

        // Y axis (upwards)
        doc.line(originX, originY, originX, originY - yLen);
        doc.line(originX, originY - yLen, originX - arrow / 2, originY - yLen + arrow);
        doc.line(originX, originY - yLen, originX + arrow / 2, originY - yLen + arrow);

        doc.setTextColor(120, 120, 120);
        doc.setFontSize(8);
        doc.setFont(undefined, "normal");

        // Labels inside the container box
        const inset = 2;

        // Origin label (bottom-left)
        doc.text("0", originX + inset, originY - inset);

        // X axis labels (bottom-right)
        doc.text(xLabel, originX + xLen - 14, originY - inset);
        doc.text(xMaxLabel, originX + xLen - 8, originY - inset - 3);

        // Y axis labels (top-left)
        doc.text(yLabel, originX + inset, originY - yLen + 6);
        doc.text(yMaxLabel, originX + inset, originY - yLen + 12);
    }

    private hexToRgb(hex: string): { r: number; g: number; b: number } | null {
        const normalized = hex.trim();
        const match = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(normalized);
        if (!match) return null;
        return {
            r: parseInt(match[1], 16),
            g: parseInt(match[2], 16),
            b: parseInt(match[3], 16),
        };
    }
}
