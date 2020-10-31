package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/scayle/common-go"
)

func healthPort(defaultPort int) int {
	p := os.Getenv("PRODUCT_HEALTH_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return defaultPort
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		panic("invalid format for the environment variable PRODUCT_HEALTH_PORT")
	}
	return port
}

func mongoHealth() bool {
	response, err := http.Get("http://localhost:27017")
	if err != nil {
		fmt.Println("err", err.Error())
		return false
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return false
	}
	responseBody := buf.String()

	if strings.HasPrefix(responseBody, "It looks like you are trying to access MongoDB over HTTP on the native driver port.") {
		return true
	}
	return false
}

func main() {
	stop := make(chan bool)

	counter := 0

	common.RegisterConsulService("mongodb",
		common.WithDefaultPort(27017),
		common.WithRegistrationModifier(func(registration *api.AgentServiceRegistration) {
			defaultPort := 8101

			// setup simple health detection using a small webserver
			registration.Check = new(api.AgentServiceCheck)
			registration.Check.HTTP = fmt.Sprintf("http://%s:%d/healthcheck", registration.Address, healthPort(defaultPort))
			registration.Check.Interval = "5s"
			registration.Check.Timeout = "3s"

			fmt.Println(registration.Address, healthPort(defaultPort))
			http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Healthcheck")
				if !mongoHealth() {
					if counter > 5 {
						// The http handler would recover from a panic. -> use of os.Exit
						fmt.Println("could not access mongo db for several times")
						os.Exit(1)
					}

					counter++
					w.WriteHeader(http.StatusServiceUnavailable)
					_, err := w.Write([]byte{})
					if err != nil {
						// The http handler would recover from a panic. -> use of os.Exit
						fmt.Println(err)
						os.Exit(1)
					}
					return
				}
				_, err := fmt.Fprintf(w, `Mongodb is alive!`)
				if err != nil {
					// The http handler would recover from a panic. -> use of os.Exit
					fmt.Println(err)
					os.Exit(1)
				}
			})

			go func() {
				err := http.ListenAndServe(fmt.Sprintf(":%d", healthPort(defaultPort)), nil)
				log.Printf("healthcheck webserver failed %v", err)
				close(stop)
			}()
		}))
	<-stop
}
