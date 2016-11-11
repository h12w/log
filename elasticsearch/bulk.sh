#!/bin/sh

PORT=32773

curl -i -X POST --data-binary @req.json "http://127.0.0.1:$PORT/idx-test/m/_bulk"
sleep 3
curl "http://127.0.0.1:$PORT/idx-test/_search?pretty=true"

