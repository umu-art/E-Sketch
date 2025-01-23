import { Injectable, UnauthorizedException } from '@nestjs/common';
import * as crypto from 'crypto';
import { checkToken, saveToken, saveTokens, TokenData } from '../api/redis.api';

@Injectable()
export class TokenService {

  constructor() {
  }

  public async checkToken(boardId: string, token: string) {
    if (!await checkToken({ boardId, token })) {
      throw new UnauthorizedException('Invalid or expired token');
    }
  }

  public async createToken(boardId: string): Promise<string> {
    const token = crypto.randomBytes(32).toString('hex');
    await saveToken({ boardId, token });
    return token;
  }

  public async createTokens(boardIds: string[]): Promise<TokenData[]> {
    let tokens: TokenData[] = [];
    for (const boardId of boardIds) {
      tokens.push({
        boardId: boardId,
        token: crypto.randomBytes(32).toString('hex'),
      });
    }
    await saveTokens(tokens);
    return tokens;
  }
}