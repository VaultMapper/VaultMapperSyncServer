import type { Capsule } from "../packets/Capsule.ts";
import { MovePacketCapsule, CellPacketCapsule, DATA } from "../packets/Capsule.ts";
import CellPacket from "../packets/CellPacket.ts";
import MovePacket from "../packets/MovePacket.ts";
import PacketType from "../packets/PacketType.ts";
import VaultCell from "./VaultCell.ts";
import VaultManager from "./VaultManager.ts";
import VaultPlayer from "./VaultPlayer.ts";

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
    console.log(`[Vault ${this.id}] Received packet from ${player.uuid}: ${PacketType[packet.type as unknown as keyof typeof PacketType]}`); // TODO: Remove this line
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
  }

  public broadcast(packet: Capsule, excludeUuid: string): void {
    for (const player of this.players.values()) {
      if (player.uuid === excludeUuid) continue;
      player.sendPacket(packet);
    }
  }
}
