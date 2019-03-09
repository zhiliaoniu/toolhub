while :

do
	ab -n 1000 -c 100 http://127.0.0.1:8001/rec
	sleep 1
	ab -n 1000 -c 100 http://127.0.0.1:8001/hot
	sleep 1
	ab -n 1000 -c 100 http://127.0.0.1:8001/get
	sleep 1
done
