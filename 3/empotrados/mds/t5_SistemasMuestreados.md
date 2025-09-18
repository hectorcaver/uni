% Sistemas Muestreados.
    Sistemas empotrados 1
% Héctor Lacueva Sacristán
% 24/04/2025

\newpage

# Introducción

|Tratamiento digital de la señal|
|:-:|
|![Esquema conexión computador y mundo analógico](./mds/resources/trat_digital.png)|
|Tratamiento digital de señal simplificado|
|![Tratamiento digital de la señal](./mds/resources/trat_digital_simply.png)|
|Control|
|![Control](./mds/resources/trat_digital_ctrl.png)|

## Definiciones

- **Señales discretas**: p.e. las señales generadas por un computador.
  - Se representan por una secuencia {$X_k$}.
    - p.e. ${X_k} = {0, 1.5, 1.66, 2, 2.77, \cdots}$
- **Señales muestreadas**: una señal continua al ser muestreada da lugar a una señal discreta.

El computador tiene un algoritmo **Ecuaciones en diferencias**:

- $x(k+1) = ax(k) + bx(k-1) + cx(k-2) + du(k)$.
- Ecuaciones en diferencias = transformador de secuencias.

## Conversor analógico-digital (CAD) o (ADC)

En MSP430:

- Un solo convertidor, 16 canales multiplexados (12 al exterior).
- **Precisión de 12 bits**, aproximaciones sucesivas.
- Modos: conversión continua, conversión única.
- `Flag` indica fin de conversión. Puede generar interrupciones.
- VREFL $\le$ V $\le$ VREFH $\rightarrow$ 0x000 $\le$ CONVERSION $\le$ 0xFFF.
- Hasta 200 kHz.
- Inicio conversión:
  - **Software**: escritura del bit ADCCTL0.ADCSC = 1.
  - **Hardware**: Pin externo, RTC overflow, TB1.1.
- Características adicionales:
  - Sensor de temperatura.
  - Comparador.


# Muestreo y reconstrucción

## Teorema de muestreo (Shannon)

Una señal continua $x(t)$ cuya transformada de Fourier $X(\omega)$ sea de banda limitada $(\omega_s)$, estará completamente determinada por la secuencia ${x_k}$ obtenida por el muestreo de la misma si:
$$\omega_m \ge 2\omega_s \Rightarrow T \le \frac{\pi}{\omega_s}$$

- Si $\omega_m > 2\omega_s$: reconstrucción perfecta por filtrado.
- Si $\omega_m < 2\omega_s$: no se puede reconstruir.
- En la práctica existen señales de las que no se puede decir que sean de banda limitada $\rightarrow$ filtros.