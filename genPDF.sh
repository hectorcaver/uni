#!/bin/bash

# =============================
# Script para convertir Markdown → PDF con Pandoc
# Añade automáticamente bibliografía si existe
# =============================

if [ "$#" -lt 1 ]; then
  echo "Uso: $0 [-i] [-u] <archivo_md_entrada> [ <archivo_pdf_salida> ]"
  echo "Opciones:"
  echo "-u <arvhivo_md_entrada> ($0 -u ej.md --> genera el fichero ej.pdf)"
  echo "-i --> inicializa el directorio, con las carpetas necesarias"
  exit 1
fi

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
  if [ "$1" == "-i" ]; then
    mkdir mds
    mkdir pdfs
    mkdir resources
    exit 0
  fi
  if [ "$#" -ne 2 ]; then
    echo "Uso: $0 <archivo_markdown_entrada> <archivo_pdf_salida>"
    exit 1
  fi
  INPUT_FILE="$1.md"
  OUTPUT_FILE="$2.pdf"
  BASENAME="$1"
fi

# =============================
# Configuración
# =============================
PAPERSIZE="A4"                                        # Tamaño de papel
GEOMETRY="top=2cm, bottom=1.5cm, left=2cm, right=2cm" # Márgenes
TOC_DEPTH=3                                           # Profundidad del índice
INPUT_PATH="mds/$INPUT_FILE"
OUTPUT_PATH="pdfs/$OUTPUT_FILE"

# =============================
# Detección automática de bibliografía
# =============================
BIBLIO_OPTS=""

# Buscar bibliografía en mds/ con el mismo nombre base
for EXT in bib json yaml yml; do
  BIB_FILE="mds/${BASENAME}.${EXT}"
  if [ -f "$BIB_FILE" ]; then
    echo "→ Se ha detectado bibliografía: $BIB_FILE"
    BIBLIO_OPTS="--citeproc --bibliography=$BIB_FILE"
    break
  fi
done

# =============================
# Conversión con Pandoc
# =============================
pandoc -f markdown+smart \
  --standalone \
  --filter=pandoc-crossref \
  -V papersize="$PAPERSIZE" \
  -V geometry="$GEOMETRY" \
  -V toc-title="Índice" \
  -V colorlinks=true \
  -V urlcolor=blue \
  --toc \
  --toc-depth=$TOC_DEPTH \
  --pdf-engine=pdflatex \
  --number-sections \
  --syntax-highlighting=idiomatic \
  --from=markdown \
  --to=latex \
  $BIBLIO_OPTS \
  "$INPUT_PATH" -o "$OUTPUT_PATH"

# =============================
# Comprobar resultado
# =============================
if [ $? -eq 0 ]; then
  echo "✅ Conversión completada: $OUTPUT_FILE"
  read -p "¿Desea abrir el archivo? (Y/N): " respuesta
  case "$respuesta" in
  [Yy]*) open "$OUTPUT_PATH" ;;
  *) ;;
  esac
else
  echo "❌ Error en la conversión"
  exit 1
fi
