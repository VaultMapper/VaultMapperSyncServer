enum RoomName {
  UNKNOWN = "UNKNOWN",
  BLACKSMITH = "Blacksmith",
  COVE = "Cove",
  CRYSTAL_CAVES = "Crystal Caves",
  DIG_SITE = "Dig Site",
  DRAGON = "Dragon",
  FACTORY = "Factory",
  LIBRARY = "Library",
  MINE = "Mine",
  MUSH_ROOM = "Mush Room",
  PAINTING = "Painting",
  VENDOR = "Vendor",
  VILLAGE = "Village",
  WILD_WEST = "Wild West",
  X_MARK = "X-mark",
  CUBE = "Cube",
}

function parseRoomName(roomName: string): RoomName {
  switch (roomName) {
    case "0":
      return RoomName.UNKNOWN;
    case "1":
      return RoomName.BLACKSMITH;
    case "2":
      return RoomName.COVE;
    case "3":
      return RoomName.CRYSTAL_CAVES;
    case "4":
      return RoomName.DIG_SITE;
    case "5":
      return RoomName.DRAGON;
    case "6":
      return RoomName.FACTORY;
    case "7":
      return RoomName.LIBRARY;
    case "8":
      return RoomName.MINE;
    case "9":
      return RoomName.MUSH_ROOM;
    case "10":
      return RoomName.PAINTING;
    case "11":
      return RoomName.VENDOR;
    case "12":
      return RoomName.VILLAGE;
    case "13":
      return RoomName.WILD_WEST;
    case "14":
      return RoomName.X_MARK;
    case "15":
      return RoomName.CUBE;
    default:
      throw new Error(`Unknown room name: ${roomName}`);
  }
}

export default RoomName;
export { parseRoomName };
