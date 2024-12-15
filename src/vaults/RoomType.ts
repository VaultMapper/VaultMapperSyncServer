enum RoomType {
  START,
  BASIC,
  ORE,
  CHALLENGE,
  OMEGA,
}

function parseRoomType(roomType: string): RoomType {
  switch (roomType) {
    case "0":
      return RoomType.START;
    case "1":
      return RoomType.BASIC;
    case "2":
      return RoomType.ORE;
    case "3":
      return RoomType.CHALLENGE;
    case "4":
      return RoomType.OMEGA;
    default:
      throw new Error(`Invalid room type: ${roomType}`);
  }
}

export default RoomType;
export { parseRoomType };
