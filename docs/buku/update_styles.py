from docx import Document
from docx.shared import RGBColor
from docx.oxml.ns import qn
from docx.oxml import OxmlElement

def set_style_shading(doc, style_name, hex_color):
    """
    Sets the background shading for a specific paragraph style.
    hex_color: str, e.g., "F4F4F4" (no #)
    """
    styles = doc.styles
    try:
        style = styles[style_name]
    except KeyError:
        print(f"Style '{style_name}' not found. Creating it...")
        style = styles.add_style(style_name, 1) # 1 = PARAGRAPH
        
    # Access the XML element of the style
    style_elm = style.element
    
    # Check if pPr exists, create if not
    pPr = style_elm.get_or_add_pPr()
    
    # Check if shd exists, create/update it
    shd = pPr.find(qn('w:shd'))
    if shd is None:
        shd = OxmlElement('w:shd')
        pPr.append(shd)
        
    # Set the attributes for shading
    shd.set(qn('w:val'), 'clear')
    shd.set(qn('w:color'), 'auto')
    shd.set(qn('w:fill'), hex_color)
    
    # Update Font settings (Consolas, 11pt, No Bold)
    rPr = style_elm.get_or_add_rPr()
    
    # Fonts
    rFonts = rPr.get_or_add_rFonts()
    rFonts.set(qn('w:ascii'), 'Courier New')
    rFonts.set(qn('w:hAnsi'), 'Courier New')
    
    # Size 10pt (20 half-points)
    sz = rPr.get_or_add_sz()
    sz.set(qn('w:val'), '20')
    
    # Force No Bold
    b = rPr.get_or_add_b()
    b.set(qn('w:val'), '0')
    
    # Add Padding using Borders (w:pBdr) + Indentation (w:ind)
    # 1. Indent text IN by 12pt (240 twips)
    ind = pPr.find(qn('w:ind'))
    if ind is None:
        ind = OxmlElement('w:ind')
        pPr.append(ind)
    ind.set(qn('w:left'), '240')  # 240 twips = 12pt
    ind.set(qn('w:right'), '240')
    
    # 2. Add Border OUT by 12pt
    pBdr = pPr.find(qn('w:pBdr'))
    if pBdr is None:
        pBdr = OxmlElement('w:pBdr')
        pPr.append(pBdr)
    
    for border_name in ['top', 'left', 'bottom', 'right']:
        border = pBdr.find(qn(f'w:{border_name}'))
        if border is None:
            border = OxmlElement(f'w:{border_name}')
            pBdr.append(border)
            
        border.set(qn('w:val'), 'single')
        border.set(qn('w:sz'), '4')       # 1/2 point width
        border.set(qn('w:space'), '12')   # 12 points spacing (padding matches indent)
        border.set(qn('w:color'), hex_color) # Same as background
    
    print(f"Updated style '{style_name}': Font=Courier New 10pt, Shading=#{hex_color}, Padding=Fixed (Indent+Border)")

def main():
    template_path = "Template Buku Penerbit FSM 2025.docx"
    
    try:
        doc = Document(template_path)
        
        # Update "Source Code" style
        set_style_shading(doc, "Source Code", "F4F4F4") 
        
        doc.save(template_path)
        print("Successfully updated template styles.")
        
    except Exception as e:
        print(f"Error updating template: {e}")

if __name__ == "__main__":
    main()
