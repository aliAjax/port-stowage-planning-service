package dispatch

import (
	"context"
	"testing"
	"time"

	"github.com/example/port-stowage-planner/internal/domain"
)

func TestDispatchStaleRejected(t *testing.T) {
	s := NewSimulator()
	now := time.Now()
	instr := domain.WorkInstruction{ID: "i1", PlanID: "p1", ContainerID: "c1", Sequence: 1}
	fresh := Lease(instr, now, time.Minute)
	if err := s.Dispatch(context.Background(), fresh); err != nil {
		t.Fatal(err)
	}
	expired := Lease(instr, now.Add(-2*time.Minute), time.Minute)
	if err := s.Dispatch(context.Background(), expired); err == nil {
		t.Fatalf("re-dispatch with an expired lease must be rejected")
	}
}

func TestAckValidatesStored(t *testing.T) {
	s := NewSimulator()
	now := time.Now()
	instr := domain.WorkInstruction{ID: "i1", PlanID: "p1", ContainerID: "c1", Sequence: 1}
	first := Lease(instr, now, time.Minute)
	if err := s.Dispatch(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	second := Lease(instr, now.Add(5*time.Second), time.Minute)
	if err := s.Dispatch(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	err := s.Ack(first, domain.Ack{InstructionID: "i1", LeaseToken: first.LeaseToken, Result: "done", At: now.Add(6 * time.Second)})
	if err == nil {
		t.Fatalf("ack with a stale token must be rejected")
	}
	if got := s.Seen["i1"]; got.LeaseToken != second.LeaseToken {
		t.Fatalf("stored lease must remain the latest dispatch")
	}
}

func TestAckKeepsStoredFields(t *testing.T) {
	s := NewSimulator()
	now := time.Now()
	instr := domain.WorkInstruction{ID: "i1", PlanID: "p1", ContainerID: "c1", Sequence: 1}
	dispatched := Lease(instr, now, time.Minute)
	if err := s.Dispatch(context.Background(), dispatched); err != nil {
		t.Fatal(err)
	}
	staleCopy := dispatched
	staleCopy.ContainerID = "tampered"
	if err := s.Ack(staleCopy, domain.Ack{InstructionID: "i1", LeaseToken: dispatched.LeaseToken, Result: "done", At: now.Add(time.Second)}); err != nil {
		t.Fatal(err)
	}
	if got := s.Seen["i1"]; got.ContainerID != "c1" {
		t.Fatalf("ack must not overwrite stored fields with a caller copy, got %s", got.ContainerID)
	}
}

func TestLeaseEmptyTokenInvalid(t *testing.T) {
	now := time.Now()
	i := domain.WorkInstruction{ID: "i1", LeaseUntil: now.Add(time.Minute)}
	if ValidLease(i, now) {
		t.Fatalf("lease without a token must be invalid")
	}
}
