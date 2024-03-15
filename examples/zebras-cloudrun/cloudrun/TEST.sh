#!/bin/bash

export APP=https://zebras-example-y2k5eue77q-uw.a.run.app

for i in {1..100}
do

echo
echo $i: creating zebras
curl $APP/v1/zebras -d '{"id":"1", "name":"alice"}'
curl $APP/v1/zebras -d '{"id":"2", "name":"bob"}' 
curl $APP/v1/zebras -d '{"id":"3", "name":"charley"}' 

echo
echo $i: fetching zebras
curl $APP/v1/zebras

echo
echo $i: fetching zebra 2
curl $APP/v1/zebras/2

done
