# Multiprocesadores

## Práctica 1: evaluación

Para calcular nº instrucciones en ejecución (en código ensamblador):

- esc.avx si que tiene LEN iteraciones
- vec.* tienen menos iteraciones

Dadas:

- $ESCALAR = 100I$
- $VEC = 25I$

El factor de reducción es: $F_R = \frac{ESCALAR}{VEC}$

La reducción es: $R = \frac{ESCALAR - VEC}{100}$

Para calcular el Speed-Up:

- $SpeedUp = \frac{T_{esc}}{T_v}$


Para calcular R:

- $R = \frac{N_{FLOP}}{T_{ns}} GFLOPS$


**Importante**:

- No contar instrucciones de memoria
- Una instrucción vectorial hace N operaciones de cálculo
- En el caso de FMA hará el doble de operaciones, suma y multiplicación conjuntas.
- Velocidad pico: velocidad en un estado ideal (todo operaciones de cálculo).
  - Calcular el límite de la velocidad.

1. Me meto en el LAB, ejecuto el comando para ver máquinas abiertas.
2. Hago ssh a una máquina encendida.
3. Buscas en internet el procesador de la máquina
4. Buscas en la documentación oficial la familia del procesador (Coffee Lake)
5. Buscar número de unidades funcionales y su frecuencia.
6. Frecuencia turbo máxima (buscamos la mayor).
7. Buscar en ``Wikichip coffee lake``.
8. Buscar ruta de datos.
9. Buscar cuantas UFs pueden ejecutar una instrucción de cálculo cada ciclo.
10. Hay que calcular tres picos distintos:
    1.  pico escalar
    2.  pico vectorial sin fma
    3.  pico vectorial con fma

