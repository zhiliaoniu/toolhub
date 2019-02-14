#!/bin/bash

#test
ss="conn=100  tcp=0   udp=2   http=3"
echo $ss | awk -F "[= ]+" '
{
	print $1, $2; 
	count+=$2; 
	count+=$4; 
	count+=$6; 
	count += $8;
	if(count > 100){
		print "more than" count
	}
} 
END{
	print count
}'
exit 0
