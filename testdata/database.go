package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"minik8s/util/log"
	"sync"
)

type BuyRequest struct {
	Balance int `json:"balance"`
}

type CancelRequest struct {
	UUID string `json:"uuid"`
}

type order struct {
	Change int    `json:"change"`
	Status string `json:"status"`
	UUID   string `json:"uuid"`
}

var stockTrain int
var stockFlight int
var stockHotel int

var TrainTickets []string
var FlightTickets []string
var HotelOrder []string

var TrainReserved []bool
var FlightReserved []bool
var HotelReserved []bool

var mutexTrain sync.Mutex
var mutexFlight sync.Mutex
var mutexHotel sync.Mutex

func init() {
	stockTrain = 10
	stockFlight = 10
	stockHotel = 10
	TrainTickets = make([]string, 10)
	FlightTickets = make([]string, 10)
	HotelOrder = make([]string, 10)
	TrainReserved = make([]bool, 10)
	FlightReserved = make([]bool, 10)
	HotelReserved = make([]bool, 10)
	for i := 0; i < 10; i++ {
		TrainTickets[i] = uuid.NewString()
		FlightTickets[i] = uuid.NewString()
		HotelOrder[i] = uuid.NewString()
		TrainReserved[i] = false
		FlightReserved[i] = false
		HotelReserved[i] = false
	}
}

func BuyTrainTicket(context *gin.Context) {
	mutexTrain.Lock()
	defer mutexTrain.Unlock()
	var request BuyRequest
	if err := context.ShouldBind(&request); err != nil {
		reply := order{
			Change: -1,
			Status: "Failed",
			UUID:   "",
		}
		context.JSON(200, reply)
		return
	}
	if stockTrain == 0 || request.Balance < 100 {
		reply := order{
			Change: request.Balance,
			Status: "Failed",
			UUID:   "",
		}
		context.JSON(200, reply)
	} else {
		for index, reserved := range TrainReserved {
			if !reserved {
				TrainReserved[index] = true
				stockTrain--
				reply := order{
					Change: request.Balance - 100,
					Status: "Succeeded",
					UUID:   TrainTickets[index],
				}
				context.JSON(200, reply)
				return
			}
		}
	}
}

func ReserveFlight(context *gin.Context) {
	mutexFlight.Lock()
	defer mutexFlight.Unlock()
	var request BuyRequest
	if err := context.ShouldBind(&request); err != nil {
		reply := order{
			Change: -1,
			Status: "Failed",
			UUID:   "",
		}
		context.JSON(200, reply)
		return
	}
	if stockFlight == 0 || request.Balance < 800 {
		reply := order{
			Change: request.Balance,
			Status: "Failed",
			UUID:   "",
		}
		context.JSON(200, reply)
	} else {
		for index, reserved := range FlightReserved {
			if !reserved {
				FlightReserved[index] = true
				stockFlight--
				reply := order{
					Change: request.Balance - 800,
					Status: "Succeeded",
					UUID:   FlightTickets[index],
				}
				context.JSON(200, reply)
				return
			}
		}
	}
}

func ReserveHotel(context *gin.Context) {
	mutexHotel.Lock()
	defer mutexHotel.Unlock()
	var request BuyRequest
	if err := context.ShouldBind(&request); err != nil {
		reply := order{
			Change: -1,
			Status: "Failed",
			UUID:   "",
		}
		context.JSON(200, reply)
		return
	}
	if stockHotel == 0 || request.Balance < 400 {
		reply := order{
			Change: request.Balance,
			Status: "Failed",
			UUID:   "",
		}
		context.JSON(200, reply)
	} else {
		for index, reserved := range HotelReserved {
			if !reserved {
				HotelReserved[index] = true
				stockHotel--
				reply := order{
					Change: request.Balance - 400,
					Status: "Succeeded",
					UUID:   HotelOrder[index],
				}
				context.JSON(200, reply)
				return
			}
		}
	}
}

func CancelFlight(context *gin.Context) {
	mutexFlight.Lock()
	defer mutexFlight.Unlock()
	var request CancelRequest
	if err := context.ShouldBind(&request); err != nil {
		return
	}
	for index, _ := range FlightTickets {
		if request.UUID == FlightTickets[index] {
			FlightReserved[index] = false
			stockFlight++
			return
		}
	}
}

func CancelTrain(context *gin.Context) {
	mutexTrain.Lock()
	defer mutexTrain.Unlock()
	var request CancelRequest
	if err := context.ShouldBind(&request); err != nil {
		return
	}
	for index, _ := range TrainTickets {
		if request.UUID == TrainTickets[index] {
			TrainReserved[index] = false
			stockTrain++
			return
		}
	}
}
func main() {
	server := gin.Default()
	server.POST("/buytrainticket", BuyTrainTicket)
	server.POST("/reserveflight", ReserveFlight)
	server.POST("/reservehotel", ReserveHotel)
	server.POST("/cancelflight", CancelFlight)
	server.POST("/canceltrain", CancelTrain)
	err := server.Run(":19293")
	if err != nil {
		log.Fatal("server start error", err)
		return
	}
}
