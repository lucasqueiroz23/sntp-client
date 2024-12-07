#!/bin/bash

make

for server in time.nist.gov time-d-g.nist.gov 129.6.15.28 time-b-g.nist.gov 129.6.15.30 time-d-g.nist.gov time-c-wwv.nist.gov 132.163.97.4;
do
    ./client $server
    sleep 4
done
