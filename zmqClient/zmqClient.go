package zmqClient

import (
	"GoPack/fileHandling"
	"context"
	"encoding/json"
	"fmt"
	zmq "github.com/go-zeromq/zmq4"
	"log"
)

const microAexporterPort = "5555"
const microBweatherPort = "5556"
const microClistSortPort = "5557"
const microDpackedPercentPort = "5558"

// communicateWithMicroservice will open a connection to a microservice via
// zeroMQ, will send a request, wait for the reply, close the connection, and
// then return the response message
func communicateWithMicroservice(port string, m zmq.Msg) zmq.Msg {

	// set up a context & socket
	ctx := context.Background()
	//socket := zmq.NewReq(ctx, zmq.WithDialerRetry(time.Second))
	socket := zmq.NewReq(ctx)
	// make sure to close when we are done
	defer socket.Close()

	// connect to the microservice
	address := fmt.Sprint("tcp://localhost:" + port)
	if err := socket.Dial(address); err != nil {
		//log.Fatal(err)
		var empty []byte
		return zmq.NewMsg(empty)
	}

	// send message to the microservice
	if err := socket.Send(m); err != nil {
		log.Fatal(err)
	}

	// get response from microservice
	response, err := socket.Recv()
	if err != nil {
		log.Fatal(err)
	}

	return response
}

// Microservice A ------------------------------------------

const EMAILSERVER = "smtp.gmail.com"
const EMAILPORT = "587"

type ExportRequest struct {
	List     fileHandling.PackingList `json:"packing_list"`
	Username string                   `json:"username"`
	Password string                   `json:"password"`
	Server   string                   `json:"server"`
	Port     string                   `json:"port"`
}

func SendExportRequest(list fileHandling.PackingList, username string, password string) string {

	var newRequest ExportRequest

	newRequest.List = list
	newRequest.Username = username
	newRequest.Password = password
	newRequest.Server = EMAILSERVER
	newRequest.Port = EMAILPORT

	// package request
	requestAsJSON, packingErr := json.Marshal(newRequest)
	if packingErr != nil {
		log.Fatal(packingErr)
	}

	// form message
	message := zmq.NewMsg(requestAsJSON)

	// call function to send message & receive message reply
	response := communicateWithMicroservice(microAexporterPort, message)

	// receives a response as a string
	return response.String()
}

// Microservice B ------------------------------------------
//
//type WeatherRequest struct {
//	Location   string `json:"Location"`
//	DepartDate string `json:"DepartDate"`
//	ReturnDate string `json:"ReturnDate"`
//}
//
//type WeatherResponse struct {
//	Data     map[string]ForecastData `json:"Data"`
//	ErrorMsg string                  `json:"ErrorMsg"`
//}
//
//type ForecastData struct {
//	TempMax             float64
//	TempMin             float64
//	PrecipitationChance float64
//}

//func SendWeatherRequest(location string, departDate string, returnDate string) WeatherResponse {
//
//	var newRequest WeatherRequest
//	newRequest.Location = location
//	newRequest.DepartDate = departDate
//	newRequest.ReturnDate = returnDate
//
//	// package request
//	requestAsJSON, packingErr := json.Marshal(newRequest)
//	if packingErr != nil {
//		log.Fatal(packingErr)
//	}
//
//	// form message
//	message := zmq.NewMsg(requestAsJSON)
//
//	// call function to send message & receive message reply
//	response := communicateWithMicroservice(microBweatherPort, message)
//
//	// unpack response
//	var forecast WeatherResponse
//	unpackErr := json.Unmarshal(response.Bytes(), &forecast)
//	if unpackErr != nil {
//		log.Fatal(unpackErr)
//	}
//
//	return forecast
//}

// Microservice C -------------------------------------------

//type SortRequest struct {
//	SortOn string                   `json:"SortOn"`
//	List   fileHandling.PackingList `json:"List"`
//}
//
//func SendListSortRequest(list fileHandling.PackingList, sortOn string) fileHandling.PackingList {
//
//	var newRequest SortRequest
//	newRequest.List = list
//	newRequest.SortOn = sortOn
//
//	// package request
//	requestAsJSON, packingErr := json.Marshal(newRequest)
//	if packingErr != nil {
//		log.Fatal(packingErr)
//	}
//
//	// form message
//	message := zmq.NewMsg(requestAsJSON)
//
//	// call function to send message & recieve message reply
//	response := communicateWithMicroservice(microClistSortPort, message)
//
//	// unpack response
//	var sortedList fileHandling.PackingList
//	unpackErr := json.Unmarshal(response.Bytes(), &sortedList)
//	if unpackErr != nil {
//		log.Fatal(unpackErr)
//	}
//
//	return sortedList
//}

// Microservice D -------------------------------------------
//
//type PackedPercentageResults struct {
//	TotalItems       int            `json:"TotalItems"`
//	PackedItems      int            `json:"PackedItems"`
//	TotalInLocation  map[string]int `json:"TotalInLocation"`
//	PackedByLocation map[string]int `json:"PackedByLocation"`
//}
//
//func SendPackedPercentRequest(packingList fileHandling.PackingList) PackedPercentageResults {
//
//	// package packing list
//	listAsJSON, packingErr := json.Marshal(packingList)
//	if packingErr != nil {
//		log.Fatal(packingErr)
//	}
//
//	// form message
//	message := zmq.NewMsg(listAsJSON)
//
//	// call function to send message & receive message reply
//	response := communicateWithMicroservice(microDpackedPercentPort, message)
//
//	// unpack response
//
//	var packedPercentage PackedPercentageResults
//	unpackErr := json.Unmarshal(response.Bytes(), &packedPercentage)
//	if unpackErr != nil {
//		//log.Fatal(unpackErr)
//		// will likely error out if too many requests to the zmq service
//		return PackedPercentageResults{}
//	}
//
//	return packedPercentage
//}
