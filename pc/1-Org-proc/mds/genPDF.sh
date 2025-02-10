#!/bin/bash

# Comprobar que se han pasado los argumentos correctos
if [ "$#" -ne 2 ]; then
  echo "Uso: $0 <archivo_markdown_entrada> <archivo_pdf_salida>"
  exit 1
fi

# Variables de configuración
INPUT_FILE="$1"                # Archivo de entrada Markdown pasado como argumento
OUTPUT_FILE="$2"               # Archivo de salida PDF pasado como argumento
PAPERSIZE="A4"                 # Tamaño de papel
GEOMETRY="top=2cm, bottom=1.5cm, left=2cm, right=2cm" # Márgenes
TOC_DEPTH=3                      # Profundidad del índice

# Comando pandoc
pandoc \
  -f markdown \
  -V papersize="$PAPERSIZE" \
  -V geometry="$GEOMETRY" \
  -V toc-title="Índice" \
  --toc \
  --toc-depth="$TOC_DEPTH" \
  --pdf-engine=pdflatex \
  -V fontsize=12pt \
  -V linkcolor=purple \
  -V subparagraph=false \
  "$INPUT_FILE.md" -o "../pdfs/$OUTPUT_FILE.pdf"

# Mensaje de éxito
echo "PDF generado exitosamente: ../pdfs/$OUTPUT_FILE"
