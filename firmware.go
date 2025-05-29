package main

import (
	"time"
	"github.com/LukaGiorgadze/gonull"
)

type Firmware[T any] struct {
	Name string
	ForWheel string
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

func difference[T comparable](a, b []Firmware[T]) []Firmware[T] {
    // reorder the input,
    // so that we can check the longer slice over the shorter one
    longer, shorter := a, b
    if len(b) > len(a) {
        longer, shorter = b, a
    }

	for index, item := range longer {
		newItem := item // idk if Go for loops have their variables as copies or refs, so ima play it safe
		newItem.TimeCreated = time.Time{}
		newItem.TimeDiscovered = time.Time{}
		longer[index] = newItem
	}

	for index, item := range shorter {
		newItem := item // idk if Go for loops have their variables as copies or refs, so ima play it safe
		newItem.TimeCreated = time.Time{}
		newItem.TimeDiscovered = time.Time{}
		shorter[index] = newItem
	}

    mb := make(map[Firmware[T]]struct{}, len(shorter))
    for _, x := range shorter {
        mb[x] = struct{}{}
    }
    var diff []Firmware[T]
    for _, x := range longer {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
}
