module labs

go 1.17

require (
	github.com/google/go-cmp v0.5.5
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220218161850-94dd64e39d7c // indirect
)

require "github.com/gbass2/CloudNativeArchitecture/tree/main/labs/lab7/movieapi" v0.0.0
replace "github.com/gbass2/CloudNativeArchitecture/tree/main/labs/lab7/movieapi" v0.0.0 => "/Volumes/SHARE/UNCC/Cloud Native Architecture/CloudNativeArchitecture/labs/lab7/movieapi"
