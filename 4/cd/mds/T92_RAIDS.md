% Ejerncicios Tema 9  
  Centros de Datos UZ, 2025-26  
% Lacueva Sacristán, Héctor
% 09/11/2025

\newpage

# T9.1 Raids (Tiempo: 3h)

IA: Sí, ChatGPT.

## Configuraciones de RAIDs de 4 discos

### RAID 01

|Configuración de 4 discos RAID 01|
|:-:|
|![Configuración de 4 discos RAID 01](resources/raid_01.png){ width=50% }|

RAID 01 (0+1, 4 discos): dos stripes (RAID 0) y mirror entre esas stripes. Estructuralmente parecido a RAID 10 en rendimiento ideal, pero más frágil ante fallos múltiples.

### RAID 10

|Configuración de 4 discos RAID 10|
|:-:|
|![Configuración de 4 discos RAID 10](resources/raid_10.png){ width=50% }|

RAID 10 (1+0, 4 discos): dos mirrors y stripe entre ellos. Los datos están duplicados por pares.

### RAID 5

|Configuración de 4 discos RAID 5|
|:-:|
|![Configuración de 4 discos RAID 5](resources/raid_5.png){ width=50% }|

RAID 5 (4 discos): datos + paridad distribuida (XOR). Cada fila de bloques tiene su bloque de paridad en un disco distinto, permite reconstruir 1 disco perdido.

RAID 5 (8 discos): igual esquema escalado, paridad distribuida entre 8 discos (se usa el equivalente a 1 disco para paridad).

## Detalles de las Configuraciones

### Suposiciones

- **Capacidad disco** = 1 TB.
- **Coste** = 100€/TB.
- Throughput secuencial por disco = 200 MB/s.
- IOPS/disco = 150 IOPS/disco.
- **Tasa de fallo**: 5% en 3 años.
- Tamaño a reconstruir = 1 TB = $10^{12} B$ para simplificación.
- **Velocidad efectiva de recuperación**: 100 MB/s.

### Costes y capacidad

|Topología|Discos|Capacidad usable|Overhead|Coste total	€| €/TB usable|
|:-:|:-:|:-:|:-:|:-:|:-:|
| RAID 10| 4 |	2 TB | 50%   | 400€  | 200€/TB     |
| RAID 01| 4	| 2 TB | 50%	  | 400€	| 200€/TB	    |
| RAID 5-4d	| 4	| 3 TB | 25%	  | 400€	| 133,33€/TB  |
| RAID 5-8d | 8 | 7 TB | 12,5%	| 800€  |	114,29€/TB	|

### Rendimiento

| Topología |  Lectura secuencial | Escritura secuencial |   Random Read IOPS  |  Random Write IOPS |
| :--------: | :-----------------: | :------------------: | :-----------------: | :----------------: |
|   RAID 10  |  4 × 200 = 800 MB/s |  2 × 200 = 400 MB/s  |  4 × 150 = 600 IOPS | 2 × 150 = 300 IOPS |
|   RAID 01  |  4 × 200 = 800 MB/s |  2 × 200 = 400 MB/s  |  4 × 150 = 600 IOPS | 2 × 150 = 300 IOPS |
|  RAID 5-4d |  3 × 200 = 600 MB/s |  3 × 200 = 600 MB/s  |  3 × 150 = 450 IOPS |      ~112 IOPS     |
|  RAID 5-8d | 7 × 200 = 1400 MB/s |  7 × 200 = 1400 MB/s | 7 × 150 = 1050 IOPS |      ~263 IOPS     |


### Fallos 

| Topología | Tolerancia a fallos |            Probabilidad de pérdida (en 3 años)            |
| :---------: | :-----------:| :-------------------------------------------------------: |
|  RAID 10  |  1 (2 si distintos) | ~0,50% (Fallos >= 3 discos + Fallos = 2 y distinto espejo) |
|  RAID 01  |  1 (2 si mismo stripe) | ~0,50% (Fallos >= 3 discos + Fallos = 2 y distinto stripe) |
| RAID 5-4d |          1          | ~1,40% (Fallos >= 2 discos)                |
| RAID 5-8d |          1          | ~5,73% (Fallos >= 2 discos)                |

### Rendimiento en fallos

| Config | Estado en fallo                                                     | Rendimiento en fallo                    | Comentarios técnicos                                                                         |
| :------------ | :------------------------------------------------------------------ | :-------------------------------------- | :------------------------------------------------------------------------------------------- |
| RAID 01   | Si falla un disco OK. Si falla un espejo completo, fallo total | Muy degradado o inoperativo             | Stripe depende de ambos espejos. Si un espejo falla, el stripe se rompe.                     |
| RAID 10   | Si falla un disco OK. Si falla un par, fallo total             | Degradado pero funcional                | Cada par puede seguir operando si al menos un disco está sano. Lectura sigue paralela.       |
| RAID 5-4d    | Si falla un disco OK. Si fallan >=2, Fallo total                | Lectura: OK Escritura: muy degradada | Lectura usa los discos restantes. Escritura requiere cálculo de paridad en tiempo real.      |
| RAID 5-8d    | Si falla un disco OK.Si fallan >=2, Fallo total                | Lectura: OK Escritura: muy degradada | Mayor volumen de datos implica más cálculos de paridad y más presión sobre discos restantes. |

### Tiempo de recuperación

| Config | Qué se reconstruye              | Nº discos reconstruidos | Tiempo estimado                                                               |
| ------------- | ------------------------------- | ----------------------- | ----------------------------------------------------------------------------- |
| RAID 01   | Se reconstruye el stripe completo                 | 2 discos  | 5.56 h  |
| RAID 10    | Disco individual de mirror      | 1 disco   | 2.78 h  |
| RAID 5-4d     | Disco individual usando paridad | 1 disco  | 2.78 h (sería más por la necesidad de leer de tres discos y calcular la paridad)  |
| RAID 5-8d     | Disco individual usando paridad | 1 disco   | Teórico 2.78 h (aúm mayor al RAID 5-4d porque en este caso hay que leer 7 discos en vez de 3)|

## Conclusiones

RAID 10 ofrece el mejor equilibrio entre rendimiento, tolerancia a fallos y tiempo de recuperación. Su estructura permite lecturas muy rápidas y una recuperación sencilla (solo un disco), aunque con un coste en capacidad del 50%.

RAID 01 tiene un rendimiento similar en condiciones normales, pero es más vulnerable ante fallos múltiples, ya que la pérdida de un espejo completo implica la caída total del sistema.

RAID 5 es más eficiente en espacio (solo se pierde el equivalente a un disco), pero sufre una degradación notable del rendimiento en escritura y durante la reconstrucción, especialmente en configuraciones con muchos discos (8d).

A medida que aumenta el número de discos, la probabilidad de pérdida de datos crece rápidamente, como se observa en RAID 5-8d (~5,7% en 3 años frente a ~1,4% en RAID 5-4d).

El tiempo de reconstrucción también se incrementa con el tamaño del array, lo que amplía la ventana de riesgo de pérdida adicional.

# Referencias

- Teoría, ChatGPT (ayuda con los cálculos).

- Vídeo de Youtube <https://www.youtube.com/watch?v=YYMQDZFILzE>
