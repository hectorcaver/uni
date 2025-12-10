#/bin/bash

perf stat -a -e power/energy-cores/ -r 5 ./naive_matrix_multiplication 1
perf stat -a -e power/energy-cores/ -r 5 ./naive_matrix_multiplication 2
perf stat -a -e power/energy-cores/ -r 5 ./naive_matrix_multiplication 4
perf stat -a -e power/energy-cores/ -r 5 ./naive_matrix_multiplication 8
perf stat -a -e power/energy-cores/ -r 5 ./naive_matrix_multiplication 16
