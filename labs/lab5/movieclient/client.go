// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"
	"labs/lab5/movieapi"

)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
	defaultYear = "1994"
	defaultDirector = "Quentin Tarantino"
	defaultCast = "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"
)

func main() {
	// Set up a connection to the server.
	conn, err1 := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err1 != nil {
		log.Fatalf("did not connect: %v", err1)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Get movie info from the database.
	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) == 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err2 := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err2 != nil {
		log.Fatalf("could not get movie info: %v", err2)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())

	// Adding a movie to the database.
	year := defaultYear
	director := defaultDirector
	cast := defaultCast
	if len(os.Args) > 1 {
		title = os.Args[1]
		year = os.Args[2]
		director = os.Args[3]
		cast = os.Args[4]

		r,err3 := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: title, Year: int32(year), Director: director, Cast: cast})
	}
}
