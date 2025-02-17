package main
// Author: c0d3-5t3w
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func checkRoot() {
	fmt.Println("Checking if the script is run as root...")
	if os.Geteuid() != 0 {
		fmt.Println("This script must be run as root. Please try again with 'sudo'.")
		os.Exit(1)
	}
	fmt.Println("Script is running as root.")
}

func checkHcxpcaptool() {
	fmt.Println("Checking if hcxpcapngtool is installed...")
	cmd := exec.Command("hcxpcapngtool", "--version")
	err := cmd.Run()
	if err != nil {
		fmt.Println("hcxpcapngtool is not installed. Installing hcxtools...")
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "hcxtool")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		err = installCmd.Run()
		if err != nil {
			fmt.Printf("Failed to install hcxtools: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("hcxtools installed successfully.")
	} else {
		fmt.Println("hcxpcapngtool is already installed.")
	}
}

func convertPcapToHashcat(pcapFile, outputFile string) bool {
	fmt.Printf("Converting %s to %s...\n", pcapFile, outputFile)
	cmd := exec.Command("hcxpcapngtool", "-z", outputFile, pcapFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error converting %s: %v\nOutput: %s\n", pcapFile, err, string(output))
		return false
	}
	fmt.Printf("Successfully converted %s to %s\n", pcapFile, outputFile)
	return true
}

func main() {
	fmt.Println("Starting pcap to hashcat conversion script...")
	checkRoot()
	checkHcxpcaptool()

	handshakes := "/home/pi/handshakes"
	hashcatables := "/home/pi/hashcatables"

	fmt.Printf("Checking if output directory %s exists...\n", hashcatables)
	if _, err := os.Stat(hashcatables); os.IsNotExist(err) {
		fmt.Printf("Output directory %s does not exist. Creating...\n", hashcatables)
		err = os.Mkdir(hashcatables, 0755)
		if err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", hashcatables, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully created directory %s\n", hashcatables)
	} else {
		fmt.Printf("Output directory %s already exists.\n", hashcatables)
	}

	fmt.Printf("Walking through the directory %s to find .pcap files...\n", handshakes)
	err := filepath.Walk(handshakes, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing the path %s: %v\n", path, err)
			return err
		}
		if strings.HasSuffix(info.Name(), ".pcap") {
			fmt.Printf("Found .pcap file: %s\n", path)
			outputFile := filepath.Join(hashcatables, strings.Replace(info.Name(), ".pcap", ".hc22000", 1))
			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				if convertPcapToHashcat(path, outputFile) {
					fmt.Printf("Converted %s to %s\n", path, outputFile)
				} else {
					fmt.Printf("Failed to convert %s\n", path)
					time.Sleep(100 * time.Millisecond) 
				}
			} else {
				fmt.Printf("Skipping %s, already converted\n", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", handshakes, err)
		os.Exit(1)
	}
	fmt.Println("Finished pcap to hashcat conversion script.")
}
