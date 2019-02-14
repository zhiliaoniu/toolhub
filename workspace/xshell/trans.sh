#!/bin/bash 
rm -fr *.xsh
while read line;do cp bak/45.255.134.122.xsh ${line}.xsh; sed -i "s/45\.255\.134\.122/$(echo ${line})/g" ${line}.xsh;done < a
sz -b *.xsh
