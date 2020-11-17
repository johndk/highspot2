package main

import (
	"flag"
	"fmt"
	"highspot2/brokers"
	"log"
	"os"
)

type CommandLine struct {
	Host   string
	Output string
	Help   bool
}

var cmdline = CommandLine{}

// The main program.
func main() {
	// Parse the command line arguments.
	flag.Parse()

	if cmdline.Help {
		// Print usage and exit 0.
		flag.Usage()
		os.Exit(0)
	}

	//
	// Create an ingester and execute the take-home exercise
	//

	ingester := brokers.NewIngestor(cmdline.Host, cmdline.Output)
	err := ingester.DoIngest()
	if err != nil {
		log.Fatalf("Error encountered. %v", err)
	}

	log.Printf("The output file %v was successfully created.", cmdline.Output)
}

// Initialize the command line arguments. Print usage highspot -h.
func init() {
	flag.StringVar(&cmdline.Host, "u", "http://localhost:8080/data", "The input file host.")
	flag.StringVar(&cmdline.Output, "o", "output.json", "The output file path.")
	flag.BoolVar(&cmdline.Help, "h", false, "Print the help text.")
	flag.Usage = printUsage
}

func printUsage() {
	fmt.Print("The Highspot take-home coding exercise.\n\n")
	fmt.Print("Usage: highspot [arguments]\n\n")
	fmt.Print("The arguments are:\n\n")
	flag.PrintDefaults()
}
