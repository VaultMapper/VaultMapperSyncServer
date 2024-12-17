import "jsr:@std/dotenv/load";
import ColorCache from "./ColorCache.ts";
import VaultManager from "./vaults/VaultManager.ts";
import DB from "./DB.ts";

await ColorCache.loadFromDB();

const HOST = Deno.env.get("HOST") || "127.0.0.1";
const PORT = parseInt(Deno.env.get("PORT") || "25284");

Deno.serve({ hostname: HOST, port: PORT }, async (req, info) => {
  if (req.headers.get("upgrade") != "websocket") {
    if (new URL(req.url).pathname === "/stats") {
      if (!Deno.env.get("STATS_AUTH")) {
        return new Response(null, { status: 501 });
      }
      if (req.headers.get("authorization") !== Deno.env.get("STATS_AUTH")) {
        return new Response(null, { status: 401 });
      }
      const currentVaults = VaultManager.getVaults().map((vault) => {
        return { id: vault.id, players: vault.getPlayers().map((player) => player.uuid) };
      });
      const stats = await DB.getStats();
      return new Response(JSON.stringify({ total: currentVaults.length, currentVaults, stats }), {
        headers: { "content-type": "application/json" },
      });
    }
    return new Response(null, { status: 501 });
  }
  const headers = req.headers;
  const url = new URL(req.url);
  const ip = info.remoteAddr.hostname;
  const nginxIP = headers.get("x-real-ip");

  const vaultID = url.pathname.split("/")[1];
  const uuid = url.pathname.split("/")[2];
  // ex. wss://vmsync.boykiss.ing/vault_0504c51e-421a-4512-ad28-2f67c865ac72/b66daa18-7360-4fd1-b60f-15a30cb0dccc

  const username = (await (await fetch(`https://sessionserver.mojang.com/session/minecraft/profile/${uuid}`)).json()).name || "<Unknown Username>";

  if (!uuid || !vaultID) {
    if (Deno.env.get("WEBHOOK_URL")) {
      await fetch(Deno.env.get("WEBHOOK_URL")!, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: `Player failed to connect (missing args) from ${username} (${uuid}) (${ip} / ${nginxIP}) to ${vaultID} - Possible Attacker`,
        }),
      });
    }

    console.log("Player %s failed to connect (missing args) from %s (%s) to %s - Possible Attacker", uuid, ip, nginxIP, vaultID);
    // setTimeout(() => socket.close(1008), 1); // deno bug workaround
    return new Response(null, { status: 400 });
  }
  if (!uuid.match(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/) || !vaultID.match(/^vault_[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/)) {
    if (Deno.env.get("WEBHOOK_URL")) {
      await fetch(Deno.env.get("WEBHOOK_URL")!, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          content: `Player failed to connect (regex) from ${username} (${uuid}) (${ip} / ${nginxIP}) to ${vaultID} - Possible Attacker`,
        }),
      });
    }

    console.log("Player %s failed to connect (regex) from %s (%s) to %s - Possible Attacker", uuid, ip, nginxIP, vaultID);
    // setTimeout(() => socket.close(1008), 1); // deno bug workaround
    return new Response(null, { status: 400 });
  }

  const { socket, response } = Deno.upgradeWebSocket(req);

  if (Deno.env.get("WEBHOOK_URL")) {
    const players = await Promise.all(
      VaultManager.getOrCreateVault(vaultID)
        .getPlayers()
        .map((player) => player.uuid)
        .map(async (puuid) => ((await (await (await fetch(`https://sessionserver.mojang.com/session/minecraft/profile/${puuid}`)).json()).name) as string) || "<Unknown Username>"),
    );
    players.push(username);

    await fetch(Deno.env.get("WEBHOOK_URL")!, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        embeds: [
          {
            title: `${username} connected`,
            description: `${vaultID}`,
            fields: [{ name: `Players (${players.length})`, value: `${players.join("\n")}` }],
            color: parseInt(ColorCache.getColor(uuid).slice(1), 16),
            footer: { text: "VMSync" },
          },
        ],
      }),
    });
  }

  console.log("Player %s connected from %s (%s) to %s", uuid, ip, nginxIP, vaultID);

  VaultManager.connect(socket, uuid, vaultID);

  return response;
});
