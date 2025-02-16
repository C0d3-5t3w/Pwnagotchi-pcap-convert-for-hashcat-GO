package main
// Author: c0d3-5t3w
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("This script must be run as root. Please try again with 'sudo'.")
		os.Exit(1)
	}
}

func checkHcxpcaptool() {
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
	}
}

func convertPcapToHashcat(pcapFile, outputFile string) bool {
	cmd := exec.Command("hcxpcapngtool", "-z", outputFile, pcapFile)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error converting %s: %v\n", pcapFile, err)
		return false
	}
	return true
}

func main() {
	checkRoot()
	checkHcxpcaptool()

	handshakes := "/home/pi/handshakes"
	hashcatables := "/home/pi/hashcatables"

	if _, err := os.Stat(hashcatables); os.IsNotExist(err) {
		err = os.Mkdir(hashcatables, 0755)
		if err != nil {
			fmt.Printf("Failed to create directory %s: %v\n", hashcatables, err)
			os.Exit(1)
		}
	}

	err := filepath.Walk(handshakes, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".pcap") {
			outputFile := filepath.Join(hashcatables, strings.Replace(info.Name(), ".pcap", ".hc22000", 1))
			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				if convertPcapToHashcat(path, outputFile) {
					fmt.Printf("Converted %s to %s\n", path, outputFile)
				} else {
					fmt.Printf("Failed to convert %s\n", path)
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
}
