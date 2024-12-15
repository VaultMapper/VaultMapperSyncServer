import { Buffer } from "node:buffer";
import * as mc from "minecraft-protocol";
import data from "minecraft-data";
import DB from "./DB.ts";

interface Packet {
  channel: string;
  data: any;
}

const HOST = "0.0.0.0";
const PORT = 25565;
const VERSION = "1.18.2";
const PACKET_CHANNEL = "vaultmapper:auth";

const server = mc.createServer({
  "online-mode": false, // TODO: change
  host: HOST,
  port: PORT,
  version: VERSION,
  motd: "Vault Mapper Auth",
  maxPlayers: 1,
});

const mcData = data(VERSION);

server.on("listening", () => {
  console.log(`Minecraft Auth server now listening on ${HOST}:${PORT}`);
});

server.on("playerJoin", async (client) => {
  const loginPacket = mcData.loginPacket;

  client.write("login", {
    ...loginPacket,
    entityId: client.id,
    hashedSeed: [0, 0],
    maxPlayers: server.maxPlayers,
    viewDistance: 10,
    reducedDebugInfo: false,
    enableRespawnScreen: true,
    isDebug: false,
    isFlat: false,
  });

  const kicktimer = setTimeout(() => {
    client.end("§cAuthentication failed, please try again later.");
    console.log(`Authentication failed for ${client.username} (timed out)`);
  }, 15000);

  const chatMessageTimer = setInterval(() => {
    client.write("chat", {
      message: JSON.stringify({
        text: "Currently attempting authentication. You'll be disconnected once the process is complete.",
      }),
      position: 0,
      sender: "me",
    });
  }, 5000);

  const player = await DB.getOrCreatePlayerToken(client.uuid, client.username);

  const tokenPacket: AuthPacket = {
    type: AuthPacketType.Token,
    data: player.token,
  };
  const packetBuffer = createCustomPayloadPacket(0x01, tokenPacket);
  client.write("custom_payload", {
    channel: PACKET_CHANNEL,
    data: packetBuffer,
  });

  client.on("packet", (data: Packet) => {
    if (data.channel !== PACKET_CHANNEL) {
      return;
    }

    console.log(data.data, typeof data.data);

    const json = Buffer.from(data.data).toString().substring(2);

    let authPacket: AuthPacket;
    try {
      authPacket = JSON.parse(json);
    } catch (error) {
      console.log(error);
      return;
    }

    switch (authPacket.type) {
      case AuthPacketType.TokenAck: {
        if (authPacket.data !== player.token) {
          client.end("§cAuthentication failed, do you have the latest version of VaultMapper installed?");
          console.log(`Authentication failed for ${client.username} (invalid TokenAck)`);
          break;
        }

        client.end("§aAuthentication successful");
        console.log(`Authentication successful for ${client.username}`);
        break;
      }
    }
  });

  client.on("end", () => {
    // cancel the timers if the client disconnects
    clearInterval(kicktimer);
    clearInterval(chatMessageTimer);
  });
});

function createCustomPayloadPacket(packetId: number, data: AuthPacket): Buffer {
  const json = JSON.stringify(data);
  return Buffer.from([packetId, ...Buffer.from(json)]);
}

interface AuthPacket {
  type: AuthPacketType;
  data: string;
}

enum AuthPacketType {
  Token = "token",
  TokenAck = "token_ack",
}
