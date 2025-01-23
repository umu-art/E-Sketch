import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { PaintService } from './service/paint.service';
import { TokenService } from './service/token.service';
import { DispatcherService } from './service/dispatcher.service';

@Module({
  imports: [],
  controllers: [AppController],
  providers: [PaintService, TokenService, DispatcherService],
})
export class AppModule {
}
