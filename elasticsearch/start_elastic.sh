#!/bin/sh

sudo sysctl -w vm.max_map_count=262144
docker run -p 9200:9200 -d elasticsearch
sleep 10
curl "http://127.0.0.1:9200/_cat/health?v"
