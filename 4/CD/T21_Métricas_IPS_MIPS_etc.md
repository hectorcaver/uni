% Ejercicios Tema 2  
  Centros de Datos UZ, 2025-26  
% Lacueva Sacristán, Héctor
% 16/09/2025

# T2.1. Métricas (45 min, 15 min calcular)

IA: No

Dados, el siguiente programa y algunos datos calcula lo siguiente:

```{fortran}
Real*8 A(1000), B(1000), C(1000)
DO i=1, 1000
        A(i) = B(i) * C(i) + 3.0
ENDDO
```

Datos:

- 8 instrucciones de LM por iteración
- 4 Bytes / instrucción
- 2 operaciones de PF
- 1 $\mu s$ todo el bucle
- 50 W

Calcular IPS, MIPS, iBW, dBW, GFLOPS, EPI y Velocidad/Vatio:

$\LARGE IPS = \frac{instr}{s} = \frac{\frac{8 \text{ instr}}{iter} \times 10^3 \text{ iter}}{10^{-6}s} = 8 \times 10^9 \text{ }\frac{instr}{s}$

$\LARGE MIPS = \frac{IPS}{10^6} = 8000 \text{ MIPS}$

$\LARGE iBW = \frac{iBytes}{s} = IPS \times 4 \text{ } \frac{Bytes}{instr} = 32 \text{ } \frac{GB}{s}$

Teniendo en cuenta que hay **dos Load y un Store y datos de 8 Bytes (Real *8)**:

$\LARGE dBW = \frac{dBytes}{s} = \frac{3 \text{ dataInstr} \times 10^3 \text{ iter} \times \frac{8 \text{ Bytes}}{dataInstr}}{iter \times 10^{-6}s} = 24 \text{ } \frac{GB}{s}$

$\LARGE GFLOPS = \frac{\frac{2 \text{ FLOP}}{iter} \times 10^3 iter}{10^{-6}s} = 2 \text{ GFLOPS}$

$\LARGE EPI = \text{ energía por instrucción } (\frac{nJ}{instr}) = \frac{Pot(W \text{ o } \frac{J}{s}) \times 10^{-9}}{IPS} = 6.25 \frac{nJ}{instr}$

$\LARGE \text{Velocidad/Vatio (GFLOPS/Watt)} = 0.04 \text{ GFLOPS/W}$

