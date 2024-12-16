import { createCanvas, loadImage } from "https://deno.land/x/canvas@v1.4.2/mod.ts";
import DB from "./DB.ts";

export default class ColorCache {
  private static cache: Map<string, string> = new Map();

  public static async loadFromDB() {
    this.cache = await DB.getPlayerColors();
  }

  public static randomColor(): string {
    const r = Math.floor(Math.random() * 256)
      .toString(16)
      .padStart(2, "0");
    const g = Math.floor(Math.random() * 256)
      .toString(16)
      .padStart(2, "0");
    const b = Math.floor(Math.random() * 256)
      .toString(16)
      .padStart(2, "0");
    return `#${r}${g}${b}`;
  }

  public static getColor(uuid: string): string {
    if (this.cache.has(uuid)) {
      return this.cache.get(uuid)!;
    } else {
      // const color = await ColorCache.getAverageSkinColor(uuid);
      const color = this.randomColor();
      this.cache.set(uuid, color);
      DB.setPlayerColor(uuid, color);
      return color;
    }
  }

  public static async getSkin(uuid: string) {
    const skinUrl = JSON.parse(atob((await (await fetch(`https://sessionserver.mojang.com/session/minecraft/profile/${uuid}`)).json()).properties[0].value)).textures.SKIN.url;
    return skinUrl;
  }

  public static async getAverageSkinColor(uuid: string) {
    const skin = await ColorCache.getSkin(uuid);
    const canvas = createCanvas(64, 64);
    const ctx = canvas.getContext("2d");
    const image = await loadImage(skin);
    ctx.drawImage(image, 0, 0);
    const red = [];
    const green = [];
    const blue = [];

    for (let x = 0; x < 64; x++) {
      for (let y = 0; y < 64; y++) {
        const {
          data: [r, g, b, a],
        } = ctx.getImageData(x, y, 1, 1);
        if (a == 0) continue;
        red.push(r);
        green.push(g);
        blue.push(b);
      }
    }

    const avgRed = clamp(Math.floor(red.reduce((total, current) => (total += current)) / red.length), 0, 255);
    const avgGreen = clamp(Math.floor(green.reduce((total, current) => (total += current)) / green.length), 0, 255);
    const avgBlue = clamp(Math.floor(blue.reduce((total, current) => (total += current)) / blue.length), 0, 255);

    return `#${avgRed.toString(16)}${avgGreen.toString(16)}${avgBlue.toString(16)}`;
  }
}

function clamp(num: number, min: number, max: number) {
  return Math.min(Math.max(num, min), max);
}
