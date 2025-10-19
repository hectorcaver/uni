% Ejercicios Tema 4  
  Centros de Datos UZ, 2025-26  
% Lacueva Sacristán, Héctor
% 08/10/2025


# T4.1 Certificación Energy Star (45 min)

IA: No

Calcular la potencia máxima en estado idle que se le permite a un servidor de la siguientes características para que obtenga la certificación Energy Star:

- No gestionado (unmanaged)
- Uniprocesador
- 8 GB RAM
- 2 HDs
- 2 I/O devices (uno con dos puertos de 1 Gbit y el otro con 1 Gbit port).

Los valores y el procedimiento se han extraido de la documentación [@energystar2013]: 

**1. Primero se debe obtener el "Base Idle Power Allowance":**

![Tabla que relaciona las diferentes categorías con especificaciones y potencia base. Obtenida de la documentación [@energystar2013]](resources/tabla_base_idle_power.png)
    Dado que contamos con un único procesador y es no gestionado pertenece al grupo A y por tanto $\textbf{BIPA = 47 W}$.

\newpage

**2. Ahora hay que obtener las "Additional Idle Power Allowances":**

![Tabla que contiene los "additional idle allowances" para los componentes extras. Obtenida de la documentación [@energystar2013]](resources/tabla_additional_idle_power.png)

  - 2 HDD: $\displaystyle 2HDD \times \frac{8W}{1HDD} = 16 W$
  - 8 GB RAM: $\displaystyle 4 GB_{RAM} \times \frac{0,75W}{1 GB_{RAM}} = 3 W$
  - Las tarjetas I/O no suman ya que no cumplen los requisitos (ninguna tiene más de 2 puertos).


**3. Calculamos la "Idle Allowance" sumando todos los valores que hemos obtenido anteriormente:**


$$
\textbf{Idle Allowance = 47 W + 16 W + 3 W = 66 W}
$$

  Por lo tanto la potencia máxima en estado idle para un servidor de estas características sería de 66 W.

# Referencias
\small
::: {#refs}
:::
\normalsize

