import { PrismaClient } from "../generated/client/index.js";
import type { Player } from "../generated/client/index.d.ts";
import { randomInt } from "node:crypto";

export { Player };

export default class DB {
  private static prisma = new PrismaClient();

  public static async getOrCreatePlayerToken(playerUuid: string, playerName: string): Promise<Player> {
    const token = TokenGenerator.generate();

    return await this.prisma.player.upsert({
      where: {
        playerUuid,
      },
      update: {
        playerName,
      },
      create: {
        playerUuid,
        playerName,
        token,
      },
    });
  }

  public static async validatePlayerToken(playerUuid: string, token: string): Promise<boolean> {
    const user = await this.prisma.player.findUnique({
      where: {
        playerUuid,
      },
    });

    if (!user) {
      return false;
    }
    return user.token === token;
  }

  public static async setPlayerColor(playerUuid: string, color: string): Promise<void> {
    await this.prisma.colorCache.upsert({
      where: {
        playerPlayerUuid: playerUuid,
      },
      update: {
        color,
      },
      create: {
        playerPlayerUuid: playerUuid,
        color,
      },
    });
  }

  public static async getPlayerColors(): Promise<Map<string, string>> {
    const players = await this.prisma.colorCache.findMany();
    const map = new Map<string, string>();
    for (const player of players) {
      map.set(player.playerPlayerUuid, player.color);
    }
    return map;
  }

  public static async setOrCreateStat(stat: string, value: number): Promise<void> {
    await this.prisma.stat.upsert({
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
    await this.prisma.stat.upsert({
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

class TokenGenerator {
  private static readonly PREFIX = "VM_";
  private static readonly CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-";
  private static readonly LENGTH = 32;
  private static readonly GENERATE_LENGTH = TokenGenerator.LENGTH - TokenGenerator.PREFIX.length;

  public static generate(): string {
    let token = TokenGenerator.PREFIX;
    for (let i = 0; i < TokenGenerator.GENERATE_LENGTH; i++) {
      token += TokenGenerator.CHARSET.charAt(randomInt(0, TokenGenerator.CHARSET.length));
    }
    return token;
  }
}
