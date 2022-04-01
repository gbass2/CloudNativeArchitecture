#!/bin/bash

go run create.go cloudnative4090
echo
sleep 3
go run upload_bucket.go cloudnative4090 test.txt
sleep 3
go run list_bucket.go cloudnative4090
sleep 3
go run delete_objects_bucket.go cloudnative4090
sleep 3
echo
echo
go run list_bucket.go cloudnative4090
sleep 3
go run delete_bucket.go cloudnative4090
sleep 3
