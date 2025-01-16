import { Controller, Get, Query, Res, UnauthorizedException } from '@nestjs/common';
import { Response } from 'express';
import { AppService } from './app.service';
import { Mutex } from 'async-mutex';
import * as crypto from 'crypto';

class TokenData {
  token: string;
  timestamp: number;
}

@Controller()
export class AppController {

  private tokenMap = new Map<string, TokenData[]>();
  private tokenMutex = new Mutex();

  constructor(private readonly appService: AppService) {
    setInterval(() => this.clearExpiredTokens(), 60 * 1000);
  }

  @Get('/actuator')
  async actuator() {
    return 'OK';
  }

  @Get('/preview')
  async getPreview(
    @Query('boardId') boardId: string,
    @Query('token') token: string,
    @Res() res: Response,
  ) {
    await this.checkToken(boardId, token);

    const jpegBuffer = await this.appService.getPreview(boardId, 300, 200);
    res.contentType('image/jpeg');
    res.send(jpegBuffer);
  }

  @Get('/internal/preview')
  async getPreviewInternal(
    @Query('boardId') boardId: string,
    @Query('width') width: number,
    @Query('height') height: number,
    @Query('xLeft') xLeft: number,
    @Query('yUp') yUp: number,
    @Query('xRight') xRight: number,
    @Query('yDown') yDown: number,
    @Res() res: Response,
  ) {
    const jpegBuffer = await this.appService.getPreviewPart(boardId, width, height, xLeft, yUp, xRight, yDown);
    res.contentType('image/jpeg');
    res.send(jpegBuffer);
  }

  private async checkToken(boardId: string, token: string) {
    let isValidToken = false;

    await this.tokenMutex.runExclusive(() => {
      const storedTokens = this.tokenMap.get(boardId) || [];
      const validToken = storedTokens.find(t => t.token === token && !this.isTokenExpired(t.timestamp));
      if (validToken) {
        isValidToken = true;
        this.tokenMap.set(boardId, storedTokens.filter(t => t.token !== token));
      }
    });

    if (!isValidToken) {
      throw new UnauthorizedException('Invalid or expired token');
    }
  }

  @Get('/internal/get-token')
  async getToken(@Query('boardId') boardId: string): Promise<string> {
    const token = crypto.randomBytes(32).toString('hex');

    await this.tokenMutex.runExclusive(() => {
      const tokens = this.tokenMap.get(boardId) || [];
      tokens.push({ token, timestamp: Date.now() });
      this.tokenMap.set(boardId, tokens);
    });

    return token;
  }

  private async clearExpiredTokens(): Promise<void> {
    await this.tokenMutex.runExclusive(() => {
      for (const [boardId, tokens] of this.tokenMap.entries()) {
        const validTokens = tokens.filter(t => !this.isTokenExpired(t.timestamp));
        this.tokenMap.set(boardId, validTokens);
      }
    });
  }

  private isTokenExpired(timestamp: number): boolean {
    const tokenAge = Date.now() - timestamp;
    return tokenAge > 5 * 60 * 1000;
  }
}