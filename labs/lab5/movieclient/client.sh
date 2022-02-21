#!/bin/bash

go run client.go "Elf" "2003" "Jon Favreau" "Will Ferrell,James Caan,Zooey Deschanel"
go run client.go "Avengers: Infinity War" "2018" "Anthony Russo" "Robert Downey Jr.,Chris Hemsworth,Chris Evans"
go run client.go "The Shawshank Redemption" "1994" "Frank Darabont" "Tim Robbins,Morgan Freeman,Bob Gunton"
go run client.go "The Green Mile" "1999" "Frank Darabont" "Tom Hanks,Michael Clarke Duncan,David Morse"
echo ""
echo ""
go run client.go "Elf"
go run client.go "Avengers: Infinity War"
go run client.go "The Shawshank Redemption"
go run client.go "The Green Mile"
