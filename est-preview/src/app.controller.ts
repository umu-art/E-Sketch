import { Controller, Get, Query, Res } from '@nestjs/common';
import { Response } from 'express';
import { PaintService } from './service/paint.service';
import { TokenService } from './service/token.service';
import { TokenData } from './api/redis.api';
import { DispatcherService } from './service/dispatcher.service';

@Controller()
export class AppController {


  constructor(private readonly dispatcherService: DispatcherService,
              private readonly paintService: PaintService,
              private readonly tokenService: TokenService) {
  }

  @Get('/actuator')
  async actuator() {
    return 'OK';
  }

  @Get('/internal/get-token')
  async getToken(@Query('boardId') boardId: string): Promise<string> {
    return await this.tokenService.createToken(boardId);
  }

  @Get('/internal/get-tokens')
  async getTokens(@Query('boardIds') boardIds: string): Promise<TokenData[]> {
    return await this.tokenService.createTokens(boardIds.split(','));
  }

  @Get('/preview')
  async getPreview(
    @Query('boardId') boardId: string,
    @Query('token') token: string,
    @Res() res: Response,
  ) {
    await this.tokenService.checkToken(boardId, token);

    const jpegBuffer = await this.dispatcherService.getPreview(boardId);
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
    const jpegBuffer = await this.paintService.getPreviewPart(boardId, width, height, xLeft, yUp, xRight, yDown);
    res.contentType('image/jpeg');
    res.send(jpegBuffer);
  }
}