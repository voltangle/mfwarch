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
	TimeCreated time.Time // extracted from the description
	Misc T // for any custom data from different firmware sources
}

type FirmwareList struct {
	LastUpdateTime time.Time
	ChangesSinceLastUpdate int
	Begode []Firmware[BegodeFirmwareMisc]
}
