#!/bin/bash

# =============================
# Script para convertir Markdown → PDF con Pandoc
# =============================

# Configuración por defecto
INCLUDE_TOC=false
NUMBER_SECTIONS=false

mostrar_uso() {
  echo "Uso: $0 [opciones] <nombre_base_archivo>"
  echo "Opciones:"
  echo "  -i             Inicializa el directorio (crea carpetas mds, pdfs, resources)"
  echo "  -u <nombre>    Convierte <nombre>.md a <nombre>.pdf"
  echo "  -t             Incluye Tabla de Contenidos (Índice)"
  echo "  -n             Numera los encabezados (1. Introducción, 1.1...)"
  echo "  -h             Muestra esta ayuda"
  exit 1
}

# Procesar opciones
# Nota: "u:" indica que -u requiere un argumento
while getopts "i u: t n h" opt; do
  case $opt in
  i)
    mkdir -p mds pdfs resources
    echo "✅ Directorios creados."
    exit 0
    ;;
  u)
    BASENAME="$OPTARG"
    ;;
  t)
    INCLUDE_TOC=true
    ;;
  n)
    NUMBER_SECTIONS=true
    ;;
  h | *)
    mostrar_uso
    ;;
  esac
done

# Desplazar los argumentos procesados por getopts
shift $((OPTIND - 1))

# Si no se usó -u, intentar tomar el primer argumento libre como nombre
if [ -z "$BASENAME" ]; then
  if [ -n "$1" ]; then
    BASENAME="$1"
  else
    mostrar_uso
  fi
fi

INPUT_FILE="${BASENAME}.md"
OUTPUT_FILE="${BASENAME}.pdf"
INPUT_PATH="mds/$INPUT_FILE"
OUTPUT_PATH="pdfs/$OUTPUT_FILE"

# Validar que el archivo existe
if [ ! -f "$INPUT_PATH" ]; then
  echo "❌ Error: No se encuentra el archivo $INPUT_PATH"
  exit 1
fi

# =============================
# Configuración Dinámica de Pandoc
# =============================

# Configuración de Índice
TOC_OPTS=""
if [ "$INCLUDE_TOC" = true ]; then
  TOC_OPTS="--toc --toc-depth=3 -V toc-title=Índice"
  echo "→ Opción: Con Índice"
fi

# Configuración de Numeración
NUM_OPTS=""
if [ "$NUMBER_SECTIONS" = true ]; then
  NUM_OPTS="--number-sections"
  echo "→ Opción: Con Numeración"
fi

# Detección automática de bibliografía
BIBLIO_OPTS=""
for EXT in bib json yaml yml; do
  BIB_FILE="mds/${BASENAME}.${EXT}"
  if [ -f "$BIB_FILE" ]; then
    echo "→ Bibliografía detectada: $BIB_FILE"
    BIBLIO_OPTS="--citeproc --bibliography=$BIB_FILE"
    break
  fi
done

# =============================
# Ejecución de Pandoc
# =============================
pandoc \
  -f markdown \
  -V papersize="A4" \
  -V geometry="top=2cm, bottom=1.5cm, left=2cm, right=2cm" \
  -V colorlinks=true \
  -V urlcolor=blue \
  --pdf-engine=pdflatex \
  $TOC_OPTS \
  $NUM_OPTS \
  $BIBLIO_OPTS \
  "$INPUT_PATH" -o "$OUTPUT_PATH"

# =============================
# Comprobar resultado
# =============================
if [ $? -eq 0 ]; then
  echo "✅ Conversión completada: $OUTPUT_PATH"
  read -p "¿Desea abrir el archivo? (Y/N): " respuesta
  case "$respuesta" in
  [Yy]*) open "$OUTPUT_PATH" 2>/dev/null || xdg-open "$OUTPUT_PATH" 2>/dev/null ;;
  *) ;;
  esac
else
  echo "❌ Error en la conversión"
  exit 1
fi
