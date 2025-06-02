package main

import (
	"time"
	"github.com/LukaGiorgadze/gonull"
)

type FirmwareType int

const (
	FirmwareTypeMotherboard FirmwareType = iota
	FirmwareTypeBMS
)

func (fwt FirmwareType)String() string {
	switch fwt {
		case FirmwareTypeMotherboard: return "motherboard"
		case FirmwareTypeBMS: return "bms"
	}
	return "unknown"
}

type Firmware[T any] struct {
	Name string
	ForWheel string
	Type FirmwareType
	VersionCode string
	VersionNumber gonull.Nullable[int]
	Description string
	Hash string // SHA256 hash of the firmware file
	TimeDiscovered time.Time
	TimeCreated time.Time // extracted from the description for Begode
	AvailableUpstream bool
	OriginalFileName string
	// hash of the firmware that this one silently replaced.
	// If none, just an empty string.
	// silently replaced - firmware file changed without changing name, date, desc, etc
	SilentReplacementOf string
	Misc T // for any custom data from different firmware sources
}

type FirmwareList struct {
	LastUpdateTime time.Time
	ChangesSinceLastUpdate int
	Begode []Firmware[BegodeFirmwareMisc]
}

func difference[T comparable](slice1, slice2 []Firmware[T]) []Firmware[T] {
    var diff []Firmware[T]

	for index, item := range slice1 {
		item.TimeCreated = time.Time{}
		item.TimeDiscovered = time.Time{}
		slice1[index] = item
	}

	for index, item := range slice2 {
		item.TimeCreated = time.Time{}
		item.TimeDiscovered = time.Time{}
		slice2[index] = item
	}

    // Loop two times, first to find slice1 strings not in slice2,
    // second loop to find slice2 strings not in slice1
    for i := 0; i < 2; i++ {
        for _, s1 := range slice1 {
            found := false
            for _, s2 := range slice2 {
                if s1 == s2 {
                    found = true
                    break
                }
            }
            // String not found. We add it to return slice
            if !found {
                diff = append(diff, s1)
            }
        }
        // Swap the slices, only if it was the first loop
        if i == 0 {
            slice1, slice2 = slice2, slice1
        }
    }

    return diff
}
