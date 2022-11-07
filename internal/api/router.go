package api

import (
	"challenge/internal/cache"
	"challenge/internal/entity"
	"challenge/internal/logger"	
	"challenge/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Router struct {
	r          *mux.Router
	repository repository.Repository
	cache *cache.Cache
}

//cria um novo router
func New(repo repository.Repository, cache *cache.Cache) *Router {
	return &Router{ 
		r:          mux.NewRouter(),
		repository: repo,
		cache: cache,
	}
}

//função handler que trabalha com os metodos get e post http para URL que termina em "/tasks"
// o metodo get lista tarefas de acordo com o que foi proposto no navegador
//o metodo post cria tarefa de acordo com o que for indicado
//ambos metodos usam json para encodar e decodar mensagens 
func (ro *Router) HandlerList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		comp := r.FormValue("completed")
		var task []entity.Task
		var err error
		if comp == "" {
			task, err = ro.repository.ListTask()
		} else {
			if comp == "true" {
				task, err = ro.repository.ListComp("yes")

			}
			if comp == "false" {
				task, err = ro.repository.ListComp("no")
			}
		}

		if err != nil {
			logger.Log().Warn("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(task)
	case "POST":
		var task entity.Task
		json.NewDecoder(r.Body).Decode(&task)

		id, err := ro.repository.NewTask(task)
		if err != nil {
			logger.Log().Warn("%v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		task.ID = int(id)
		jt, err := json.Marshal(task)
		if err != nil {
			logger.Log().Warn("%v", err)
		}

		err = ro.cache.Set(strconv.Itoa(task.ID), string(jt))
		if err != nil{
			logger.Log().Warn("coudn't set task: %v", err)
		}

		json.NewEncoder(w).Encode(id)

	}
}

//função handler que trabalha com os metodos get e put http para URL que termina em "/tasks/id". 
//O metodo get expoe tarefa que for solicitada no navegador.
//O metodo put atualiza tarefa que for solicitada.
//Ambos metodos usam json para encodar e decodar mensagens 
func (ro *Router) HandlerGetTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
		if err != nil {
			logger.Log().Warn("%v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		tk, err := ro.cache.Get(strconv.Itoa(id))
		if err == nil {
			json.NewEncoder(w).Encode(tk)
			return
		}

		task, err := ro.repository.GetTask(id)
		if err != nil {
			logger.Log().Warn("%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return 
		}

		jt, err := json.Marshal(task)
		if err != nil {
			logger.Log().Warn("%v", err)
		}

		err = ro.cache.Set(strconv.Itoa(id), string(jt))
		if err != nil{
			logger.Log().Warn("coudn't set task: %v", err)
			return
		}
		
		json.NewEncoder(w).Encode(task)

	case "PUT":

		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
		if err != nil {
			logger.Log().Warn("%v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		
		err = ro.repository.UpdateTask(id)
		if err != nil {
			logger.Log().Warn("%v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		tk, err := ro.cache.Get(strconv.Itoa(id))
		if err == nil {
			jt, err := json.Marshal(tk)
			if err != nil {
				logger.Log().Warn("%v", err)
			}
			
			err = ro.cache.Set(strconv.Itoa(id), string(jt))
			if err != nil {
				err = ro.cache.Del(strconv.Itoa(id))
				if err != nil {
					logger.Log().Warn("%v", err)
				}
			}
		}

		if err != nil {
			logger.Log().Warn("%v", err)
		}

		w.WriteHeader(http.StatusOK) //retorna 200
	}
}
