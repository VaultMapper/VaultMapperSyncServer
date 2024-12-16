import CellType from "./CellType.ts";
import RoomName from "./RoomName.ts";
import RoomType from "./RoomType.ts";

export default class VaultCell {
  x: number;
  z: number;
  inscripted: boolean;
  cellType: CellType;
  roomType: RoomType;
  roomName: RoomName;
  marked: boolean;

  constructor(x: number, z: number, inscripted: boolean, cellType: CellType, roomType: RoomType, roomName: RoomName, marked: boolean) {
    this.x = x;
    this.z = z;
    this.inscripted = inscripted;
    this.cellType = cellType;
    this.roomType = roomType;
    this.roomName = roomName;
    this.marked = marked;
  }
}
