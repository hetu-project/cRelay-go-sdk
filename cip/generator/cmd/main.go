package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nbd-wtf/go-nostr/cip/generator"
)

func main() {
	// Define command line flags
	inputFile := flag.String("input", "", "Path to the CIP definition JSON file")
	flag.Parse()

	// Check if input file is provided
	if *inputFile == "" {
		fmt.Println("Error: Input file is required")
		flag.Usage()
		os.Exit(1)
	}

	// Read and parse the input file
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	var def generator.CIPDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		log.Fatalf("Failed to parse CIP definition: %v", err)
	}

	// Generate the CIP implementation
	if err := generator.GenerateCIP(def); err != nil {
		log.Fatalf("Failed to generate CIP: %v", err)
	}

	fmt.Printf("Successfully generated CIP implementation in cip/%s\n", def.Package)
}
