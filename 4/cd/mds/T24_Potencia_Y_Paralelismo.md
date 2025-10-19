% Ejercicios Tema 2  
  Centros de Datos UZ, 2025-26  
% Lacueva Sacristán, Héctor
% 24/09/2025


# T2.4. Potencia y paralelismo (45 min)

IA: No

Un programa se ejecuta en un core a 4 GHz y 0,9V consumiendo 50 W (40 W dinámica, 10 W estática). Puede ejecutarse en dos cores a 2 GHz por core y 0,7V.

1. **Calcula la nueva potencia, asumiendo que la potencia estática es proporcional a la tensión**:

    - Calculamos la nueva potencia estática:

    $$P_{est_{new}} = T_{new} \times \frac{P_{est_{old}}}{T_{old}} = 0,7V \times \frac{10W}{0,9V} = 7, \overline{7}W$$

    - Calculamos la nueva potencia dinámica:

        A partir de esta fórmula $P_{din} = \alpha \times C \times V_{dd}^2 \times F$ y teniendo en cuenta que C depende del número de cores funcionando (más o menos transistores funcionando):

    $$P_{din_{new}} = 2 * C_{core} \times 0,7^2 \times 2 GHz$$
    $$P_{din_{old}} = C_{core} \times 0,9^2 \times 4 GHz = 50 W$$
    $$\frac{P_{din_{new}}}{P_{din_{old}}} = \frac{2 \times 0,7^2 \times 2GHz}{0.9^2 \times 4GHz} \longrightarrow P_{din_{new}} = \frac{0,7^2}{0,9^2} \times P_{din_{old}} = 24,1975 W$$


      Entonces la potencia total final será $P_{T_{new}} = P_{din_{new}} + P_{est_{new}} = 31,975W \approx 32W$.

2. **Calcula la variación de potencia en porcentaje**:

    $$ \% \Delta P = \frac{P_{new} - P_{old}}{P_{old}} \times 100 = -36,06\%$$

