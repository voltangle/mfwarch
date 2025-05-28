package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"github.com/davecgh/go-spew/spew"
	cp "github.com/otiai10/copy"
)

func main() {
    fmt.Println("Hello from mfwarch!")

	_, err := os.Stat("firmwares")
	if os.IsNotExist(err) {
		os.Mkdir("firmwares", 0770)
		list := begodeDownloadAll("firmwares", []string{""})
		file, err := os.Create("firmwares/firmwares.json")
		if err != nil {
			fmt.Printf("mfwarch: unable to create file: %s\n", err)
			os.Exit(1)
		}

		fullList := FirmwareList {
			LastUpdateTime: time.Now(),	
			ChangesSinceLastUpdate: 0, // TODO: implement
			Begode: list,
		}

		fwlist, err := json.MarshalIndent(fullList, "", " ")
		if err != nil {
			fmt.Printf("mfwarch: unable to serialize firmware list: %s\n", err)
			os.Exit(1)
		}

		_, err = file.Write(fwlist)
		if err != nil {
			fmt.Printf("mfwarch: unable to save firmware: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	fwList, err := getLocalFirmwareList("firmwares/firmwares.json")
	if err != nil {
		fmt.Printf("mfwarch: unable to get local firmware list: %s\n", err)
		os.Exit(1)
	}
	spew.Dump(fwList)

	err = cp.Copy("firmwares", "firmwares2")
	if err != nil {
		fmt.Printf("mfwarch: unable to copy firmwares: %s\n", err)
	}

	list := begodeDownloadAll("firmwares2", []string{""})

	file, err := os.Create("firmwares/firmwares.json")
	if err != nil {
		fmt.Printf("mfwarch: unable to create file: %s\n", err)
		os.Exit(1)
	}

	fullList := FirmwareList {
		LastUpdateTime: time.Now(),	
		ChangesSinceLastUpdate: 0, // TODO: implement
		Begode: list,
	}

	fwlist, err := json.MarshalIndent(fullList, "", " ")
	if err != nil {
		fmt.Printf("mfwarch: unable to serialize firmware list: %s\n", err)
		os.Exit(1)
	}

	_, err = file.Write(fwlist)
	if err != nil {
		fmt.Printf("mfwarch: unable to save firmware: %s\n", err)
		os.Exit(1)
	}
}

func getLocalFirmwareList(filePath string) (FirmwareList, error) {
	file, err := os.Open(filePath)
	if err != nil { return FirmwareList{}, err }
	defer file.Close()

	definition, err := io.ReadAll(bufio.NewReader(file))
	if err != nil { return FirmwareList{}, err }

	var parsed FirmwareList
	err = json.Unmarshal(definition, &parsed)
	if err != nil { return FirmwareList{}, err }

	return parsed, nil
}

type None struct {}
