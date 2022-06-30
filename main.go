package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"task/internal/app/handlers"
	json2 "task/internal/app/repositories/json"

	//	repomemory "task/internal/app/repositories/memory"
	"task/internal/app/services/configmanager"
	"task/internal/app/services/mailsendingmanager"
	"task/internal/app/services/voteeventmanager"
	"task/internal/app/services/votelinkmanager"
	"task/internal/app/services/votemanager"
	"task/internal/app/store/json"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.yml", "path to config file")
}

func main() {
	flag.Parse()
	config := configmanager.NewConfig()
	err := config.Init(configPath)
	if err != nil {
		log.Fatal(err)
	}

	//	taskRepo := repomemory.NewTaskRepo()
	taskRepo := json2.NewTaskRepo(config)
	//	store := memory.NewStore(taskRepo)
	store := json.NewStore(taskRepo, config)
	vlm := votelinkmanager.NewEncryptVoteLinkManager(config)
	vem := voteeventmanager.NewVoteEventManager()
	msm := mailsendingmanager.NewDummyMailSendingManager()

	vm := votemanager.NewVoteManager(store, vlm, vem, msm)

	taskCtrl := handlers.NewTaskController(store, vm)

	mw := handlers.NewMiddleware()

	router := mux.NewRouter()

	router.HandleFunc("/task", taskCtrl.Create).Methods("POST")
	router.HandleFunc("/task/{id:[0-9]+}", taskCtrl.Get).Methods("GET")
	router.HandleFunc("/tasks", taskCtrl.GetAll).Methods("GET")
	router.HandleFunc("/task", taskCtrl.Update).Methods("PUT")
	router.HandleFunc("/task/{id:[0-9]+}", taskCtrl.Delete).Methods("DELETE")

	router.HandleFunc("/vote/{vote_link}", taskCtrl.Vote).Methods("GET")

	router.Use(mw.Logging)

	err = http.ListenAndServe(config.BindAddr, router)
	log.Fatal(err)
}
