from docx import Document
from docx.shared import Inches, Mm

def get_layout_info(docx_path):
    doc = Document(docx_path)
    
    # Usually the first section defines the main layout
    section = doc.sections[0]
    
    print(f"--- Layout Info for {docx_path} ---")
    
    # Page Size
    width_mm = section.page_width.mm
    height_mm = section.page_height.mm
    print(f"Page Size: {width_mm:.1f} x {height_mm:.1f} mm")
    
    # Margins
    print(f"Margins:")
    print(f"  Top:    {section.top_margin.mm:.1f} mm")
    print(f"  Bottom: {section.bottom_margin.mm:.1f} mm")
    print(f"  Left:   {section.left_margin.mm:.1f} mm")
    print(f"  Right:  {section.right_margin.mm:.1f} mm")
    print(f"  Gutter: {section.gutter.mm:.1f} mm")
    
    # Identify standard size
    if abs(width_mm - 210) < 2 and abs(height_mm - 297) < 2:
        print("Standard: A4")
    elif abs(width_mm - 215.9) < 2 and abs(height_mm - 279.4) < 2:
        print("Standard: Letter")
    elif abs(width_mm - 148) < 2 and abs(height_mm - 210) < 2:
        print("Standard: A5")
    elif abs(width_mm - 176) < 2 and abs(height_mm - 250) < 2:
        print("Standard: B5 (ISO)")
    elif abs(width_mm - 182) < 2 and abs(height_mm - 257) < 2:
        print("Standard: B5 (JIS)")
    else:
        print("Standard: Custom / Unknown")

if __name__ == "__main__":
    get_layout_info("Template Buku Penerbit FSM 2025.docx")
