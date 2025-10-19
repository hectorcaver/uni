% Práctica 2: 
%
%

Tarea 2: sacar tiempos de ejecución y demás.

```
perf stat -a -e power/energy-cores/ -r 5 ./naive_matrix_multiplication 1
```


```
perf stat -a -e power/energy-cores/ -r 5 ./eigen_matriz_multiplication 1
```



