% Ejercicios Tema 6 \ 
 Centros de Datos UZ, 2025-26  
% Lacueva Sacristán, Héctor
% 28/10/2025


# T6.1 Consumos anuales de AWS-Aragón (Tiempo: 2h) 

IA: Sí, para obtener información/enlaces del BOE relacionados con el tema de AWS y para que realizara los cálculos una vez obtenidos los valores de los índices.

## Leer el documento “2024 Modificación nº2 PIGA AWS Anexo 1 Memoria Justificativa.pdf” para entender mejor la climatización de AWS en Aragón.

El documento simplemente se podría resumir en lo siguiente:

- Modificaciones en la gestión del agua en los diferentes emplazamientos que permiten reducir en un 10% el consumo de agua de los centros de datos.
- Modificación del sistema de filtrado del agua, pasando de la tecnología de ósmosis inversa a la de nanofiltración (+ eficiencia).
- Petición de conexión adicional del centro de datos de El Burgo de Ebro a la red pública de agua potable y al colector
de aguas pluviales.

## Calcular el consumo anual de agua y de electricidad de sus tres centros de datos (consumos agregados de los tres).

Según los datos obtenidos del [@boeAWSinfoAG2025]:

|Centro de datos|Consumo hídrico|Consumo eléctrico|
|:-:|:-:|:-:|
|El Burgo de Ebro|111.880 m³/año|1.766 GWh/año|
|Huesca|143.400 m³/año|2.270,6 GWh/año|
|Villanueva de Gállego|56.980 m³/año|756,9 GWh/año|
|---|---|---|
|**Total agregado**|**312.260 m³/año**|**4.793,5 GWh/año**|

## Trasladar esos consumos a ciudades de un determinado tamaño, con los índices de Zaragoza capital.

Según los datos en [@indiceElectricidadZaragoza] y [@indiceAguaZaragoza] los índices de Zaragoza se podrían considerar los siguientes:

|Índice|Valor|
|:-:|:-:|
|Consumo eléctrico|**3000 GWh/año**|
|Consumo hídrico|**220 litros por persona al día (consumo general)**|
|Número de habitantes|**686.986 habitantes** según [@enterat]|

Si esto lo llevamos a ciudades con los mismos índices con diferente número de habitantes obtenemos lo siguiente:

| Ciudad (habitantes)    | Consumo total de agua (m³/año) | % que representa 312.260 m³ | Consumo eléctrico ciudad (GWh/año) | % que representa 4.793,5 GWh |
| :-: | :-: | :-: | :-: | :-: |
| 10.000 | 803.000 | **38,9 \%** | 43,66 | **10.985 \%** |
| 50.000 | 4.015.000 | **7,8 \%** | 218,3 | **2.197 \%** |
| 100.000 | 8.030.000 | **3,9 \%** | 436,6 | **1.098 \%** |
| 250.000 | 20.075.000 | **1,56 \%** | 1.091,5 | **439 \%** |
| **Zaragoza (686.986)** | **55.172.000** | **0,57 \%** | **3.000** | **159,8 \%** |

## Compara tus números con la noticia de cabecera: ¿buena calidad?

Según la noticia [@noticiaCabecera] el consumo eléctrico de AWS en aragón se corresponde al de una ciudad de 300.000 habitantes cuando según los datos del BOA y los índices de Zaragoza tendría que ser de 1,1 millones de habitantes aproximadamente.

Al no mostrarse ningún cálculo ni ningún indicio de que estudio han realizado para llegar a esos números, diría que no es de buena calidad. Cabe la posibilidad de que bajo alguna condición que no se menciona se llegue a dicho resultado.

# Referencias

\small
::: {#refs}
:::
\normalsize
