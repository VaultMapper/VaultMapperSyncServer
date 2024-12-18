import { PrismaClient } from "../generated/client/index.js";
import type { Player, PrismaClient as PC } from "../generated/client/index.d.ts";

export { Player };

export default class DB {
  private static prisma: PC = new PrismaClient();

  public static async init(): Promise<void> {
    await this.prisma.$queryRaw`PRAGMA journal_mode=WAL;`;
  }

  public static async setPlayerColor(playerUuid: string, color: string): Promise<void> {
    await this.prisma.player.upsert({
      where: {
        playerUuid,
      },
      update: {
        color,
      },
      create: {
        playerUuid,
        color,
      },
    });
  }

  public static async getPlayerColors(): Promise<Map<string, string>> {
    const players = await this.prisma.player.findMany({
      select: {
        playerUuid: true,
        color: true,
      },
    });
    const map = new Map<string, string>();
    for (const player of players) {
      map.set(player.playerUuid, player.color);
    }
    return map;
  }

  public static async getStats(): Promise<{ [stat: string]: number }> {
    const stats = await this.prisma.stats.findMany({
      select: {
        stat: true,
        value: true,
      },
    });
    const map: { [stat: string]: number } = {};
    for (const stat of stats) {
      map[stat.stat] = stat.value;
    }

    return map;
  }

  public static async setOrCreateStat(stat: string, value: number): Promise<void> {
    await this.prisma.stats.upsert({
      where: {
        stat,
      },
      update: {
        value,
      },
      create: {
        stat,
        value,
      },
    });
  }

  public static async incOrCreateStat(stat: string, value: number): Promise<void> {
    await this.prisma.stats.upsert({
      where: {
        stat,
      },
      update: {
        value: {
          increment: value,
        },
      },
      create: {
        stat,
        value,
      },
    });
  }
}
