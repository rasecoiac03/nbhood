#!/bin/bash

cat nbhs_init | awk -F"|" '{print "{\"state\":\""$1"\", \"city\":\""$2"\", \"nbh\":\""$3"\", \"feature\":\""$4"\", \"featureName\":\""$5"\", \"address\":\""$6"\"}"}' > mongo_insert.json
