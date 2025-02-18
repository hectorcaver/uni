---
title: "Tema 2: Regularización y Selección de Modelos"
subtitle: "Aprendizaje automático"
author: "Héctor Lacueva Sacristán"
---

\newpage

# Generalización

Es la capacidad de predecir la salida para nuevos datos.

![Ejemplo de ajuste polinómico. En esta imagen se observa como ajustar un modelo simple a unos datos vs. un modelo complejo. Se puede ver que el modelo simple no ajusta perfectamente los datos de entrenamiento, pero hace un trabajo decente para nuevos datos. En el segundo caso, se ajusta perfectamente a los datos de entrenamiento pero hace un trabajo pésimo para los nuevos datos. Por tanto, el segundo modelo generaliza peor que el primero. IMPORTANTE[^0]](mds/resources_t2/ejemplo_ajuste_polinomico.png)

[^0]: **Principio de la navaja de Occam**: En igualdad de condiciones, **elegir la hipótesis más simple**.


# Sobreajuste y Subajuste

![Representación gráfica del sobreajuste, subajuste y punto ideal](mds/resources_t2/repr_ajuste_funcion.png)

![Ejemplo en **regresión**. Para una función muy simple se produce subajuste (por alto sesgo), para una función muy compleja se produce sobreajuste (por alta varianza, si la complejidad es la correcta el resultado será el deseado.)](mds/resources_t2/ejemplo_ajuste_funciones_regr.png)

![Ejemplo en **clasificación**. Para una función muy simple se produce subajuste (por alto sesgo), para una función muy compleja se produce sobreajuste (por alta varianza, si la complejidad es la correcta el resultado será el deseado.)](mds/resources_t2/ejemplo_ajuste_funciones_clas.png)


## Definiciones

- Incluir definición de **Sesgo**.
- Incluir definición de **Varianza**.

## Sobreajuste

Si hay demasiados atributos, la hipótesis puede **ajustarse muy bien a los datos de entrenamiento**, pero **puede no generalizar bien a nuevos ejemplos**.

## Como evita el Sobreajuste

Hay varias opciones:

### Reducir el número de atributos

- Seleccionar manualmente los atributos a mantener.
- Selección de modelos (validación cruzada, puede ser K-Folding-validation).

### Regularización:

- Mantener los atributos, pero reducir la magnitud de los pesos.
- Funciona bien si hay muchos atributos, y cada uno contribuye un poco a la predicción.

# Evaluación de hipótesis

## Como evaluar una hipótesis

Se evalua con datos de test distintos de los de entrenamiento.

## Como evaluar varias hipótesis

Imaginemos que queremmos elegir entre estos modelos:

1. $h_{\theta}(x) = \theta_0 + \theta_1x$
2. $h_{\theta}(x) = \theta_0 + \theta_1x + \theta_2x^2$
3. $h_{\theta}(x) = \theta_0 + \theta_1x + \cdots + \theta_3x^3$
4. $h_{\theta}(x) = \theta_0 + \theta_1x + \hdots + \theta_10x^10$

Supongamos que el menor error se da para: $J_{test}(\theta^(5))$. Por tanto, elegimos el polinomio de orden 5: $\theta_0 + \cdots + \theta_5x^5$.

El **problema** es que $J_{test}(\theta^(5))$ es una estimación optimista del error de generalización, porque el parámetro extra d[^2] se ha ajustado con los datos de test.

[^2]: d: grado del polinomio.

# Selección de modelos: Validación Cruzada

## División de datos

- **Entrenamiento**: sirven para entrenar cada modelo.
- **Validación**: sirven para comparar modelos.
- **Test**: **bajo llave hasta el final**.

## Proceso de aprendizaje

- **Aprender los parámetros** con los datos de **entrenamiento**.
- **Ajustar los hyper-parámetros** con los datos de **validación**.
- **SOLO UNA VEZ, AL FINAL**, calcular la precisión con los datos de **test**.

Es importante no espiar nunca los datos de test.

**Datos de Validación y de Test**:

 - Datos de **Validación**
   - Misma distribución que los datos de entrenamiento (p.e. escoger 20% de los datos de entrenamiento).
   - los datos de validación se gastan (cambiar cada cierto tiempo).
   - Si hay pocos datos, y el entrenamiento no es muy costoso, usar [K-fold](#k-fold-cross-validation).
 - Datos de **Test**
   - Representativos de lo que esperamos encontrar en el futuro.

## K-Fold Cross-Validation

Partir los datos en **k pliegues**. Cada dato entra en el conjunto de validación una vez. Valores típicos de $K \rightarrow 5, 10$.

## Leave-one-out Cross Validation

Cuando **hay pocas muestras de entrenamiento N**, tomar **K=N**. En cada iteración:

- N-1 muestras para el entrenamiento.
- 1 muestra para la validación.

# Errores

Habitualmente usaremos $RMSE$ ya que trabaja en las mismas unidades que el valor predicho y es independiente del número de muestras.

Los errores de entrenamiento y validación sirven para comprobar sobreajuste o subajuste. El error de validación sirve para seleccionar el mejor modelo y el de test sirve para la evaluación final del modelo elegido.

## Como detecto el sobre-ajuste o sub-ajuste

### Sub-ajuste

$$
    E_{train}(\theta) es alto
    E_{validation}(\theta) \approxeq E_{train}(\theta)
$$

### Sobre-ajuste

$$
    E_{train}(\theta) es bajo
    E_{validation}(\theta) > E_{train}(\theta)
$$

# Cómo encuentro el mejor modelo

## Búsqueda exhaustiva (grid search)

- Probar todas las combinaciones posibles de los hiper-parámetros.
- Factible si son pocos y el entrenamiento es rápido.

$$
    \text{6 hyper-parámetros, 10 valores} \rightarrow \text{1.000.000 pruebas}
$$

## Búsqueda heurística

- Probar valores para el hyper-parámetro más importante y fijarlo
- Repetir con los hyper-parámetros siguientes

$$
    \text{6 hyper-parámetros, 10 valores} \rightarrow \text{60 pruebas}
$$

## Otras variaciones

- Refinamiento sucesivo: grid basta + grid fina.
- Heurística + grid para afinar.

> **SIEMPRE SE DEBE ANOTAR EN UNA TABLA LOS RESULTADOS**

