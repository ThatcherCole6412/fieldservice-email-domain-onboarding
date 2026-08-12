package main

import (
	"fmt"
	"os"
)

func main() {
	order := WorkOrder{ID: "WO-1842", Customer: "Northwind Market", CustomerEmail: os.Getenv("DEMO_EMAIL_TO"), PhotoCount: 3, DispatchStatus: "completed"}
	message, needed := FollowUpMessage(order)
	if !needed {
		fmt.Println("no follow-up required")
		return
	}
	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	domain := os.Getenv("DEMO_EMAIL_DOMAIN")
	if domain == "" {
		panic("DEMO_EMAIL_DOMAIN is required")
	}
	verification, err := client.VerifyDomain(domain)
	if err != nil {
		panic(err)
	}
	fmt.Println("domain verification:", verification.Status)
	result, err := client.SendFollowUp(order.CustomerEmail, "Service completed: "+order.ID, "<p>"+message+"</p>", "work-order-"+order.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("follow-up sent:", result.MessageID)
}
