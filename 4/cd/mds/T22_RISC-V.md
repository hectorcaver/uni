% Ejercicios Tema 2  
  Centros de Datos UZ, 2025-26  
% Lacueva Sacristán, Héctor
% 22/09/2025

# T2.2. RISC-V (2h 30min)

IA: No

## Característica RISC-V

La característica elegida son las extensiones del repertorio de instrucciones para añadir funcionalidades. La información ha sido extraida de la web [^1].

La arquitectura parte de un núcleo base simple que contiene únicamente operaciones enteras básicas, saltos y accesos a memoria. Estos, pueden no ser suficiente para obtener un procesador avanzado por lo que se añaden extensiones que amplían las capacidades del procesador en función de las necesidades de la aplicación.

Existen muchas extensiones pero las más comunes son:

- **M** (para multiplicación y división enteras).
- **A** (para instrucciones atómicas en multiprocesadores).
- **F/D** (para operaciones en coma flotante de simple, doble precisión).
- **C** (para poder utilizar instrucciones comprimidas de 16 bits que reducen el tamaño del código y el consumo).
- **V** (para operaciones vectoriales).

Trabajar de esta manera permite crear procesadores adaptados a cada situación, sin hardware innecesario y por tanto, un menor consumo y coste.

## Diagrama de organización relacionado con el RISC-V "atrevido". Describir la función de los bloques principales

Según la web [^3] el core RISC-V "Atrevido 423" es un procesador Out-Of-Order de 64-bit con soporte para Multiprocesadores ideal para *"AI Inference, Key-Value Stores, Recommendation Systems, Sparse Data & HPC"*. 

Dada la imagen obtenida de la web [^2] que nos muestra un diagrama de organización del core RISC-V "Atrevido 423" de la compañía Semidynamics podemos diferenciar los siguientes bloques:

- **Caché de instrucciones y datos separadas**, ambas pueden ser desde 4KB hasta 32KB. Algo no se corresponde en la imagen ya que en la web [^3] oficial, se especifica desde 8KB hasta 32KB.
- TLB para instrucciones y datos: se encargan de traducir direcciones virtuales a físicas.
- **Renombre de registros (Renamer)** que sirve para reducir dependencias.
- Sistema de predicción de saltos:
  - **TAGE**: sistema avanzado para predicción de saltos.
  - **BTB** (Branch Target Buffer): tabla que guarda en cada entrada, la dirección destino de la última ejecución de un salto para lanzar cuanto antes un salto.
  - **RAS** (Return Address Stack): gestiona direcciones de retorno de llamadas a funciones.
- 5 ventanas de lanzamiento (1 para Memoria, 2 para enteros, 1 para saltos, 1 para operaciones en PF).
- Así mismo cuenta con 5 UF (1 AGU para calculo de direcciones efectivas, 2 ALU para enteros (1 con soporte de criptografía y otra de manipulación de bits), 1 BR para cálculo de direcciones de salto y 1 FPU).
- Banco de registros de enteros (7 puertos de lectura y 3 de escritura) y banco de registros de PF (3 puertos de lectura y 1 de escritura).
- **Protección de memoria**:
  - PMP (Physical Memory Protection).
  - ECC (Error Correction Code).
  - MMU (Memory Management Unit).
- **Gazzillion Unit**: unidad para gestionar grandes flujos de datos (orientada a Big Data / IA).
- AXI/CHI (512b / 1024b): interfaz de bus para conectar con memoria externa o coprocesadores.
- Debug: módulo de depuración.
- PMU (Performance Monitoring Unit): mide el rendimiento del procesador.

[^1]: <https://dzone.com/articles/introduction-to-the-risc-v-architecture>

[^2]: <https://www.wnie.online/semidynamics-announces-fully-customisable-4-way-atrevido-423-risc-v-core-for-big-data-applications/>

[^3]: <https://semidynamics.com/en/products/atrevido>
