package main

type WorkOrder struct {
	ID, Customer, CustomerEmail string
	PhotoCount                  int
	DispatchStatus              string
}

func FollowUpNeeded(order WorkOrder) bool {
	return order.DispatchStatus == "completed" && order.PhotoCount > 0 && order.CustomerEmail != ""
}

func FollowUpMessage(order WorkOrder) (string, bool) {
	if !FollowUpNeeded(order) {
		return "", false
	}
	return "Work order " + order.ID + " is complete. " + order.Customer + ", please review the service photos.", true
}
