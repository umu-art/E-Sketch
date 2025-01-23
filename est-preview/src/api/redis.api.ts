import { createClient } from 'redis';

const redisClient = createClient({
  url: process.env.REDIS_URL || 'redis://localhost:6379',
  password: process.env.REDIS_PASSWORD || undefined,
});

let isConnected = false;
const maxReconnectAttempts = 5;
let reconnectAttempts = 0;

redisClient.on('error', (err) => {
  console.error('Redis Client Error', err);
  if (isConnected) {
    isConnected = false;
    reconnect();
  }
});

redisClient.on('connect', () => {
  console.log('Connected to Redis');
  isConnected = true;
  reconnectAttempts = 0;
});

async function connect() {
  try {
    await redisClient.connect();
  } catch (error) {
    console.error('Failed to connect to Redis:', error);
    reconnect();
  }
}

function reconnect() {
  if (reconnectAttempts < maxReconnectAttempts) {
    reconnectAttempts++;
    console.log(`Attempting to reconnect...`);
    setTimeout(connect, 1000);
  } else {
    console.error('Max reconnection attempts reached. Redis is unavailable.');
    process.exit(4);
  }
}

connect().catch(console.error);

process.on('SIGINT', async () => {
  if (isConnected) {
    await redisClient.disconnect();
  }
  process.exit(0);
});

export async function saveToken(token: TokenData): Promise<void> {
  try {
    await redisClient.set(`token:${token.boardId}:${token.token}`, 'valid', {
      EX: 300,
    });
  } catch (error) {
    console.error('Error storing token in Redis:', error);
    throw new Error('Failed to create token');
  }
}

export interface TokenData {
  token: string;
  boardId: string;
}

export async function saveTokens(tokens: TokenData[]): Promise<void> {
  try {
    const pipeline = redisClient.multi();
    for (const token of tokens) {
      pipeline.set(`token:${token.boardId}:${token.token}`, 'valid', {
        EX: 300,
      });
    }
    await pipeline.exec();
  } catch (error) {
    console.error('Error storing tokens in Redis:', error);
    throw new Error('Failed to save tokens');
  }
}

export async function checkToken(token: TokenData): Promise<boolean> {
  try {
    const key = `token:${token.boardId}:${token.token}`;
    const isValid = await redisClient.get(key);

    if (isValid === 'valid') {
      await redisClient.del(key);
      return true;
    }

    return false;
  } catch (error) {
    console.error('Error checking token in Redis:', error);
    return false;
  }
}