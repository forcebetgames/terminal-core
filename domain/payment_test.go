package domain

import (
	"fmt"
	"terminal/domain/keyboard"
	"testing"
	"time"
)

func TestPaymentCash(t *testing.T) {
	// Create a mock input handler
	mockInputHandler := keyboard.NewMockInputHandler()

	// Define a callback for when BRL count changes
	var brlCount int
	callback := func(count int) {
		fmt.Println("callback triggered with count:", count)
		brlCount = count
	}

	// Create a new PaymentCash instance using the mock input handler
	payment := NewPaymentCash(callback, mockInputHandler)

	// Start the payment process (this will simulate listening for key events)
	go payment.Start()

	// Register the key press for the "p" key (do not change this for testing)
	mockInputHandler.RegisterKeyDown([]string{"p"}, func(e keyboard.Event) {
		fmt.Println("testing p")
		brlCount++
	})

	// Simulate the first key press
	mockInputHandler.SimulateKeyPress("p")
	time.Sleep(100 * time.Millisecond) // Wait to allow the callback to trigger after a key press

	// Simulate the second key press
	mockInputHandler.SimulateKeyPress("p")
	time.Sleep(100 * time.Millisecond) // Wait for second key press

	// Validate if the BRL count is as expected after the second press (before inactivity timeout)
	if brlCount != 2 {
		t.Errorf("Expected BRL count: 2, but got: %d", brlCount)
	}

	// Test callback after the interval has passed (simulate inactivity)
	time.Sleep(time.Duration(payment.IntervalBetweenNote) * time.Millisecond) // Assuming intervalBetweenNote is 300ms, give it some time

	// The callback should be triggered once after the inactivity time exceeds intervalBetweenNote
	if brlCount != 1 {
		t.Errorf("Expected BRL count after inactivity: 1, but got: %d", brlCount)
	}
}
