

# Dataframes pandas

## Crear un dataframe

```python
data = list[valores]
index = list[nombre_filas]
columns = list[nombre_columnas]

# Crea una dataframe de pandas con los datos
df = pd.DataFrame(data, index, columns)
```

Por ejemplo:

```python
data = [[1,2],[3,4]]
index = ['Chocolate', 'Chucherías']
columns = ['€ Tienda1', '€ Tienda2']

df = pd.DataFrame(data, index, columns)
df
```

## Concatenar dataframes

```python
pieces = [df1, df2, df3]

# Concatena los dataframes en uno solo
pd.concat(pieces)

# Adding a column to a DataFrame is relatively fast. However, adding a row requires a copy, and may be expensive. We recommend passing a pre-built list of records to the DataFrame constructor instead of building a DataFrame by iteratively appending records to it.
```

# Métricas de evaluación del ajuste de los datos

## $RMSE$

## $R^2$

$R^2$ mide la varizanza en las predicciones relativas a una constante simple de predicción.

Si un modelo no predice mejor que usando la media de y, $R^2 = 0$. Si el modelo se ajusta a la perfección a los datos $R^2 = 1$. Por lo general, valores más grandes implican un mejor ajuste.


# Como elegir el tamaño del paso (factor de aprendizaje o $\alpha$) para SGD

Para conseguir la convergencia hacia un mínimo, tenemos que ser cuidadosos a la hora de elegir el ratio de aprendizaje ($\alpha$).

Una buena heurística para elegir un buen $\alpha$ consiste en empezar con un valor pequeño e ir incrementando. Una vez tengas suficientes valores, eliges el $\alpha$ con menor error.



