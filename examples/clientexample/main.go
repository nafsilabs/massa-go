package main

import (
	"github.com/nafsilabs/massa-go/client"
)

func main() {
	c := client.NewClient(true)
	slot, err := client.NextSlot(c)
	if err != nil {
		panic(err)
	}
	println("Next slot period:", slot)

	address := "AU1y3oYTgK8RGzLWFVAGL3JxHLdTVKxBPHrwbo2Kj8a2CSbeMug"
	balance, err := client.FetchBalance(c, address)
	if err != nil {
		panic(err)
	}
	println("Address:", address)
	println("Candidate balance:", balance.Candidate.String())
	println("Final balance:", balance.Final.String())

	//listen to events
	// startSlot := client.NewSlot(int(slot)-1, 0)
	// endSlot := client.NewSlot(int(slot)+10, 0)
	// events, err := client.ListenEvents(c, startSlot, endSlot, nil, nil, false)
	// if err != nil {
	// 	panic(err)
	// }
	// println("Events found:", len(events))
	// for _, event := range events {
	// 	println("Event data:", event.Data)
	// }
}
