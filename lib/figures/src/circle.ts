import { DefaultFigure, FigureHeader, FigureHeaderBuilder, getIntFromHeader } from './default';

export class Circle extends DefaultFigure {

  radius?: number;

  protected parseHeader(header: FigureHeader): void {
    this.radius = getIntFromHeader(header, 0);
  }

  public exportHeader(): FigureHeader {
    return new FigureHeaderBuilder()
      .add(this.radius?.toString())
      .build();
  }

  public toSvg(document: Document): SVGCircleElement {
    // TODO: Implement SVG circle creation
    return document.createElementNS('http://www.w3.org/2000/svg', 'circle');
  }
}
