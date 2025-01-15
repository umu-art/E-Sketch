import { Point } from './point';

export enum FigureType {
  LINE = 0,
  RECTANGLE = 1,
  CIRCLE = 2,
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

  protected abstract parseHeader(_header: FigureHeader): void;

  public abstract exportHeader(): FigureHeader;

  public abstract toSvg(document: Document): SVGPathElement;
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