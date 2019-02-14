#!/bin/sh

arr=( 1 2 3 4 5 )
unset arr[2]
echo ${arr[@]:2} #${arr[1]}

declare -A arr2
arr2=( [one]=one-1 )
arr2[two]=two-2
echo ${arr2[*]}
echo ${!arr2[*]}
