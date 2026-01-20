import jsPDF from "jspdf";
import autoTable from "jspdf-autotable";
import { Container, PlanDetail, PlanItem, Placement } from "@/lib/api/types";

// Helper to generate consistent colors from strings (product IDs)
function stringToColor(str: string): string {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    const c = (hash & 0x00ffffff).toString(16).toUpperCase();
    return "#" + "00000".substring(0, 6 - c.length) + c;
}

function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
    const normalized = hex.trim();
    const match = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(normalized);
    if (!match) return null;
    return {
        r: parseInt(match[1], 16),
        g: parseInt(match[2], 16),
        b: parseInt(match[3], 16),
    };
}

export async function generateStuffingReport(plan: PlanDetail, container: Container) {
    const doc = new jsPDF({
        orientation: "landscape",
        unit: "mm",
        format: "a4",
    });

    const pageWidth = doc.internal.pageSize.getWidth();
    const pageHeight = doc.internal.pageSize.getHeight();
    const margin = 20;

    // --- Page 1: Summary ---
    createSummaryPage(doc, plan, container, pageWidth, pageHeight, margin);

    // --- Page 2+: Steps ---
    if (plan.placements && plan.placements.length > 0) {
        createStepPages(doc, plan, container, pageWidth, pageHeight, margin);
    }

    doc.save(`manifest-${plan.id.substring(0, 8)}.pdf`);
}

function createSummaryPage(
    doc: jsPDF,
    plan: PlanDetail,
    container: Container,
    pageWidth: number,
    pageHeight: number,
    margin: number
) {
    let yPosition = margin;

    // Title
    doc.setTextColor(0, 64, 255);
    doc.setFontSize(28);
    doc.setFont("helvetica", "bold");
    doc.text("LOADING PLAN", margin, yPosition);
    yPosition += 15;

    // Container Info
    doc.setTextColor(0, 0, 0);
    doc.setFontSize(14);
    doc.text(`Container: ${container.name}`, margin, yPosition);
    doc.setFontSize(10);
    doc.setFont("helvetica", "normal");
    doc.text(`Dims: ${container.length_mm}x${container.width_mm}x${container.height_mm} mm`, margin, yPosition + 6);
    doc.text(`Plan ID: ${plan.id}`, margin, yPosition + 12);
    
    yPosition += 25;

    // Items Table
    const tableBody = plan.items.map(item => [
        item.label,
        `${item.length_mm} x ${item.width_mm} x ${item.height_mm}`,
        item.weight_kg.toString(),
        item.quantity.toString()
    ]);

    autoTable(doc, {
        startY: yPosition,
        head: [['Product Name', 'Dimensions (mm)', 'Weight (kg)', 'Quantity']],
        body: tableBody,
        headStyles: { fillColor: [66, 66, 66] },
    });
}

function createStepPages(
    doc: jsPDF,
    plan: PlanDetail,
    container: Container,
    pageWidth: number,
    pageHeight: number,
    margin: number
) {
    const placements = plan.placements || [];
    // Sort by step number
    const sortedPlacements = [...placements].sort((a, b) => a.step_number - b.step_number);

    // Map product details for easy lookup
    const productMap = new Map<string, PlanItem>();
    plan.items.forEach(item => productMap.set(item.product_id, item));

    for (let i = 0; i < sortedPlacements.length; i++) {
        doc.addPage();
        renderStepLandscape(doc, sortedPlacements[i], sortedPlacements, productMap, container, pageWidth, pageHeight, margin);
        
        // Footer
        doc.setTextColor(120, 120, 120);
        doc.setFontSize(8);
        doc.setFont("helvetica", "italic");
        doc.text(`Page ${i + 2} - Step ${sortedPlacements[i].step_number}`, margin, pageHeight - 10);
    }
}

function renderStepLandscape(
    doc: jsPDF,
    currentPlacement: Placement,
    allPlacements: Placement[],
    productMap: Map<string, PlanItem>,
    container: Container,
    pageWidth: number,
    pageHeight: number,
    margin: number
) {
    let yPosition = margin + 5;

    // Title
    doc.setTextColor(0, 64, 255);
    doc.setFontSize(20);
    doc.setFont("helvetica", "bold");
    const stepText = `Step ${currentPlacement.step_number}: Load ${currentPlacement.label}`;
    const textWidth = doc.getTextWidth(stepText);
    doc.text(stepText, (pageWidth - textWidth) / 2, yPosition);
    yPosition += 10;

    // Layout Calculation
    const gap = 10;
    const legendHeight = 30;
    const diagramTop = yPosition;
    const diagramHeight = pageHeight - margin - legendHeight - diagramTop;
    const diagramWidth = pageWidth - 2 * margin;
    const diagramBoxWidth = (diagramWidth - gap) / 2;

    const topViewBox = { x: margin, y: diagramTop, w: diagramBoxWidth, h: diagramHeight };
    const sideViewBox = { x: margin + diagramBoxWidth + gap, y: diagramTop, w: diagramBoxWidth, h: diagramHeight };

    // Function to calculate scale to fit container in box
    const fitScale = (boxW: number, boxH: number, contL: number, contD: number) => {
        const padding = 20;
        const usefulW = boxW - padding;
        const usefulH = boxH - padding;
        const scale = Math.min(usefulW / contL, usefulH / contD);
        return {
            scale,
            offsetX: (boxW - (contL * scale)) / 2,
            offsetY: (boxH - (contD * scale)) / 2
        };
    };

    const topFit = fitScale(topViewBox.w, topViewBox.h, container.length_mm, container.width_mm);
    const sideFit = fitScale(sideViewBox.w, sideViewBox.h, container.length_mm, container.height_mm);

    const topOriginX = topViewBox.x + topFit.offsetX;
    const topOriginY = topViewBox.y + topFit.offsetY;
    
    // Side view origin
    const sideOriginX = sideViewBox.x + sideFit.offsetX;
    const sideOriginY = sideViewBox.y + sideFit.offsetY;

    // Draw Container Outlines
    doc.setDrawColor(0, 0, 0);
    doc.setLineWidth(0.5);
    
    // Top View Container
    doc.rect(topOriginX, topOriginY, container.length_mm * topFit.scale, container.width_mm * topFit.scale);
    doc.setFontSize(10);
    doc.text("Top View (XY)", topOriginX, topOriginY - 5);

    // Side View Container
    doc.rect(sideOriginX, sideOriginY, container.length_mm * sideFit.scale, container.height_mm * sideFit.scale);
    doc.text("Side View (XZ)", sideOriginX, sideOriginY - 5);

    // Draw Placements (up to current step)
    const visiblePlacements = allPlacements
        .filter(p => p.step_number <= currentPlacement.step_number)
        .sort((a, b) => a.step_number - b.step_number);

    for (const p of visiblePlacements) {
        const isCurrent = p.step_number === currentPlacement.step_number;
        const item = productMap.get(p.product_id);
        if (!item) continue;

        const getDims = (rot: number) => {
             const l = item.length_mm, w = item.width_mm, h = item.height_mm;
             switch(rot) {
                case 1: return { l: w, w: l, h: h };
                case 2: return { l: w, w: h, h: l };
                default: return { l, w, h };
             }
        };
        const dim = getDims(p.rotation);

        const colorHex = stringToColor(p.product_id);
        const rgb = hexToRgb(colorHex);

        if (isCurrent && rgb) {
            doc.setFillColor(rgb.r, rgb.g, rgb.b);
            doc.setDrawColor(20, 20, 20); // Bold outline for current
        } else {
            doc.setFillColor(240, 240, 240); // Grey for past
            doc.setDrawColor(180, 180, 180);
        }
        
        // Top View Logic:
        const topX = topOriginX + (p.pos_x * topFit.scale);
        const topY = topOriginY + ((container.width_mm - (p.pos_y + dim.w)) * topFit.scale);
        
        doc.rect(topX, topY, dim.l * topFit.scale, dim.w * topFit.scale, isCurrent ? "FD" : "DF");

        // Side View Logic:
        const sideX = sideOriginX + (p.pos_x * sideFit.scale);
        const sideY = sideOriginY + ((container.height_mm - (p.pos_z + dim.h)) * sideFit.scale);

        doc.rect(sideX, sideY, dim.l * sideFit.scale, dim.h * sideFit.scale, isCurrent ? "FD" : "DF");
    }

    // Legend for current step
    const currentItem = productMap.get(currentPlacement.product_id);
    if (currentItem) {
        doc.setTextColor(0, 0, 0);
        doc.setFontSize(12);
        doc.text(`Item: ${currentItem.label} (ID: ${currentItem.product_id.substring(0,6)})`, margin, pageHeight - margin - 20);
        doc.text(`Position: (${currentPlacement.pos_x}, ${currentPlacement.pos_y}, ${currentPlacement.pos_z})`, margin, pageHeight - margin - 15);
    }
}
