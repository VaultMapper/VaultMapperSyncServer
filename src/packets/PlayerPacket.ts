import ColorCache from "../ColorCache.ts";

export default abstract class PlayerPacket {
  public readonly uuid: string;
  public readonly color: string;

  constructor(uuid: string) {
    this.uuid = uuid;
    this.color = ColorCache.getColor(uuid);
  }
}
