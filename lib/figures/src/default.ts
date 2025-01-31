import { Point } from './point';

export enum FigureType {
  LINE = 0,
  RECTANGLE = 1,
  ELLIPSE = 2,
}

export abstract class DefaultFigure {

  type: FigureType;
  id: string;
  points: Point[];

  public constructor(type: FigureType, id: string, header: FigureHeader, points: Point[]) {
    this.type = type;
    this.id = id;
    this.points = points;
    this.parseHeader(header);
  }

  static startProcess(type: FigureType, id: string, header: FigureHeader, point: Point) {
    return new (this.constructor as any)(type, id, header, [point]);
  };

  abstract process(cursor: Point): void;

  protected abstract parseHeader(_header: FigureHeader): void;

  public abstract exportHeader(): FigureHeader;

  public abstract toSvg(document: Document): SVGPathElement;

  public clone(): DefaultFigure {
    const newPoints = this.points.map((point) => new Point(point.x, point.y));
    return new (this.constructor as any)(this.type, this.id, this.exportHeader(), newPoints);
  }
}

export type FigureHeader = string[];

export class FigureHeaderBuilder {
  private header: FigureHeader = [];

  public add(value: string | undefined): FigureHeaderBuilder {
    if (value) {
      this.header.push(value);
    } else {
      this.header.push('');
    }

    return this;
  }

  public build(): FigureHeader {
    return this.header;
  }
}

export function getFromHeader(header: FigureHeader, index: number): string | undefined {
  if (index >= 0 && index < header.length) {
    return header[index];
  }
  return undefined;
}

export function getIntFromHeader(header: FigureHeader, index: number): number | undefined {
  const value = getFromHeader(header, index);
  if (value && !isNaN(parseInt(value))) {
    return parseInt(value);
  }
  return undefined;
}