import { createClient } from 'redis';

const redisClient = createClient({
  url: process.env.REDIS_URL || 'redis://localhost:6379',
  password: process.env.REDIS_PASSWORD || undefined,
});

redisClient.on('error', (err) => console.error('Redis Client Error', err));

export async function saveToken(token: TokenData): Promise<void> {
  try {
    await redisClient.connect();
    await redisClient.set(`token:${token.boardId}:${token.token}`, 'valid', {
      EX: 300,
    });
  } catch (error) {
    console.error('Error storing token in Redis:', error);
    throw new Error('Failed to create token');
  } finally {
    await redisClient.disconnect();
  }
}

export interface TokenData {
  token: string;
  boardId: string;
}

export async function saveTokens(tokens: TokenData[]): Promise<void> {
  try {
    await redisClient.connect();
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
  } finally {
    await redisClient.disconnect();
  }
}

export async function checkToken(token: TokenData): Promise<boolean> {
  try {
    await redisClient.connect();
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
  } finally {
    await redisClient.disconnect();
  }
}