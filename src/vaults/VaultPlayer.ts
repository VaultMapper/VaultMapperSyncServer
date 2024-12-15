import { Capsule } from "../packets/Capsule.ts";

export default class VaultPlayer {
  public readonly uuid: string;
  public readonly ws: WebSocket;
  private x: number = 0;
  private z: number = 0;
  private yaw: number = 0;

  constructor(uuid: string, ws: WebSocket) {
    this.uuid = uuid;
    this.ws = ws;
  }

  public getX(): number {
    return this.x;
  }

  public getZ(): number {
    return this.z;
  }

  public setX(x: number): void {
    this.x = x;
  }

  public setZ(z: number): void {
    this.z = z;
  }

  public getYaw(): number {
    return this.yaw;
  }

  public setYaw(yaw: number): void {
    this.yaw = yaw;
  }

  public sendPacket(packet: Capsule): void {
    this.ws.send(JSON.stringify(packet));
  }
}
