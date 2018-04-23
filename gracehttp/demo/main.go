package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"jvole.com/dsp/gracehttp"
)

/*
curl 'http://localhost:8080/sleep/?duration=20s'
kill -SIGUSR2 $pid (-SIGTERM)
curl 'http://localhost:8080/sleep/?duration=1s'
*/
func main() {

	http.HandleFunc("/sleep/", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		time.Sleep(duration)

		fmt.Fprintf(
			w,
			"started at %s slept for %d nanoseconds from pid %d.\n",
			time.Now(),
			duration.Nanoseconds(),
			os.Getpid(),
		)
	})

	log.Println(fmt.Sprintf("Serving :8080 with pid %d.", os.Getpid()))

	gracehttp.ListenAndServe(":8080", nil)

	log.Println("Server stoped.")
}
