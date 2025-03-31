% Tema 6: Aprendizaje Automático
  Aprendizaje no supervisado.
  Reducción de la dimensionalidad.
% Autor: Héctor Lacueva Sacristán 869637.
% Marzo-2025

# Aprendizaje supervisado vs. Aprendizaje NO supervisado

## Aprendizaje supervisado

- Entradas: (atributos de entrada y salida real)
    $$D = \{x^{(i)}, y^{(i)}\}^m_{i=1}$$
  - En Regresión: x e y pertenecen a los reales.
  - En Clasificación: x pertenece a los reales e y pertenece a una clase.


## Aprendizaje no supervisado

- Objetivo: encontrar información sobre la estructura de los datos.
- Las etiquetas (y) no están disponibles o son muy caras de obtener.
- Entradas: (atributos de entrada y ya no contamos con la salida real)
    $$D = \{x^{(i)}\}^m_{i=1}$$
  - x pertenece a los reales.

Se suele emplear para:

- Comprensión y reducción de la dimensión, simplificar problemas.
- Visualización de datos de alta dimension.
- Descubrir conocimiento (agrupamiento o clustering).

# Comprensión de los datos

En la imagen a continuación, **¿Es posible usar estos datos de menor dimensión para obtener buenos resultados en la clasificación?**.

**La resupuesta es SI**, se pueden distinguir distintos grupos para cada clase.

## Ventaja

Se consigue un modelo menos costoso en tiempo y memoria. Hemos pasado de 400 atributos a solamente 2.
**Por lo que computacionalmente es 200 veces menos costoso**.

## Desventaja

Se pierde información, aunque a veces es tan poca que merece la pena.
Por ejemplo, las esquinas de las imágenes casi nunca tienen información relevante y se podrían desconsiderar sin afectar mucho a la información.


# Correlación de atributos

Dependiendo de la correlación podemos definir que atributos son (o podemos considerar) independientes.

| Correlación negativa (dependencia )| Correlación 0 (independencia) | Correlación positiva (dependencia) |
|:-:|:-:|:-:|
|img|img|img|
|Valores pequeños de un atributo producen valores pequeños del segundo atributo|-|Valores grandes de un atributo producen valores grandes del segundo atributo|