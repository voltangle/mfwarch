package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const begodeFirmwareListURL = "http://one-api.begode.com/agent/dog/wheel/dev/type/page?current=1&size=1000";
const begodeWheelFirmwareListURL = "http://one-api.begode.com/api/version/one/version"
const begodeFirmwareDownloadURL = "http://one-api.begode.com/imgs/temp/%s"

func main() {
    fmt.Println("Hello from mfwarch!")
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
	fmt.Println("Downloading Begode firmwares...")

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

		fmt.Printf("Downloading from %s...\n", element.Name)
		for _, fw := range parsed.Data.Records {
			os.MkdirAll(fmt.Sprintf("firmwares/begode/%s", element.Name), 0770)
			out, err := os.OpenFile(fmt.Sprintf(
				"firmwares/begode/%s/%s.bin",
				element.Name,
				fw.BriefIntroduction), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("mfwarch: unable to create firmware file: %s\n", err)
				os.Exit(1)
			}
			defer out.Close()

			res, err := http.Get(fmt.Sprintf(begodeFirmwareDownloadURL, fw.ApkWrap))
			if err != nil {
				fmt.Printf("mfwarch: unable to download firmware file: %s\n", err)
				os.Exit(1)
			}
			defer res.Body.Close()

			_, err = io.Copy(out, res.Body)
			if err != nil {
				fmt.Printf("mfwarch: unable to save firmware: %s\n", err)
				os.Exit(1)
			}
		}
	}
}

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
