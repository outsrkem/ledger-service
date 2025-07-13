package uuid

import (
	"fmt"
	"testing"
)

func TestNewV5(t *testing.T) {
	fmt.Println("\nUUID v5 examples:")
	domains := []string{"example.com", "google.com", "github.com"}
	for i, domain := range domains {
		uuidV5 := NewV5(NamespaceDNS, []byte(domain))
		fmt.Printf("%d. %s (%s)\n", i+1, domain, uuidV5)
		fmt.Printf("  Version: %d, Variant: %d, Valid: %t\n",
			uuidV5.Version(), uuidV5.Variant(), uuidV5.IsValid())
	}

}

func TestUUID_Namespace(t *testing.T) {
	// 4. Custom namespace example
	fmt.Println("\nCustom namespace example:")
	customNamespace, _ := NewV4() // Use v4 as custom namespace
	name := []byte("custom-identifier")
	customUUID := NewV5(customNamespace, name)

	fmt.Printf("Namespace: %s\n", customNamespace)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Generated UUID: %s\n", customUUID)
}
func TestNewV52(t *testing.T) {
	for i := 0; i < 5; i++ {
		customNamespace, _ := NewV4()
		name := []byte(customNamespace.String())
		customUUID := NewV5(customNamespace, name)
		fmt.Printf("Generated UUID: %s\n", customUUID.NoHyphenString())
	}
}
