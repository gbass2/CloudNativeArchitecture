// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"
	"strconv"
	"strings"
	"labs/lab5/movieapi"
	"google.golang.org/grpc"
)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
	defaultYear = 1994
	defaultDirector = "Quentin Tarantino"
	defaultCast = "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Get movie info from the database.
	// Contact the server and print out its response.
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	title := defaultTitle

	if len(os.Args) < 3 {
		title = os.Args[1]

		r1, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
		if err != nil {
			log.Fatalf("could not get movie info: %v", err)
		}
		log.Printf("Movie Info for %s %d %s %v", title, r1.GetYear(), r1.GetDirector(), r1.GetCast())
	}

	// Set movie info in database.
	if len(os.Args) > 3 {
		// Getting the movie info from the command line args.
		title = os.Args[1]
		year, _ := strconv.Atoi(os.Args[2])
		director := os.Args[3]
		cast := os.Args[4]

		// Splitting cast into a slice
		castSlice := strings.Split(cast,",")

		// Adding movie to database
		r2,_ := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: title, Year: int32(year), Director: director, Cast: castSlice})

		// If there was an error message then exit.
		if strings.TrimSpace(r2.GetMessage()) != "" {
			log.Fatalf(r2.GetMessage())
		} else {
			// Query database for movie info of recently movie added
			log.Printf("Movie added to database with no errors")
			r3, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})

			if err != nil {
				log.Fatalf("could not get movie info: %v", err)
			}

			log.Printf("Movie Info for %s %d %s %v", title, r3.GetYear(), r3.GetDirector(), r3.GetCast())
		}
	}
}
