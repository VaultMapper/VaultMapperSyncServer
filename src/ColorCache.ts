import DB from "./DB.ts";

export default class ColorCache {
  private static cache: Map<string, string> = new Map();

  public static async loadFromDB() {
    this.cache = await DB.getPlayerColors();
  }

  public static getColor(uuid: string): string {
    if (this.cache.has(uuid)) {
      return this.cache.get(uuid)!;
    } else {
      const color = this.generateColor();
      this.cache.set(uuid, color);
      DB.setPlayerColor(uuid, color);
      return color;
    }
  }

  private static generateColor(): string {
    return "#" + Math.floor(Math.random() * 16777215).toString(16); // TODO: implement the skin to color algorithm
  }
}
