% Aplicaciones concurrentes.
    Ejecutivos Cíclicos
    Sistemas empotrados 1
% Héctor Lacueva Sacristán
% 24/04/2025

\newpage

# Introducción

## Definiciones

### Tareas concurrentes

Se dice que dos o más tareas son concurrentes si
pueden ejecutarse simultáneamente, de forma que alguna de ellas
comience a ejecutarse antes que termine alguna otra.

### Aplicación concurrente

Es un programa que especifica dos o más tareas
concurrentes o, de forma equivalente, un programa cuya ejecución se
realiza según varios flujos de control que avanzan en paralelo.

### Paralelísmo virtual

**Monoprocesador**: Un único procesador va alternando la ejecución de
las diversas tareas $\rightarrow$ **entrelazado**.

### Paralelísmo real

- **Multiprocesador**: Cada proceso se ejecuta en un procesador diferente.
Zona de memoria común (datos comunes).
- **Sistema distribuido**: Multiprocesador sin memoria compartida.

### Importante

En algunos casos excepcionales, puede compartirse una variable sin exclusión mutua, si se cumple:

- Un proceso la escribe y uno la lee
- La escritura y lectura se realizan en memoria: el compilador no debe utilizar copias locales en registros para optimizar el uso de la variable.
- La escritura y lectura de la variable son atómicas: "debe ser leída con una única instrucción del procesador".

# Marco temporal de una tarea

Un sistema empotrado puede ser un sistema concurrente, compuesto por
tareas. Desde el punto de vista del tiempo, varios tipos de tareas:

- **Tareas periódicas**: son aquellas que se activan regularmente en
instantes de tiempo separados por un período de tiempo determinado.
- **Tareas esporádicas**: son aquellas que se activan de forma irregular,
cada vez que se producen ciertos eventos externos.
- **Tareas de fondo**

## Definición

El **marco temporal de una tarea**: es el conjunto de atributos temporales asociados a dicha tarea:

- **T**: Periodo de ejecución (P) o separación mínima entre eventos (S).
- **C**: Tiempo de cómputo máximo.
- **R**: Tiempo de respuesta.
- **D**: Plazo de respuesta (deadline).

|Estructura temporal de una tarea|
|:-:|
|![Estructura temporal de una tarea](./mds/resources/marco_temporal.png)|

Para que la ejecución de una tarea sea correcta, debe comenzar después de su activación y terminar antes de su plazo de respuesta.

# Planificación

## Objetivo

Planificar el uso de los recursos del sistema (en particular, el procesador), para poder garantizar los requisitos temporales de las tareas.

Un método de planificación consta de:

- **Un algoritmo de planificación**, que determina el orden de acceso de las tareas a los recursos del sistema
- **Un método de análisis** que permite calcular el comportamiento temporal del sistema
  - Para comprobar que los requisitos están garantizados en todos los casos
  - Se estudia siempre el **peor caso**
  - Es necesario **conocer la duración** de las tareas en el **peor caso**


## Métodos

- Planificación estática "off-line": **Planificación cíclica**
- Planificaciones basadas en prioridades:
  - Prioridades estáticas:
    - Prioridad al más frecuente (Rate monotonic)
    - Prioridad al más urgente (Deadline monotonic)
  - Prioridades dinámicas:
    - Proximidad del plazo de respuesta (Earliest deadline first)
    - Prioridad al de menor holgura (Least Laxity First)

## Ejecutivos cíclicos

### Definición

Un ejecutivo cíclico es una **estructura de control o programa cíclico que entrelaza de forma explícita la ejecución de diversas tareas periódicas en un único procesador**. El entrelazado es fijo y está definido en el denominado plan principal que es construido antes de poner en marcha el sistema.

### Plan principal

Especificación del entrelazado de varias tareas periódicas durante un período de tiempo (ciclo principal) de tal forma que su ejecución cíclica garantiza el cumplimiento de los plazos de las tareas.

- La duración del ciclo principal es igual al mínimo común múltiplo de los periodos de las tareas.
  - **M = `mcm(Ti)`**
  - Se supone tiempo entero
  - El comportamiento temporal del sistema se repite cada ciclo principal.

### Planes secundarios

Cada **plan principal es dividido en uno o más planes secundarios o marcos ("frames")** que se ejecutarán de forma secuencial.

Cada **comienzo/fin de un marco en el ejecutivo cíclico se sincroniza con el reloj**. Son puntos donde se fuerza la corrección del tiempo real.

Por simplicidad, en la práctica, **la duración de todos los marcos es la misma**. A esta duración se le denomina **ciclo secundario**.
- Si las acciones definidas en un marco acaban antes de que concluya el ciclo secundario el ejecutivo cíclico espera.
- **Si las acciones definidas en un marco no han acabado al terminar el ciclo secundario**, se produce un **error**: **desbordamiento de marco**.
- Si la **duración de una acción es superior al ciclo secundario** debe ser descompuesta en **subacciones**.

### Propiedades

- **No hay concurrencia en la ejecución**.
  - Cada ciclo secundario es una secuencia de llamadas a procedimientos.
  - No se necesita un núcleo de ejecución multitarea.
- **Los procedimientos pueden compartir datos**.
  - No se necesitan mecanismos de exclusión mutua como los semáforos o monitores.

### Determinación de los ciclos

Sea un conjunto de tareas periódicas {$P_i / i = 1..n$}, con requisitos temporales representados por ternas $(C_i, T_i, D_i)$.

- **Ciclo principal**:
  - $M = mcm(T_i)$
- **Ciclo secundario**:
  - $m \le min(D_i)$
  - $m \ge max(C_i)$
  - $\exists k: M = km$
  - $\forall i: m + (m - mcd(m,T_i)) \le D_i$
    - Garantiza que entre el instante de activación de cada tarea y su plazo límite exista un marco o ciclo secundario completo.
    - $m - mcd(m,T_i)$ es el retraso máximo entre la activación de una tarea y el comienzo del siguiente marco
    - esta condición incluye a la primera.

### Planificación

#### Objetivo:

Asignación de procesos (o subprocesos) a los marcos de forma que se cumplan los requisitos temporales.

#### Planteamiento:

Búsqueda en el espacio de estados
- Estado: asignación parcial
- Algoritmo: búsqueda en profundidad con retroceso
  - Se pretende encontrar una única solución
- Guiado de la búsqueda:
  - **Heurísticas sobre la siguiente tarea a asignar**
    - Primero el más urgente o el más frecuente
    - Primero el de tiempo de proceso más grande
  - Si varios marcos cumplen las condiciones:
    - El primer marco que cumpla las condiciones
    - El marco con menor tiempo de cómputo libre

### Partición de tareas

Hay casos en que un conjunto de tareas que no es planificable:
- Si el tiempo de cómputo de uno es mayor que el plazo de algún otro:
  - $C_i > D_k$ => No existe valor de m que cumpla a la vez:
    - $m \le min(D_i)$
    - $m \ge max(C_i)$
- Si para una ejecución de una tarea no queda ningún marco con suficiente tiempo libre

**Solución**: descomponer la tarea demasiado larga $Pi=(C_i,T_i,D_i)$ en varias subtareas $P_{ij}=(C_{ij},T_{ij},D_{ij})$:

- $T_{ij}=T_i ; D_{ij}=D_i$
- $C_{i1} + C_{i2} + C_{i3} + \cdots = C_i$
- no partir secciones críticas
- relación de precedencia en cada una de las ejecuciones a mantener en la planificación: $P_{i1} \rightarrow P_{i2} \rightarrow P_{i3} \rightarrow \cdots$

### Partición con secciones críticas

Si es preciso partir alguna tarea en subtareas, **no deben partirse las secciones críticas** con objeto de preservar la **exclusión mutua**.

### Tareas esporádicas

Tarea esporádica E: atención a eventos externos aperiódicos
- $S_E$ : separación mínima entre eventos
- $D_E$ : plazo límite ( normalmente $D_E \le S_E$ )
- $C_E$ : tiempo de cómputo máximo
En un ejecutivo cíclico las tareas esporádicas pueden programarse de dos formas: 

#### Muestreo periódico del evento

Se programa la tarea esporádica como una tarea periódica que consulta si ha llegado un evento, y en tal caso, lo procesa.

- Transformamos la tarea esporádica en periódica con $D' = T' \le D_E / 2$, y planificamos de la forma convencional.
  - Si la tarea periódica cumple sus plazos, cualquier evento se atiende en su plazo $D_E$.

#### Por interrupción

La llegada del evento producen una interrupción, y el evento se trata inmediatamente.
Basta con reservar tiempo en cada marco para atender el máximo número de eventos que pueden llegar en un marco: $C_{res} = \lceil \frac{m}{S_E} \rceil C_E$