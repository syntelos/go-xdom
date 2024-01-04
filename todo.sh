#!/bin/bash

egrep -n TODO $(find . -type f -name '*.go') | egrep -v ':var TODO'

if [ -f TODO.txt ]
then
    cat -n TODO.txt | sed '1d'
fi
