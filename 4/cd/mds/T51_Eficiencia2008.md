% Ejercicios Tema 5  
  Centros de Datos UZ, 2025-26  
% Lacueva Sacristán, Héctor
% 26/10/2025 


# T5.1 Eficiencia 2008 (Asus) (Tiempo total 1h)

IA: No.

![Fig. 1: Gráfica sacada de [@specasus2008] sobre el procesador Intel Xeon L5430 en el ASUS RS160-E5 (2008), el eje X representa la carga del sistema, desde "activo en espera" (izquierda), hasta 100\% (derecha). El eje Y (izquierda) representa la métrica R (ssj_ops/watt). El eje Y (derecha) muestra la métrica $P_{IT}$ (W).](resources/T51_asus2008.png)

## Compara los siguientes tres modos de operación para una carga media:

### Trabajar con tres servidores al 20\%

Cada uno de los servidores consume alrededor de 113 W. Por lo tanto en este modo de operación se consume un total de **339 W**.

En cuanto a ssj_ops/watt contamos con 501 por cada servidor. Con los tres sigue siendo la misma cantidad.

### Un servidor al 60\% con los otros dos apagados

Un servidor al 60% consume alrededor de 150 W y tiene un rendimiento de 1120 ssj_ops/watt. Los servidores apagados ni gastan ni suman rendimiento.

### Un servidor al 60\% con uno apagado y el otro idle

El serviodr al 60% (150 W y 1120 ssj_ops/watt, 168.211 ssj_ops totales), el apagado (0 W y 0 ssj_ops/watt) y el que está en idle (89,4 W y 0 ssj_ops).

Esto hace un total de 239,4 W y $\frac{168.211 \text{\ ssj\_ops}}{239,4 W} = 702.63$ ssj_ops/watt.

### Comparación

|Modo opreación|$P_{IT}$ (W)| ssj_ops/watt |
|:-:|:-:|:-:|
|Tres servidores al 20%|339 W|501 ssj_ops/watt|
|Un servidor al 60% y dos apagados|150 W|1120 ssj_ops/watt|
|Un servidor al 60%, otro apagado y otro idle|239,4 W|702,63 ssj_ops/watt|

El modo más eficiente por vatio es claramente un servidor al 60% con los otros apagados (1120 ssj_ops/W).


# T5.2 Eficiencia 2025 (Dell Inc. PowerEdge R7715)

![Fig. 1: Gráfica sacada de [@specdell2025] sobre el procesador AMD EPYC 9845, el eje X representa la carga del sistema, desde "activo en espera" (izquierda), hasta 100\% (derecha). El eje Y (izquierda) representa la métrica R (ssj_ops/w). El eje Y (derecha) muestra la métrica $P_{IT}$ (W).](resources/T52_dell2025.png)

## Compara los siguientes tres modos de operación para una carga media:

### Trabajar con tres servidores al 20\%

Cada uno de los servidores consume alrededor de 163 W. Por lo tanto en este modo de operación se consume un total de **489 W**.

En cuanto a ssj_ops/watt contamos con 21.874 por cada servidor. Con los tres sigue siendo la misma cantidad.


### Un servidor al 60\% con los otros dos apagados

Un servidor al 60% consume alrededor de 267 W y tiene un rendimiento de 40.100 ssj_ops/watt. Los servidores apagados ni gastan ni suman rendimiento.

### Un servidor al 60\% con uno apagado y el otro idle

El servidor al 60% (267 W y 40.100 ssj_ops/watt), el apagado (0 W y 0 ssj_ops/watt) y el que está en idle (70 W y 0 ssj_ops).

Esto hace un total de 337 W y $\frac{10.717.256 \text{\ ssj\_ops}}{337 W} = 31801,94$  ssj_ops/watt.

### Comparación

|Modo opreación|$P_{IT}$ (W)| ssj_ops/watt |
|:-:|:-:|:-:|
|Tres servidores al 20%|489 W|21.874 ssj_ops/watt|
|Un servidor al 60% y dos apagados|267 W|40.100 ssj_ops/watt|
|Un servidor al 60%, otro apagado y otro idle|337 W|31.801,94 ssj_ops/watt|

Nuevamente, el modo más eficiente por vatio es un servidor al 60% y los otros apagados (40 100 ssj_ops/watt). El modo mixto pierde eficiencia por el consumo del idle, quedando en 31.801,94 ssj_ops/watt, pese a ello, sigue siendo mucho mejor que el escenario de tres servidores al 20% en rendimiento por vatio, pero menos óptimo que concentrar la carga en un solo servidor activo.

# Conclusión

En 17 años, la eficiencia energética y el rendimiento han mejorado de forma abrumadora:

- Rendimiento total: el procesador de 2025 logra unas 60 veces más ssj_ops totales que el de 2008.
- Eficiencia energética: el salto en ssj_ops/watt es de 501 $\rightarrow$ 21.874, es decir, una mejora de más de 4300% en eficiencia por vatio.

En resumen, el procesador moderno no solo consume más potencia absoluta, sino que convierte la energía en trabajo útil con una eficiencia exponencialmente mayor, haciendo posible más computación con un coste energético mucho menor por unidad de rendimiento.

# Referencias

\small
::: {#refs}
:::
\normalsize

