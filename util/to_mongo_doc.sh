#!/bin/bash

cat nbhs_init | awk -F"|" '{print "{\"nbh\":\""$1"\", \"feature\":\""$2"\", \"featureName\":\""$3"\", \"address\":\""$4"\"}"}' > mongo_insert.json
