package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LukaGiorgadze/gonull"
	"github.com/davecgh/go-spew/spew"
)

const begodeFirmwareListURL = "http://one-api.begode.com/agent/dog/wheel/dev/type/page?current=1&size=1000";
const begodeWheelFirmwareListURL = "http://one-api.begode.com/api/version/one/version"
const begodeFirmwareDownloadURL = "http://one-api.begode.com/imgs/%s"

type BegodeAPIResponse[T any] struct {
	Code int
	Data BegodeAPIResponseHeader[T]
}

type BegodeAPIResponseHeader[T any] struct {
	Current int
	Pages int
	Records []T
}

type BegodeFirmwareList struct {
	Code string
	ID int
	Name string
	Size int
	Total int
}

type BegodeFirmwareDetails struct {
	ApkWrap string
	AppType int
	BriefIntroduction string
	CompulsoryUpgrading int
	FirmwareTypeCode string
	Name string
	PerfectBug string
	PkgName string
	VersionCode int
	VersionName int
}

type BegodeFirmwareMisc struct {
	ListSection string
	UpgradeRequired bool
	AppType int
}

func begodeDownloadAll(destination string, existingFwHashes []string) []Firmware[BegodeFirmwareMisc] {
	res, err := http.Get(begodeFirmwareListURL)
	if err != nil {
		fmt.Printf("mfwarch: error while running: %s\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("mfwarch: couldn't get response body: %s\n", err)
		os.Exit(1)
	}
	var parsed BegodeAPIResponse[BegodeFirmwareList]
	json.Unmarshal(resBody, &parsed)

	var rawLists []struct {
		res BegodeAPIResponse[BegodeFirmwareDetails]
		name string
	}

	for _, element := range parsed.Data.Records {
		req, err := http.NewRequest(http.MethodGet, begodeWheelFirmwareListURL, nil)
		if err != nil {
			fmt.Printf("mfwarch: error making http request: %s\n", err)
			os.Exit(1)
		}
		q := req.URL.Query()
		q.Add("firmwareTypeCode", element.Code)
		req.URL.RawQuery = q.Encode()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("mfwarch: error making http request: %s\n", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("mfwarch: couldn't read response for firmware details: %s\n", err)
			os.Exit(1)
		}

		var parsed BegodeAPIResponse[BegodeFirmwareDetails]
		json.Unmarshal(resBody, &parsed)

		list := struct {
			res BegodeAPIResponse[BegodeFirmwareDetails]
			name string
		} {
			res: parsed,
			name: element.Name,
		}
		rawLists = append(rawLists, list)
	}

	var result []Firmware[BegodeFirmwareMisc]

	for _, item := range rawLists {
		for _, fw := range item.res.Data.Records {
			res, err := http.Get(fmt.Sprintf(begodeFirmwareDownloadURL, fw.ApkWrap))
			if err != nil {
				fmt.Printf("mfwarch: unable to download firmware file: %s\n", err)
				os.Exit(1)
			}
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			hash := sha256.New()
			hash.Write([]byte(resBody))
			hashString := fmt.Sprintf("%x", hash.Sum(nil))

			path := fmt.Sprintf("%s/%s.bin", destination, hashString)
			
			_, err = os.Stat(path)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				fmt.Printf("mfwarch: unable to stat file: %s\n", err)
				os.Exit(1)
			}
			if err != nil {
				file, err := os.Create(path)
				if err != nil {
					fmt.Printf("mfwarch: unable to create file: %s\n", err)
					os.Exit(1)
				}

				_, err = file.Write(resBody)
				if err != nil {
					fmt.Printf("mfwarch: unable to save firmware: %s\n", err)
					os.Exit(1)
				}
			}

			result = append(result, Firmware[BegodeFirmwareMisc]{
				Name: fw.BriefIntroduction,
				ForWheel: "", // FIXME: detect wheel model with pkgName
				VersionCode: fw.PkgName,
				VersionNumber: gonull.NewNullable(int(fw.VersionName)),
				Description: fw.PerfectBug,
				Hash: hashString,
				TimeDiscovered: time.Now(),
				AvailableUpstream: true,
				OriginalFileName: strings.Replace(fw.ApkWrap, "temp/", "", 1),
				Misc: BegodeFirmwareMisc {
					ListSection: item.name,
					UpgradeRequired: fw.CompulsoryUpgrading != 0,
					AppType: fw.AppType,
				},
			})
		}
	}

	// WARNING: test code, remove before prod

	hash := sha256.New()

	randomIdx := rand.IntN(len(result))
	randomDuplicate := result[randomIdx]
	randomDuplicate.Hash = fmt.Sprintf("%x", hash.Sum([]byte{byte(randomIdx)}))
	result = append(result, randomDuplicate)
	spew.Dump(randomDuplicate)

	return result
}

