import type { StuffingPlanData, ItemData } from "../types";

export interface ReportGeneratorConfig {
    brandLogoUrl?: string;
    companyName?: string;
    stepImages?: string[]; // Base64 images for each step
    stepAspectRatios?: number[]; // Aspect ratios for each step image
    fullScreenshot?: string; // Full screenshot with all items for summary page
    fullScreenshotAspectRatio?: number; // Aspect ratio for the full screenshot
}

export class ReportGenerator {
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
        const colPackages = rightX + (availableWidth * 0.40);
        const colVolume = rightX + (availableWidth * 0.55);
        const colWeight = rightX + (availableWidth * 0.75);

        doc.setFontSize(10);
        doc.setFont(undefined, "bold");
        doc.text("Name", rightX, rightY);
        doc.text("Packages", colPackages, rightY);
        doc.text("Volume", colVolume, rightY);
        doc.text("Weight", colWeight, rightY);
        rightY += 1;

        doc.line(rightX, rightY, containerX + containerWidth - 20, rightY);
        rightY += 6;

        // Items data with better formatting
        doc.setFontSize(9);
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

            doc.text(item.label, rightX, rightY);
            doc.text(count.toString(), colPackages + 5, rightY);
            doc.text(`${totalVolume.toFixed(1)} m3`, colVolume, rightY);
            doc.text(`${totalWeight.toFixed(1)} kg`, colWeight, rightY);
            rightY += 6;
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
        const itemMap = new Map<string, ItemData>();
        data.items.forEach((item) => itemMap.set(item.item_id, item));

        for (let i = 0; i < sortedPlacements.length; i++) {
            doc.addPage();
            this.renderStepLandscape(
                doc,
                i,
                sortedPlacements[i],
                data.items,
                config,
                pageWidth,
                pageHeight,
                margin
            );
        }

        // Footer on last page
        doc.setTextColor(120, 120, 120);
        doc.setFontSize(8);
        doc.setFont(undefined, "italic");
        doc.text(`Generated on ${new Date().toLocaleString()}`, margin, pageHeight - 10);
    }

    private renderStepLandscape(
        doc: any,
        index: number,
        placement: any,
        items: ItemData[],
        config: ReportGeneratorConfig,
        pageWidth: number,
        pageHeight: number,
        margin: number
    ) {
        let yPosition = margin + 5;

        doc.setTextColor(...this.BRAND_BLUE);
        doc.setFontSize(20);
        doc.setFont(undefined, "bold");
        const stepText = `Step ${placement.step_number}`;
        const textWidth = doc.getTextWidth(stepText);
        doc.text(stepText, (pageWidth - textWidth) / 2, yPosition);
        yPosition += 12;

        if (config.stepImages?.[index]) {
            try {
                const maxImgWidth = pageWidth - 2 * margin - 40;
                const maxImgHeight = pageHeight / 2;
                const aspectRatio = config.stepAspectRatios?.[index] || 1;

                let imgWidth, imgHeight;

                if ((maxImgWidth / maxImgHeight) > aspectRatio) {
                    imgHeight = maxImgHeight;
                    imgWidth = maxImgHeight * aspectRatio;
                } else {
                    imgWidth = maxImgWidth;
                    imgHeight = imgWidth / aspectRatio;
                }
                
                const imgX = (pageWidth - imgWidth) / 2;

                doc.addImage(config.stepImages[index], "PNG", imgX, yPosition, imgWidth, imgHeight);
                yPosition += imgHeight + 10;
            } catch (error) {
                console.error("Failed to add step image:", error);
                yPosition += 10;
            }
        }

        doc.setTextColor(0, 0, 0);
        doc.setFontSize(12);
        doc.setFont(undefined, "normal");

        const currentItem = items.find(item => item.item_id === placement.item_id);

        if (currentItem) {
            // Calculate width of the label text
            const labelTextWidth = doc.getTextWidth(currentItem.label);
            
            // Define dimensions for the color square
            const squareSize = 5; // 5mm
            const paddingBetweenSquareAndText = 3; // 3mm

            // Calculate total width of the block (square + padding + text)
            const totalBlockWidth = squareSize + paddingBetweenSquareAndText + labelTextWidth;

            // Calculate new itemInfoX to center the entire block
            const itemInfoX = (pageWidth - totalBlockWidth) / 2;
            
            // Draw color square
            doc.setFillColor(currentItem.color_hex);
            doc.rect(itemInfoX, yPosition - 4, squareSize, squareSize, "F");
            doc.setDrawColor(100, 100, 100);
            doc.rect(itemInfoX, yPosition - 4, squareSize, squareSize, "S");

            // Draw label
            doc.text(currentItem.label, itemInfoX + squareSize + paddingBetweenSquareAndText, yPosition);
        }
    }
}
