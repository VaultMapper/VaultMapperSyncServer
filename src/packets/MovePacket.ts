import PlayerPacket from "./PlayerPacket.ts";

export default class MovePacket extends PlayerPacket {
  public readonly x: number;
  public readonly z: number;
  public readonly yaw: number;

  constructor(uuid: string, x: number, z: number, yaw: number) {
    super(uuid);
    this.x = x;
    this.z = z;
    this.yaw = yaw;
  }
}
