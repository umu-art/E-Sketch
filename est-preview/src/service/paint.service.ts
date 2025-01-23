import { Injectable, Logger } from '@nestjs/common';
import * as JSDOM from 'jsdom';
import { Board } from 'paint/dist';
import { DefaultFigure } from 'figures/dist';
import * as sharp from 'sharp';
import { getAllFigures } from '../api/back.api';

@Injectable()
export class PaintService {

  private readonly logger = new Logger(PaintService.name);

  async getPreview(boardId: string, width: number, height: number): Promise<Buffer> {
    return this.getPreviewPart(boardId, width, height, undefined, undefined, undefined, undefined);
  }

  async getPreviewPart(boardId: string,
                       width: number, height: number,
                       xLeft: number, yUp: number,
                       xRight: number, yDown: number,
  ): Promise<Buffer> {
    let startTime = Date.now();

    return getAllFigures(boardId)
      .then(figures => {
        this.logger.log(`Fetched figures in ${Date.now() - startTime}ms`);
        return this.createSvgWithFigures(figures, width, height, xLeft, yUp, xRight, yDown);
      })
      .then(svg => {
        this.logger.log(`Created SVG in ${Date.now() - startTime}ms`);
        return this.convertSvgToJpeg(svg);
      })
      .then(jpeg => {
        this.logger.log(`Converted to JPEG in ${Date.now() - startTime}ms`);
        return jpeg;
      });
  }

  private async createSvgWithFigures(figures: DefaultFigure[],
                                     width: number, height: number,
                                     xLeft: number, yUp: number,
                                     xRight: number, yDown: number,
  ): Promise<SVGElement> {

    const dom = new JSDOM.JSDOM('<!DOCTYPE html><html><body></body></html>');
    const document = dom.window.document;


    const svgElement = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    svgElement.setAttribute('width', `${width}`);
    svgElement.setAttribute('height', `${height}`);

    if (xLeft && yUp && xRight && yDown) {
      svgElement.setAttribute('viewBox', `${xLeft}, ${yUp}, ${xRight - xLeft}, ${yDown - yUp}`);
    } else {
      const leftUp = this.getLeftUpPoint(figures);
      const rightDown = this.getRightDownPoint(figures);
      svgElement.setAttribute('viewBox', `${leftUp.x}, ${leftUp.y}, ${rightDown.x - leftUp.x}, ${rightDown.y - leftUp.y}`);
    }

    const board = new Board(svgElement);

    figures.forEach(figure => {
      board.upsertFigure(figure);
    });

    return svgElement;
  }

  private getLeftUpPoint(figures: DefaultFigure[]): { x: number, y: number } {
    const minX = Math.min(...figures.flatMap(f => f.points.map(p => p.x)));
    const minY = Math.min(...figures.flatMap(f => f.points.map(p => p.y)));
    return { x: minX, y: minY };
  }

  private getRightDownPoint(figures: DefaultFigure[]): { x: number, y: number } {
    const maxX = Math.max(...figures.flatMap(f => f.points.map(p => p.x)));
    const maxY = Math.max(...figures.flatMap(f => f.points.map(p => p.y)));
    return { x: maxX, y: maxY };
  }

  private async convertSvgToJpeg(svgElement: SVGElement): Promise<Buffer> {
    return sharp(Buffer.from(svgElement.outerHTML))
      .flatten({ background: { r: 255, g: 255, b: 255 } })
      .jpeg({ quality: 90 })
      .toBuffer();
  }
}
