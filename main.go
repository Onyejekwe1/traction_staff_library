package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"traction_staff_library/config"

	_ "github.com/lib/pq"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/sakto"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Laporta1."
	dbname   = "traction_library"
	sslmode  = "disable"
)

var CurrentLocalTime = sakto.GetCurDT(time.Now(), "Africa/Lagos")

func main() {
	// Tracking the time we started the server
	os.Setenv("TZ", config.SiteTimeZone)
	fmt.Println("Starting the web server at ", CurrentLocalTime)

	a := DBConnect{}

	a.Initialize(host, port, user, password, dbname, sslmode)
	// Organizing our db connection string
	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	//
	//db, err := sql.Open("postgres", psqlInfo)
	//if err != nil {
	//	itrlog.Error("Error occurred", err)
	//	panic(err)
	//}
	//defer db.Close()
	//
	//// Trying to connect to our DB
	//err = db.Ping()
	//if err != nil {
	//	itrlog.Error("Error occurred", err)
	//	panic(err)
	//}

	var dir string
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()

	// Create cross-site request forgery (CSRF) protection in every http requests.
	// 32-byte-long-auth-key []string{config.SiteDomainName}
	csrfMiddleware := csrf.Protect(
		[]byte(config.SecretKeyCORS),
		csrf.TrustedOrigins([]string{config.SiteDomainName}),
	)

	// This is related to the CORS config to allow all origins []string{"*"} or specify only allowed IP or hostname.
	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{config.SiteDomainName}),
	)

	r.Use(cors)
	r.Use(csrfMiddleware)
	r.Use(loggingMiddleware)
	r.Use(mux.CORSMethodMiddleware(r))

	// Will serve our static files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Addr: "127.0.0.1:8081",
		// Good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a ago routine so it doesn't block
	go func() {
		msg := `Web server started at `
		fmt.Println(msg, CurrentLocalTime)
		itrlog.Info("Web server started at ", CurrentLocalTime)
		if err := srv.ListenAndServe(); err != nil {
			itrlog.Error(err)
		}
	}() // Note the parentheses calls the function.

	// Buffered channels = queue
	c := make(chan os.Signal, 1) // Queue with a capacity of 1.

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	fmt.Println("Shutdown web server at " + CurrentLocalTime.String())
	itrlog.Warn("Server has been shutdown at ", CurrentLocalTime.String())
	os.Exit(0)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		req := "IP:" + sakto.GetIP(r) + ":" + r.RequestURI + ":" + CurrentLocalTime.String()
		fmt.Println(req)
		itrlog.Info(req)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
