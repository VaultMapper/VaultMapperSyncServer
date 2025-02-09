import type { Capsule } from "../packets/Capsule.ts";
import { MovePacketCapsule, CellPacketCapsule, DATA } from "../packets/Capsule.ts";
import CellPacket from "../packets/CellPacket.ts";
import MovePacket from "../packets/MovePacket.ts";
import PacketType from "../packets/PacketType.ts";
import CellType from "./CellType.ts";
import RoomName from "./RoomName.ts";
import RoomType from "./RoomType.ts";
import VaultCell from "./VaultCell.ts";
import VaultManager from "./VaultManager.ts";
import VaultPlayer from "./VaultPlayer.ts";

const pixelArtAlphabet3x5: { [key: string]: string[] } = {
  a: ["###", "# #", "###", "# #", "# #"],
  b: ["## ", "# #", "## ", "# #", "## "],
  c: ["###", "#  ", "#  ", "#  ", "###"],
  d: ["## ", "# #", "# #", "# #", "## "],
  e: ["###", "#  ", "###", "#  ", "###"],
  f: ["###", "#  ", "###", "#  ", "#  "],
  g: ["###", "#  ", "# #", "# #", "###"],
  h: ["# #", "# #", "###", "# #", "# #"],
  i: ["###", " # ", " # ", " # ", "###"],
  j: ["###", "  #", "  #", "# #", "## "],
  k: ["# #", "# #", "## ", "# #", "# #"],
  l: ["#  ", "#  ", "#  ", "#  ", "###"],
  m: ["# #", "###", "# #", "# #", "# #"],
  n: ["#  ", "## ", "# #", "# #", "# #"],
  o: ["###", "# #", "# #", "# #", "###"],
  p: ["###", "# #", "###", "#  ", "#  "],
  q: ["###", "# #", "# #", "###", "  #"],
  r: ["## ", "# #", "## ", "# #", "# #"],
  s: ["###", "#  ", "###", "  #", "###"],
  t: ["###", " # ", " # ", " # ", " # "],
  u: ["# #", "# #", "# #", "# #", "###"],
  v: ["# #", "# #", "# #", "# #", " # "],
  w: ["# #", "# #", "# #", "###", "# #"],
  x: ["# #", "# #", " # ", "# #", "# #"],
  y: ["# #", "# #", " # ", " # ", " # "],
  z: ["###", "  #", " # ", "#  ", "###"],
  " ": ["   ", "   ", "   ", "   ", "   "],
};
export default class Vault {
  public readonly id: string;
  private players: Map<string, VaultPlayer> = new Map();
  private cells: VaultCell[] = [];

  constructor(id: string) {
    this.id = id;
  }

  public addOrUpdatePlayer(player: VaultPlayer): void {
    this.players.set(player.uuid, player);
  }

  public removePlayer(player: VaultPlayer): void {
    this.players.delete(player.uuid);

    if (this.players.size === 0) {
      VaultManager.deleteVault(this.id); // TODO: add a timer to only delete after it's been empty for a while
    }
  }

  public getPlayers(): VaultPlayer[] {
    return Array.from(this.players.values());
  }

  public getPlayer(uuid: string): VaultPlayer | undefined {
    return this.players.get(uuid);
  }

  public addOrUpdateCell(cell: VaultCell): void {
    const index = this.cells.findIndex((c) => c.x === cell.x && c.z === cell.z);
    if (index === -1) {
      this.cells.push(cell);
    } else {
      this.cells[index] = cell;
    }
  }

  public getCells(): VaultCell[] {
    return this.cells;
  }

  public getCell(x: number, z: number): VaultCell | undefined {
    return this.cells.find((c) => c.x === x && c.z === z);
  }

  public handlePacket(packet: Capsule, player: VaultPlayer) {
    try {
      switch (packet.type) {
        case PacketType.CELL: {
          const capsule: CellPacketCapsule = packet as CellPacketCapsule;
          this.addOrUpdateCell(CellPacket.toVaultCell(DATA(capsule)));
          break;
        }
        case PacketType.MOVE: {
          const capsule: MovePacketCapsule = packet as MovePacketCapsule;
          player.setX(DATA<typeof capsule.t>(capsule).x);
          player.setZ(DATA<typeof capsule.t>(capsule).z);
          player.setYaw(DATA<typeof capsule.t>(capsule).yaw);
          this.addOrUpdatePlayer(player);

          packet = new MovePacketCapsule(new MovePacket(player.uuid, DATA<typeof capsule.t>(capsule).x, DATA<typeof capsule.t>(capsule).z, DATA<typeof capsule.t>(capsule).yaw)); //fix color
          break;
        }
      }

      this.broadcast(packet, player.uuid);
    } catch (e) {
      if (Deno.env.get("WEBHOOK_URL")) {
        fetch(Deno.env.get("WEBHOOK_URL")!, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            content: `Failed to handle packet from ${player.uuid} in ${this.id} - ${e}`,
          }),
        });
      }
      console.error(e);
    }
  }

  public broadcast(packet: Capsule, excludeUuid: string): void {
    for (const player of this.players.values()) {
      if (player.uuid === excludeUuid) continue;
      player.sendPacket(packet);
    }
  }

  public drawPixelArtMessage(message: string, x: number, z: number, inscripted: boolean, cellType: CellType, roomType: RoomType, roomName: RoomName, marked: boolean): void {
    const lineSpacing = 3;
    const lines = message.split("\n");
    let longestLine = 0;
    for (const line of lines) {
      if (line.length > longestLine) {
        longestLine = line.length;
      }
    }
    z -= Math.round((lines.length * 6) / 2);
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      const lineLength = line.length * 4;
      const lineX = x - Math.round(lineLength / 2);
      for (let j = 0; j < line.length; j++) {
        const char = line[j].toLowerCase();
        if (pixelArtAlphabet3x5[char]) {
          const pixels = pixelArtAlphabet3x5[char];
          for (let k = 0; k < pixels.length; k++) {
            for (let l = 0; l < pixels[k].length; l++) {
              if (pixels[k][l] === "#") {
                this.addOrUpdateCell(new VaultCell(lineX + j * 4 + l, z + i * (6 + lineSpacing) + k, inscripted, cellType, roomType, roomName, marked));
              }
            }
          }
        }
      }
    }
  }
}
