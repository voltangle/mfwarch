package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/LukaGiorgadze/gonull"
)

const begodeFirmwareListURL = "http://one-api.begode.com/agent/dog/wheel/dev/type/page?current=1&size=1000";
const begodeWheelFirmwareListURL = "http://one-api.begode.com/api/version/one/version"
const begodeFirmwareDownloadURL = "http://one-api.begode.com/imgs/%s"

type BegodeWheelModel int

const (
	BegodeExtreme BegodeWheelModel = iota
	BegodeRace
	BegodePanther
	BegodeT4
	BegodeT4V2
	BegodeT4Pro
	BegodeT4Max
	BegodeMaster1
	BegodeMaster2_2
	BegodeMaster2_3
	BegodeMaster2_4
	BegodeMaster4
	BegodeMasterX
	BegodeMasterPro
	BegodeMasterPro3
	BegodeEX20S
	BegodeEX30
	BegodeA1
	BegodeA2
	BegodeA5
	BegodeMten3
	BegodeMten4
	BegodeMten5
	BegodeMtenMini
	BegodeETMax
	BegodeBlitz
	BegodeBlitzPro_0
	BegodeBlitzPro_1
	BegodeXWay
	BegodeFalcon
	BegodeC8
	BegodeFuture
	BegodeMSuperPro_0
	BegodeMSuperPro_1
	BegodeMonsterPro
	BegodeMSuperRSC30
	BegodeMSuperRSC38
	BegodeMCM5
	BegodeNikola
	ExtremeBullCommanderMax
	ExtremeBullRocket
	ExtremeBullGriffin
	ExtremeBullGTProPlus
	ExtremeBullGTPro
	ExtremeBullGT
	ExtremeBullCommanderPro
	ExtremeBullCommanderPro50S
	ExtremeBullXM
	ExtremeBullCommanderMini
	ExtremeWheelES
	ExtremeWheelK4
	ExtremeWheelK6_0
	ExtremeWheelK6_1
)

func (model BegodeWheelModel)String() string {
	switch model {
		case BegodeExtreme: return "Extreme"
		case BegodeRace: return "Race"
		case BegodePanther: return "Panther"
		case BegodeT4: return "T4"
		case BegodeT4Pro: return "T4 Pro"
		case BegodeT4Max: return "T4 Max"
		case BegodeMaster1: return "Master1"
		case BegodeMaster2_2: return "Master2.2"
		case BegodeMaster2_3: return "Master2.3"
		case BegodeMaster2_4: return "Master2.4"
		case BegodeMaster4: return "Master4"
		case BegodeMasterX: return "Master X"
		case BegodeMasterPro: return "Master Pro"
		case BegodeMasterPro3: return "Master Pro 3"
		case BegodeEX20S: return "EX20S"
		case BegodeEX30: return "EX30"
		case BegodeA1: return "A1"
		case BegodeA2: return "A2"
		case BegodeA5: return "A5"
		case BegodeMten3: return "Mten3"
		case BegodeMten4: return "Mten4"
		case BegodeMten5: return "Mten5"
		case BegodeMtenMini: return "Mten Mini"
		case BegodeETMax: return "ET Max"
		case BegodeBlitz: return "Blitz"
		case BegodeBlitzPro_0: return "Blitz Pro"
		case BegodeBlitzPro_1: return "Blitz Pro"
		case BegodeXWay: return "X-Way"
		case BegodeFalcon: return "Falcon"
		case BegodeC8: return "C8"
		case BegodeFuture: return "Future"
		case BegodeMSuperPro_0: return "MSuper Pro"
		case BegodeMSuperPro_1: return "MSuper Pro"
		case BegodeMonsterPro: return "Monster Pro"
		case BegodeMSuperRSC30: return "MSuper RS C30"
		case BegodeMSuperRSC38: return "MSuper RS C38"
		case BegodeMCM5: return "MCM5"
		case BegodeNikola: return "Nikola"
		case ExtremeBullCommanderMax: return "CommanderMax"
		case ExtremeBullRocket: return "Rocket"
		case ExtremeBullGriffin: return "Griffin"
		case ExtremeBullGTProPlus: return "GT Pro+"
		case ExtremeBullGTPro: return "GT Pro"
		case ExtremeBullGT: return "GT"
		case ExtremeBullCommanderPro: return "Commander Pro"
		case ExtremeBullCommanderPro50S: return "Commander Pro 50S"
		case ExtremeBullXM: return "XM"
		case ExtremeBullCommanderMini: return "Commander Mini"
		case ExtremeWheelES: return "ES"
		case ExtremeWheelK4: return "K4"
		case ExtremeWheelK6_0: return "K6"
		case ExtremeWheelK6_1: return "K6"
	}

	return "Unknown"
}

func makeBegodeModel(name string, code string) BegodeWheelModel {
	brandCode := code[0:2]
	hardwareCode, err := strconv.ParseInt(code[2:7], 10, 32)
	name = strings.ToLower(name)
	if err != nil {
		fmt.Printf("mfwarch: failed to parse name %s and code %s: %s\n", name, code, err)
		return -1
	}

	// TODO: also identify BMS model by name or hardware code

	switch brandCode {
	case "GW":
		switch hardwareCode {
			case 18250: return BegodeExtreme
			case 20262: return BegodePanther
			case 20270: return BegodeRace
			case 16121: return BegodeT4
			case 16122: return BegodeT4V2
			case 16250: return BegodeT4Pro
			case 16251: return BegodeT4Max
			case 20140: return BegodeMaster1
			case 20149: return BegodeMaster2_2
			case 20148: return BegodeMaster2_3
			case 20150: return BegodeMaster2_4
			case 20151: return BegodeMaster4
			case 20041: return BegodeMasterX
			case 23040: return BegodeMasterPro
			case 23250: return BegodeMasterPro3
			case 20250: return BegodeEX30
			case 20130: return BegodeEX20S
			case 14210: return BegodeA1
			case 15110: return BegodeA2
			case 15111: return BegodeA5
			case 10010: return BegodeMten3
			case 10110: return BegodeMten4
			case 12110: return BegodeMten5
			case 11210: return BegodeMtenMini
			case 20260: return BegodeETMax
			case 20361: return BegodeBlitzPro_1
			case 20351: return BegodeBlitz
			case 20360: return BegodeBlitzPro_0
			case 20263: return BegodeXWay
			case 16210: return BegodeFalcon
			case 18210: return BegodeC8
			case 15010: return BegodeFuture
			case 19020:
				if strings.Contains(name, "rs") { return BegodeMSuperRSC30 }
				return BegodeMSuperPro_0
			case 19120: 
				if strings.Contains(name, "rs") { return BegodeMSuperRSC38 }
				return BegodeMSuperPro_1
			case 24020: return BegodeMonsterPro
			case 14010: return BegodeMCM5
			case 17020: return BegodeNikola
		}
	case "JN":
		switch hardwareCode {
			case 20080: 
			if strings.Contains(name, "griffin") {
				return ExtremeBullGriffin
			}
			return ExtremeBullCommanderMax
			case 15250: return ExtremeBullRocket
			case 20261: return ExtremeBullGTProPlus
			case 20260: return ExtremeBullGTPro
			case 20251: return ExtremeBullGT
			case 20252: return ExtremeBullCommanderPro50S
			case 20122: return ExtremeBullCommanderPro
			case 18150: return ExtremeBullCommanderMini
		}
	case "JL":
		switch hardwareCode {
			case 24021: return ExtremeWheelES	
			case 13051: return ExtremeWheelK6_1
			case 13050: return ExtremeWheelK6_0
		}
	}

	return -1
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

type BegodeFirmwareMisc struct {
	ListSection string
	UpgradeRequired bool
	AppType int
}

func begodeDownloadAll(destination string) []Firmware[BegodeFirmwareMisc] {
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
				ForWheel: makeBegodeModel(fw.BriefIntroduction, fw.PkgName).String(),
				Type: FirmwareTypeMotherboard, // FIXME: detect this shit
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

	return result
}
