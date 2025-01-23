import { Injectable, Logger } from '@nestjs/common';
import { PaintService } from './paint.service';
import { getModifiedBoardIds } from '../api/back.api';
import { load, save } from '../api/s3.api';

@Injectable()
export class DispatcherService {

  private readonly logger = new Logger(DispatcherService.name);

  constructor(private readonly paintService: PaintService) {
    setInterval(this.syncPreviews.bind(this), 60 * 1000);
  }

  public async getPreview(boardId: string): Promise<Buffer> {
    let onS3 = await load(boardId);
    if (onS3) {
      return onS3;
    }

    return this.paintService.getPreview(boardId, 500, 300);
  }

  async syncPreviews() {
    const modifiedBoardIds = await getModifiedBoardIds(1);
    this.logger.log('syncing previews for boards:', modifiedBoardIds.length);
    const batchSize = 5;

    for (let i = 0; i < modifiedBoardIds.length; i += batchSize) {
      const batch = modifiedBoardIds.slice(i, i + batchSize);
      await Promise.all(batch.map(boardId => this.repaintPreview(boardId)));
    }
  }

  async repaintPreview(boardId: string): Promise<void> {
    let preview = await this.paintService.getPreview(boardId, 500, 300);
    await save(boardId, preview);
  }
}