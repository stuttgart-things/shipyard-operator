#!/bin/bash

files=$(ls ~/.kube/*.yaml )
i=1

for j in $files
do
echo "$i.$j"
i=$(( i + 1 ))
done

echo "Enter number"
read input
echo ${file[$input]}