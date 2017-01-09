#!/bin/bash

mongoimport -h $1 --db nbhood --collection nbhs --file mongo_insert.json
