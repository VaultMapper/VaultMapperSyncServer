package icons

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/NodiumHosting/VaultMapperSyncServer/proto"
)

var iconMap = map[proto.RoomName]image.Image{
	proto.RoomName_ROOMNAME_BLACKSMITH:       ReadIcon("./icons/blacksmith.png"),
	proto.RoomName_ROOMNAME_COVE:             ReadIcon("./icons/cove.png"),
	proto.RoomName_ROOMNAME_CRYSTAL_CAVES:    ReadIcon("./icons/crystal_caves.png"),
	proto.RoomName_ROOMNAME_CUBE:             ReadIcon("./icons/cube.png"),
	proto.RoomName_ROOMNAME_DIG_SITE:         ReadIcon("./icons/dig_site.png"),
	proto.RoomName_ROOMNAME_DRAGON:           ReadIcon("./icons/dragon.png"),
	proto.RoomName_ROOMNAME_FACTORY:          ReadIcon("./icons/factory.png"),
	proto.RoomName_ROOMNAME_LIBRARY:          ReadIcon("./icons/library.png"),
	proto.RoomName_ROOMNAME_MINE:             ReadIcon("./icons/mine.png"),
	proto.RoomName_ROOMNAME_MUSH_ROOM:        ReadIcon("./icons/mush_room.png"),
	proto.RoomName_ROOMNAME_PAINTING:         ReadIcon("./icons/painting.png"),
	proto.RoomName_ROOMNAME_VENDOR:           ReadIcon("./icons/vendor.png"),
	proto.RoomName_ROOMNAME_VILLAGE:          ReadIcon("./icons/village.png"),
	proto.RoomName_ROOMNAME_WILD_WEST:        ReadIcon("./icons/wild_west.png"),
	proto.RoomName_ROOMNAME_X_MARK:           ReadIcon("./icons/x_mark.png"),
	proto.RoomName_ROOMNAME_RAID:             ReadIcon("./icons/raid.png"),
	proto.RoomName_ROOMNAME_LAB:              ReadIcon("./icons/laboratory.png"),
	proto.RoomName_ROOMNAME_PIRATE_CAVE:      ReadIcon("./icons/pirate_cave.png"),
	proto.RoomName_ROOMNAME_GARDEN:           ReadIcon("./icons/woldian_garden.png"),
	proto.RoomName_ROOMNAME_ARCADE:           ReadIcon("./icons/arcade.png"),
	proto.RoomName_ROOMNAME_COMET:            ReadIcon("./icons/comet_observatory.png"),
	proto.RoomName_ROOMNAME_PLAYZONE:         ReadIcon("./icons/playzone.png"),
	proto.RoomName_ROOMNAME_HELLISH_DIG_SITE: ReadIcon("./icons/hellish_digsite.png"),
	proto.RoomName_ROOMNAME_BOSS:             ReadIcon("./icons/boss_room.png"),
	proto.RoomName_ROOMNAME_CHROMATIC_CAVES:  ReadIcon("./icons/chromatic_caves.png"),
	proto.RoomName_ROOMNAME_DIAMOND_CAVES:    ReadIcon("./icons/diamond_caves.png"),
	proto.RoomName_ROOMNAME_EMERALD_CAVES:    ReadIcon("./icons/emerald_caves.png"),
	proto.RoomName_ROOMNAME_END:              ReadIcon("./icons/raw_end.png"),
	proto.RoomName_ROOMNAME_MODDED_CAVES:     ReadIcon("./icons/raw_modded_caves.png"),
	proto.RoomName_ROOMNAME_NETHER:           ReadIcon("./icons/raw_nether.png"),
	proto.RoomName_ROOMNAME_QUARRY:           ReadIcon("./icons/raw_quarry.png"),
	proto.RoomName_ROOMNAME_FARM:             ReadIcon("./icons/the_farm.png"),
}

func GetIcon(roomName *proto.RoomName) image.Image {
	return iconMap[*roomName]
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
