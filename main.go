package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
	"github.com/davecgh/go-spew/spew"
	cp "github.com/otiai10/copy"
)

func main() {
    fmt.Println("Hello from mfwarch!")

	_, err := os.Stat("firmwares")
	if os.IsNotExist(err) {
		os.Mkdir("firmwares", 0770)
		list := begodeDownloadAll("firmwares")
		file, err := os.Create("firmwares/firmwares.json")
		if err != nil {
			fmt.Printf("mfwarch: unable to create file: %s\n", err)
			os.Exit(1)
		}

		fullList := FirmwareList {
			LastUpdateTime: time.Now(),	
			ChangesSinceLastUpdate: 0,
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

	err = cp.Copy("firmwares", "firmwares2")
	if err != nil {
		fmt.Printf("mfwarch: unable to copy firmwares: %s\n", err)
	}

	fmt.Println("mfwarch: downloading Begode firmwares...")
	list := begodeDownloadAll("firmwares2")
	for index, item := range list {
		item.TimeDiscovered = time.Time{}
		item.TimeCreated = time.Time{}
		list[index] = item
	}

	fmt.Println("mfwarch: finding differences...")
	diff := difference(list, fwList.Begode)
	spew.Dump(diff)
	fmt.Println(len(diff))
	if len(diff) == 0 {
		fmt.Println("mfwarch: no differences found")
		list = fwList.Begode
	} else { // process the changes
		fmt.Println("mfwarch: found differences")
		for _, item := range diff {
			// check if new list has a new firmware
			if slices.ContainsFunc(list, func(item2 Firmware[BegodeFirmwareMisc]) bool {
				return item2.Hash == item.Hash
			}) {
				// search for silent replacements
				idx := slices.IndexFunc(fwList.Begode,
					func(fw Firmware[BegodeFirmwareMisc]) bool {
						return fw.Name == item.Name &&
						fw.VersionCode == item.VersionCode &&
						fw.Description == item.Description &&
						fw.Misc.ListSection == item.Misc.ListSection &&
						fw.Hash != item.Hash
				})

				if idx != -1 {
					item.SilentReplacementOf = fwList.Begode[idx].Hash
				}
				// search for reappearances
				idx = slices.IndexFunc(fwList.Begode,
					func(fw Firmware[BegodeFirmwareMisc]) bool {
						return fw.Name == item.Name &&
						fw.VersionCode == item.VersionCode &&
						fw.Description == item.Description &&
						fw.Misc.ListSection == item.Misc.ListSection &&
						fw.Hash == item.Hash
				})
				if idx != -1 {
					fwList.Begode[idx].AvailableUpstream = true
					continue
				}

				item.TimeDiscovered = time.Now()
				fwList.Begode = append(fwList.Begode, item)
			}

			// if any firmwares were removed
			idx := slices.IndexFunc(fwList.Begode, func(item2 Firmware[BegodeFirmwareMisc]) bool {
				return item2.Hash == item.Hash
			})
			if idx != -1 {
				fwList.Begode[idx].AvailableUpstream = false
			}
		}
	}

	os.Remove("firmwares2/firmwares.json")
	file, err := os.Create("firmwares2/firmwares.json")
	if err != nil {
		fmt.Printf("mfwarch: unable to create file: %s\n", err)
		os.Exit(1)
	}

	fullList := FirmwareList {
		LastUpdateTime: time.Now(),	
		ChangesSinceLastUpdate: len(diff),
		Begode: fwList.Begode,
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

	os.RemoveAll("firmwares")
	os.Rename("firmwares2", "firmwares")
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
