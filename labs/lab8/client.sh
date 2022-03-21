#!/bin/bash

curl -X GET "http://localhost:8000/list"
sleep 3
curl -X POST "http://localhost:8000/create?item=popcorn&price=10"
sleep 3
curl -X GET "http://localhost:8000/price?item=popcorn"
sleep 3
curl -X POST "http://localhost:8000/update?item=popcorn&price=5"
sleep 3
curl -X POST "http://localhost:8000/delete?item=popcorn"
