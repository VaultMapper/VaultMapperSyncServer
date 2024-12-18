const SERVER_HOST = "127.0.0.1";
const SERVER_PORT = "25284";
const CONNECTIONS = 1000; // Number of concurrent connections
const MESSAGES_PER_CONNECTION = 10000; // Number of messages per connection
const CONNECTIONS_PER_VAULT = 12; // Number of connections per vault
const VAULT_IDS: string[] = [];

for (let i = 0; i < Math.ceil(CONNECTIONS / CONNECTIONS_PER_VAULT); i++) {
  VAULT_IDS.push(`vault_${crypto.randomUUID()}`);
}

async function stressTest() {
  const tasks: Promise<void>[] = [];

  for (let i = 0; i < CONNECTIONS; i++) {
    tasks.push(simulateClient(i));
  }

  await Promise.all(tasks);
}

async function simulateClient(id: number) {
  await null;

  const uuid = crypto.randomUUID();
  const VAULT_ID = VAULT_IDS[id % VAULT_IDS.length];
  const url = `ws://${SERVER_HOST}:${SERVER_PORT}/${VAULT_ID}/${uuid}`;
  const ws = new WebSocket(url);
  let color = "";
  for (let i = 0; i < 6; i++) {
    color += Math.floor(Math.random() * 16).toString(16);
  }

  ws.onopen = async () => {
    console.log(`Client ${id} connected`);

    for (let i = 0; i < MESSAGES_PER_CONNECTION; i++) {
      const moveCapsule = {
        type: 3,
        data: {
          uuid,
          color,
          x: Math.floor(Math.random() * 48) - 24,
          z: Math.floor(Math.random() * 48) - 24,
          yaw: Math.random() * 360,
        },
      };

      ws.send(JSON.stringify(moveCapsule));
      const cellCapsule = {
        type: 2,
        data: {
          x: Math.floor(Math.random() * 48) - 24,
          z: Math.floor(Math.random() * 48) - 24,
          i: Math.random() < 0.5,
          c: Math.floor(Math.random() * 4),
          r: Math.floor(Math.random() * 3),
          n: Math.floor(Math.random() * 3),
          m: Math.random() < 0.5,
        },
      };
      ws.send(JSON.stringify(cellCapsule));
      await new Promise((resolve) => setTimeout(resolve, Math.random() * 1000));
    }

    ws.close();
  };

  ws.onerror = (error) => {
    const errorMessage = (error as ErrorEvent).message;
    console.error(`Client ${id} error: ${errorMessage}`);
  };
}

stressTest();
