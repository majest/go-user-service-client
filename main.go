package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	srv "github.com/majest/go-microservice/server"
	cl "github.com/majest/user-service-client/client"
	"github.com/majest/user-service/server"
	"github.com/sony/gobreaker"
	"golang.org/x/net/context"
)

var consulIP string
var consulPort int

func init() {
	flag.StringVar(&consulIP, "consulip", "192.168.99.101", "Consul node ip")
	flag.IntVar(&consulPort, "consulport", 8500, "Consul node port")
	flag.Parse()
}

func newConfig() *srv.ClientConfig {
	return &srv.ClientConfig{
		ServiceName:     "UserService",
		GRPCSettings:    []grpc.DialOption{grpc.WithInsecure(), grpc.WithTimeout(100 * time.Millisecond)},
		MaxQPS:          100,
		MaxAttempts:     3,
		MaxTime:         100 * time.Millisecond,
		BreakerSettings: gobreaker.Settings{},
		ConsulIP:        consulIP,
		ConsulPort:      consulPort,
	}
}

func main() {
	go http.ListenAndServe(":36660", nil)

	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stdout)
	config := newConfig()

	svc, err := cl.CreateClient(ctx, logger, config)
	if err != nil {
		fmt.Println(err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/user/{data}", makeHandler(ctx, svc, logger))
	logger.Log(http.ListenAndServe(":8090", router))
}

func makeHandler(ctx context.Context, svc server.UserService, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := svc.FindOne(&server.UserFindOneRequest{&server.UserSearch{Id: mux.Vars(r)["data"]}})
		if err != nil {
			fmt.Println(err.Error())
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			logger.Log("resp_unmarshal", err.Error())
		}

		fmt.Fprintln(w, fmt.Sprintf("%s", respBytes))
	}
}
