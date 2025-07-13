package uuid

import (
	"fmt"
	"testing"
)

func TestNewV4(t *testing.T) {
	// 1. Generate single UUID v4
	fmt.Println("Single UUID v4:")
	uuidV4, err := NewV4()
	if err != nil {
		fmt.Printf("Generation error: %v\n", err)
	} else {
		fmt.Printf("Standard format: %s\n", uuidV4)
		fmt.Printf("No hyphens: %s\n", uuidV4.NoHyphenString())
		fmt.Printf("Version: %d, Variant: %d, Valid: %t\n\n",
			uuidV4.Version(), uuidV4.Variant(), uuidV4.IsValid())
	}

}

func TestNewV4Batch(t *testing.T) {
	// 2. Batch generate UUID v4
	fmt.Println("Batch generated UUID v4:")
	batch, err := NewV4Batch(20)
	if err != nil {
		fmt.Printf("Batch generation error: %v\n", err)
	} else {
		for i, u := range batch {
			fmt.Printf("%d. %s (no hyphens: %s)\n", i+1, u, u.NoHyphenString())
		}
	}
}
