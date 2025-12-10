#/bin/bash

perf stat -a -e power/energy-cores/ -r 5 ./pi_taylor_parallel 100000000 1
perf stat -a -e power/energy-cores/ -r 5 ./pi_taylor_parallel 100000000 2
perf stat -a -e power/energy-cores/ -r 5 ./pi_taylor_parallel 100000000 4
perf stat -a -e power/energy-cores/ -r 5 ./pi_taylor_parallel 100000000 8
perf stat -a -e power/energy-cores/ -r 5 ./pi_taylor_parallel 100000000 16
