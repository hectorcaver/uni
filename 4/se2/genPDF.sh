#!/bin/bash

# Comprobar si se ha pasado la opción -u
if [ "$1" == "-u" ]; then
  if [ "$#" -ne 2 ]; then
    echo "Uso: $0 -u <nombre_archivo_sin_extension>"
    exit 1
  fi
  BASENAME="$2"
  INPUT_FILE="${BASENAME}.md"
  OUTPUT_FILE="${BASENAME}.pdf"
else
  if [ "$#" -ne 2 ]; then
    echo "Uso: $0 <archivo_markdown_entrada> <archivo_pdf_salida>"
    exit 1
  fi
  INPUT_FILE="$1.md"
  OUTPUT_FILE="$2.pdf"
fi

# Variables de configuración
PAPERSIZE="A4"                                        # Tamaño de papel
GEOMETRY="top=2cm, bottom=1.5cm, left=2cm, right=2cm" # Márgenes
TOC_DEPTH=3                                           # Profundidad del índice

# Convertir el archivo Markdown a PDF usando pandoc
pandoc \
  -f markdown \
  -V papersize="$PAPERSIZE" \
  -V geometry="$GEOMETRY" \
  -V toc-title="Índice" \
  --toc \
  --toc-depth=$TOC_DEPTH \
  --pdf-engine=pdflatex \
  "mds/$INPUT_FILE" -o "pdfs/$OUTPUT_FILE"

# Verificar si la conversión fue exitosa
if [ $? -eq 0 ]; then
  echo "Conversión completada: $OUTPUT_FILE"

  read -p "¿Desea abrir el archivo? (Y/N): " respuesta
  case "$respuesta" in
  [Yy]*) open "pdfs/$OUTPUT_FILE" ;;
  *) ;;
  esac
else
  echo "Error en la conversión"
  exit 1
fi
