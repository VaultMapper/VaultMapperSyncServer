import ColorCache from "../ColorCache.ts";
import DB from "../DB.ts";
import { Capsule, CellPacketCapsule, JoinPacketCapsule, LeavePacketCapsule } from "../packets/Capsule.ts";
import CellPacket from "../packets/CellPacket.ts";
import JoinPacket from "../packets/JoinPacket.ts";
import LeavePacket from "../packets/LeavePacket.ts";
import CellType from "./CellType.ts";
import RoomType from "./RoomType.ts";
import Vault from "./Vault.ts";
import VaultPlayer from "./VaultPlayer.ts";

export default class VaultManager {
  private static vaults: Map<string, Vault> = new Map();

  public static getOrCreateVault(uuid: string): Vault {
    let vault = this.vaults.get(uuid);
    if (!vault) {
      vault = new Vault(uuid);
      this.vaults.set(uuid, vault);

      DB.incOrCreateStat("global:vaults_total", 1);
    }
    return vault;
  }

  public static getVaults(): Vault[] {
    return Array.from(this.vaults.values());
  }

  public static deleteVault(uuid: string): void {
    const vault = this.vaults.get(uuid);
    if (!vault) return;
    DB.incOrCreateStat("global:cells_total", vault.getCells().length || 0);
    DB.incOrCreateStat(
      "global:rooms_total",
      this.vaults
        .get(uuid)!
        .getCells()
        .filter((cell) => cell.cellType === CellType.ROOM).length || 0,
    );
    DB.incOrCreateStat(
      "global:omega_rooms_total",
      this.vaults
        .get(uuid)!
        .getCells()
        .filter((cell) => cell.roomType === RoomType.OMEGA).length || 0,
    );
    DB.incOrCreateStat(
      "global:challenge_rooms_total",
      this.vaults
        .get(uuid)!
        .getCells()
        .filter((cell) => cell.roomType === RoomType.CHALLENGE).length || 0,
    );

    this.vaults.delete(uuid);
  }

  public static connect(ws: WebSocket, uuid: string, vaultID: string): void {
    const player = new VaultPlayer(uuid, ColorCache.getColor(uuid), ws);
    const vault = this.getOrCreateVault(vaultID);

    ws.addEventListener("open", () => {
      vault.addOrUpdatePlayer(player);

      for (const cell of vault.getCells()) {
        player.sendPacket(new CellPacketCapsule(CellPacket.fromVaultCell(cell)));
      }

      for (const otherPlayer of vault.getPlayers()) {
        if (otherPlayer.uuid === player.uuid) continue;
        player.sendPacket(new JoinPacketCapsule(new JoinPacket(otherPlayer.uuid))); //idk if we even need this or if we need something else, maybe a 0 move packet?
      }

      vault.broadcast(new JoinPacketCapsule(new JoinPacket(uuid)), player.uuid);
    });

    ws.addEventListener("close", () => {
      console.log(`Player ${player.uuid} disconnected from ${vaultID}`);
      vault.removePlayer(player);
      vault.broadcast(new LeavePacketCapsule(new LeavePacket(uuid)), player.uuid);
    });

    ws.addEventListener("error", (_error) => {
      // console.error(error);
      console.log(`Player ${player.uuid} disconnected from ${vaultID} due to error`);
      vault.removePlayer(player);
      vault.broadcast(new LeavePacketCapsule(new LeavePacket(uuid)), player.uuid);
    });

    ws.addEventListener("message", (event) => {
      if (typeof event.data !== "string") return;
      if (!this.isValidJsonString(event.data)) return;
      vault.handlePacket(JSON.parse(event.data) as Capsule, player);
    });
  }

  private static isValidJsonString(str: string): boolean {
    try {
      JSON.parse(str);
    } catch (_e) {
      return false;
    }
    return true;
  }
}
