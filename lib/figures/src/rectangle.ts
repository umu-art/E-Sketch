import { DefaultFigure, FigureHeader, FigureHeaderBuilder, getFromHeader, getIntFromHeader } from './default';

export class Rectangle extends DefaultFigure {

  lineThickness?: number;
  lineColor?: string;
  fillColor?: string;

  protected parseHeader(header: FigureHeader): void {
    this.lineThickness = getIntFromHeader(header, 0);
    this.lineColor = getFromHeader(header, 1);
    this.fillColor = getFromHeader(header, 2);
  }

  public exportHeader(): FigureHeader {
    return new FigureHeaderBuilder()
      .add(this.lineThickness?.toString())
      .add(this.lineColor)
      .add(this.fillColor)
      .build();
  }

  public toSvg(): SVGPathElement {
    // TODO: Implement SVG rectangle creation
    return document.createElementNS('http://www.w3.org/2000/svg', 'rect');
  }
}
