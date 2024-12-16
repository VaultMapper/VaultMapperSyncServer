import "jsr:@std/dotenv/load";
import "./authserver.ts";
import ColorCache from "./ColorCache.ts";
import DB from "./DB.ts";
import VaultManager from "./vaults/VaultManager.ts";

await ColorCache.loadFromDB();

const HOST = "0.0.0.0";
const PORT = 25284;

Deno.serve({ hostname: HOST, port: PORT }, async (req, info) => {
  if (req.headers.get("upgrade") != "websocket") {
    if (new URL(req.url).pathname === "/stats") {
      if (req.headers.get("authorization") !== Deno.env.get("STATS_AUTH")) {
        return new Response(null, { status: 401 });
      }
      const currentVaults = VaultManager.getVaults().map((vault) => {
        return { id: vault.id, players: vault.getPlayers().map((player) => player.uuid) };
      });
      return new Response(JSON.stringify({ total: currentVaults.length, vaults: currentVaults }), {
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
  const uuid = url.searchParams.get("uuid");
  const token = url.searchParams.get("token");
  // ex. wss://sync.vaultmapper.site/vault_0504c51e-421a-4512-ad28-2f67c865ac72?uuid=b66daa18-7360-4fd1-b60f-15a30cb0dccc&token=token

  if (!uuid || !vaultID || !token) {
    console.log("Player %s failed to connect (missing args) from %s (%s) to %s - Possible Attacker", uuid, ip, nginxIP, vaultID);
    // setTimeout(() => socket.close(1008), 1); // deno bug workaround
    return new Response(null, { status: 400 });
  }
  if (!uuid.match(/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$/) || !vaultID.match(/^vault_[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$/)) {
    console.log("Player %s failed to connect (regex) from %s (%s) to %s - Possible Attacker", uuid, ip, nginxIP, vaultID);
    // setTimeout(() => socket.close(1008), 1); // deno bug workaround
    return new Response(null, { status: 400 });
  }
  if (!(await DB.validatePlayerToken(uuid, token))) {
    console.debug((await DB.getOrCreatePlayerToken(uuid, "")).token, token);

    console.log("Player %s failed to connect (invalid token) from %s (%s) to %s - Possible Attacker", uuid, ip, nginxIP, vaultID);
    // setTimeout(() => socket.close(1008), 1); // deno bug workaround
    return new Response(null, { status: 400 });
  }

  const { socket, response } = Deno.upgradeWebSocket(req);

  console.log("Player %s connected from %s (%s) to %s", uuid, ip, nginxIP, vaultID);

  VaultManager.connect(socket, uuid, vaultID);

  return response;
});
