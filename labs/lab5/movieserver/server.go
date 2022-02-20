// Package main implements a server for movieinfo service.
package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"strings"
	"labs/lab5/movieapi"
	// "labs/lab5/grpc-go"
)

const (
	port = ":50051"
)

// server is used to implement movieapi.MovieInfoServer
type server struct {
	movieapi.UnimplementedMovieInfoServer
}

// Map representing a database
var moviedb = map[string][]string{"Pulp fiction": []string{"1994", "Quentin Tarantino", "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"}}

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) GetMovieInfo(ctx context.Context, in *movieapi.MovieRequest) (*movieapi.MovieReply, error) {
	title := in.GetTitle()
	log.Printf("Received: %v", title)
	reply := &movieapi.MovieReply{}
	if val, ok := moviedb[title]; !ok { // Title not present in database
		return reply, nil
	} else {
		if year, err := strconv.Atoi(val[0]); err != nil {
			reply.Year = -1
		} else {
			reply.Year = int32(year)
		}
		reply.Director = val[1]
		cast := strings.Split(val[2], ",")
		reply.Cast = append(reply.Cast, cast...)

	}

	return reply, nil
}

func (s *server) SetMovieInfo(ctx context.Context, in *movieapi.MovieData) (*movieapi.MovieStatus) {
	title := in.GetTitle()
	year := in.GetYear()
	director := in.GetDirector()
	cast := in.GetCast()

	log.Printf("Received movie data")
	reply := &movieapi.MovieStatus{}

	if strings.TrimSpace(title) == "" {
		reply.Message = "Invalid movie title. Title is blank"
		return reply
	}
	if strings.TrimSpace(director) == "" {
		reply.Message = "Invalid movie director. Director is blank"
		return reply
	}
	if len(cast) < 0 {
		reply.Message = "Invalid movie cast. Cast is blank"
		return reply
	}

	if year < 1900 || year > 2022{
			reply.Message = "Invalid movie year. Year needs to be between 1900 and 2022"
			return reply
		}
	}

	moviedb[title] = [3]string{year,director,cast}

	reply.Message = ""
	return reply
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	movieapi.RegisterMovieInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
