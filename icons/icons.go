package icons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
)

var jsonData = `{
    "roomIcons": {
        "the_vault:gui/map/crystal_caves": [
            "the_vault:vault/rooms/challenge/crystal_caves1",
            "the_vault:vault/rooms/challenge/crystal_caves2",
            "the_vault:vault/rooms/challenge/crystal_caves3",
            "the_vault:vault/rooms/challenge/crystal_caves4",
            "the_vault:vault/rooms/challenge/crystal_caves5"
        ],
        "the_vault:gui/map/dragon": [
            "the_vault:vault/rooms/challenge/dragon1"
        ],
        "the_vault:gui/map/factory": [
            "the_vault:vault/rooms/challenge/factory1"
        ],
        "the_vault:gui/map/raid": [
            "the_vault:vault/rooms/challenge/raid1"
        ],
        "the_vault:gui/map/lab": [
            "the_vault:vault/rooms/challenge/laboratory1",
            "the_vault:vault/rooms/challenge/laboratory2"
        ],
        "the_vault:gui/map/xmark": [
            "the_vault:vault/rooms/challenge/x-mark1"
        ],
        "the_vault:gui/map/village": [
            "the_vault:vault/rooms/challenge/village1",
            "the_vault:vault/rooms/challenge/village2",
            "the_vault:vault/rooms/challenge/village3",
            "the_vault:vault/rooms/challenge/village4"
        ],
        "the_vault:gui/map/chromatic_caves": [
            "the_vault:vault/rooms/raw/chromatic_cave1"
        ],
        "the_vault:gui/map/emerald_caves": [
            "the_vault:vault/rooms/raw/emerald_cave1"
        ],
        "the_vault:gui/map/diamond_caves": [
            "the_vault:vault/rooms/raw/diamond_cave1"
        ],
        "the_vault:gui/map/farm": [
            "the_vault:vault/decor/raw/farm/farm1",
            "the_vault:vault/decor/raw/farm/farm2",
            "the_vault:vault/decor/raw/farm/farm3",
            "the_vault:vault/decor/raw/farm/farm4",
            "the_vault:vault/decor/raw/farm/farm5",
            "the_vault:vault/decor/raw/farm/farm6"
        ],
        "the_vault:gui/map/cove": [
            "the_vault:vault/rooms/omega/cove1"
        ],
        "the_vault:gui/map/digsite": [
            "the_vault:vault/rooms/omega/omega_digsite1",
            "the_vault:vault/rooms/omega/omega_digsite2"
        ],
        "the_vault:gui/map/vendor": [
            "the_vault:vault/rooms/omega/vendor1",
            "the_vault:vault/rooms/omega/vendor2",
            "the_vault:vault/rooms/omega/vendor3"
        ],
        "the_vault:gui/map/paint": [
            "the_vault:vault/rooms/omega/painting1"
        ],
        "the_vault:gui/map/mine": [
            "the_vault:vault/rooms/omega/mine1",
            "the_vault:vault/rooms/omega/mine2"
        ],
        "the_vault:gui/map/blacksmith": [
            "the_vault:vault/rooms/omega/blacksmith1"
        ],
        "the_vault:gui/map/library": [
            "the_vault:vault/rooms/omega/library1"
        ],
        "the_vault:gui/map/mushroom": [
            "the_vault:vault/rooms/omega/mush_room1",
            "the_vault:vault/rooms/omega/mush_room2"
        ]
    }
}`

var data struct {
	RoomIcons map[string][]string `json:"roomIcons"`
}

func Init() {
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON: ", err)
		return
	}
}

func GetIcon(roomName *string) image.Image {
	if (roomName == nil) || (*roomName == "") {
		return nil
	}

	if strings.Contains(*roomName, "boss") {
		return ReadIcon("./icons/boss.png")
	}

	if strings.Contains(*roomName, "raid") {
		return ReadIcon("./icons/raid.png")
	}

	for iconPathOriginal := range data.RoomIcons {
		for _, roomNameOriginal := range data.RoomIcons[iconPathOriginal] {
			if roomNameOriginal == *roomName {
				iconFileSplit := strings.Split(iconPathOriginal, "/")
				iconFile := iconFileSplit[len(iconFileSplit)-1]
				return ReadIcon("./icons/" + iconFile + ".png")
			}
		}
	}

	// if not found, try path from room name (can be the case with unexplored inscription rooms)
	rnSplit := strings.Split(*roomName, "/")
	iconName := rnSplit[len(rnSplit)-1]

	return ReadIcon("./icons/" + iconName + ".png")
}

func ReadIcon(relPath string) image.Image {
	img := readIcon(relPath)
	if img == nil {
		fmt.Println("Error reading icon: ", relPath)
	}
	return img
}

func readIcon(relPath string) image.Image {
	path, err := filepath.Abs(relPath)
	if err != nil {
		//fmt.Println("Error getting relative path: ", err)
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		//fmt.Println("Error reading icon file: ", err)
		return nil
	}

	info, err := file.Stat()
	if err != nil {
		//fmt.Println("Error getting file info: ", err)
		return nil
	}

	size := info.Size()

	imgData := make([]byte, size)
	_, err = file.Read(imgData)
	if err != nil {
		//fmt.Println("Error reading file: ", err)
		return nil
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			//fmt.Println("Error closing file: ", err)
			return
		}
	}(file)
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		//fmt.Println("Error decoding image: ", err)
		return nil
	}
	return img
}
