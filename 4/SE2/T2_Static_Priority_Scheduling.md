# Tema 2: Static Priorities Scheduling

## Objetivos

El objetivo principal es **analizar y asegurar el cumplimiento de plazos (deadlines) de tareas críticas del sistema**.

Para ello se presentan:

- Guías de diseño de sistemas de tiempo
  - Tareas periódicas
  - Tareas esporádicas
  - Comunicación a través de servidores
- Técnicas de priorización estática:
  - Rate Monotonic (RM): **priority to the most frequent**
    - RMS/RMA: Rate Monotonic Scheduling / Analysis
    - Response time = period in periodic tasks.
  - Deadline Monotonic (DM): **priority to the most urgent**
    - DMS/DMA: Deadline Monotonic Scheduling / Analysis
    - Deadline for response $\le$ period in periodic tasks

## Métricas

### Factor de Utilización (Utilization factor)

Medida de la **carga de un procesador**. El objetivo es encontrar métodos que produzcan una planificación aceptable con factores de utilización lo más altos posibles.



|Utilización de una tarea $T_i$|Utilización del sistema|
|:-:|:-:|
|$U_i = \frac{C_i}{P_i}$|$U = \sum_{i = 1}^{n}{U_i} = \sum_{i = 1}^{n}{\frac{C_i}{P_i}}$|


