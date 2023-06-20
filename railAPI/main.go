package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"

	"github.com/atanda0x/Metro-Rail-Api/dbutils"
)

// DB Driver visible to whole program
var DB *sql.DB

// TrainResources is the model for holding rail information
type TrainResources struct {
	ID              int
	DriverName      string
	Operatingstatus bool
}

// StationResource holds information about location
type StationResource struct {
	ID          int
	Name        string
	OpeningTime time.Time
	ClosingTime time.Time
}

// ScheduleRescourse links both trains and station
type ScheduleRescourse struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime time.Time
}

// Get http://localhost:8000/v1/trains/1
func (t TrainResources) getTrain(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("trian-id")
	err := DB.QueryRow("select ID, DRIVER_NAME, OPERATING_STATUS FROM train where id=?", id).Scan(&t.ID, &t.DriverName, &t.Operatingstatus)
	if err != nil {
		log.Println(err)
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusNotFound, "Train culd not be found.")
	} else {
		res.WriteEntity(t)
	}
}

// POST http://localhost:9000/v1/trains
func (t TrainResources) createTrain(req *restful.Request, res *restful.Response) {
	log.Println(req.Request.Body)
	decoder := json.NewDecoder(req.Request.Body)
	var b TrainResources
	err := decoder.Decode(&b)
	log.Println(b.DriverName, b.Operatingstatus)
	// Error handling is obvious here. so omitting....
	statement, _ := DB.Prepare("insert into train (DRIVER_NAME, OPERATING_STATUS) values (?, ?)")
	result, err := statement.Exec(b.DriverName, b.Operatingstatus)
	if err == nil {
		newID, _ := result.LastInsertId()
		b.ID = int(newID)
		res.WriteHeaderAndEntity(http.StatusCreated, b)
	} else {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// Delete http://localhost:9000/v1/trains/1
func (t TrainResources) removeTrain(req *restful.Request, res *restful.Response) {
	id := req.PathParameter("train-id")
	statement, _ := DB.Prepare("delete from train where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		res.WriteHeader(http.StatusOK)
	} else {
		res.AddHeader("Content-Type", "text/plain")
		res.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// Register addds paths and routes to container
func (t *TrainResources) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/trains").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON) //You can specify this per route as well
	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))
	container.Add(ws)
}

func main() {
	db, err := sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("Driver creation failed!!!")
	}

	// Create tables
	dbutils.Initialise(db)
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	t := TrainResources{}
	t.Register(wsContainer)
	log.Printf("start listening on localhost:9000")
	server := &http.Server{
		Addr:    ":9000",
		Handler: wsContainer,
	}
	log.Fatal(server.ListenAndServe())
}
