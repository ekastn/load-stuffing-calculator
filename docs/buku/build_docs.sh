#!/bin/bash
set -e

# Configuration
TEMPLATE_DOC="Template Buku Penerbit FSM 2025.doc"
TEMPLATE_DOCX="Template Buku Penerbit FSM 2025.docx"
SOURCE_MD="draft.md"
OUTPUT_DOCX="Draft_Layouted.docx"

echo "Checking environment..."

# 1. Ensure Template is in .docx format (required for pandoc reference-doc)
if [ ! -f "$TEMPLATE_DOCX" ]; then
    if [ -f "$TEMPLATE_DOC" ]; then
        echo "Converting legacy .doc template to .docx..."
        # LibreOffice headless conversion
        libreoffice --headless --convert-to docx "$TEMPLATE_DOC"
    else
        echo "Error: Template file '$TEMPLATE_DOC' not found!"
        echo "Please place the publisher template in this directory."
        exit 1
    fi
fi

# 2. Build the document
if [ -f "$SOURCE_MD" ]; then
    echo "Building Word document from $SOURCE_MD..."
    pandoc "$SOURCE_MD" \
        -o "$OUTPUT_DOCX" \
        --lua-filter="mermaid-filter.lua" \
        --highlight-style=pygments \
        --reference-doc="$TEMPLATE_DOCX" \
        --toc
    
    echo "Success! Created '$OUTPUT_DOCX'"
else
    echo "Error: Source file '$SOURCE_MD' not found!"
    exit 1
fi
