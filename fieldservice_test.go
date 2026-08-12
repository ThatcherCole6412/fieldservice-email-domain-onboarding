package main

import "testing"

func TestFollowUpDecisionRequiresCompletedPhotoEvidence(t *testing.T) {
	base := WorkOrder{ID: "WO-1", Customer: "Aster", CustomerEmail: "ops@example.com", PhotoCount: 2, DispatchStatus: "completed"}
	if message, ok := FollowUpMessage(base); !ok || message == "" {
		t.Fatal("completed work with photos should produce a follow-up")
	}
	base.DispatchStatus = "en_route"
	if _, ok := FollowUpMessage(base); ok {
		t.Fatal("unfinished dispatch should not produce a follow-up")
	}
}
