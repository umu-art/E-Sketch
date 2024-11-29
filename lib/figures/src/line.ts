import { DefaultFigure, FigureHeader, FigureHeaderBuilder, getFromHeader, getIntFromHeader } from './default';
import { toSvgPath } from './point-utils';

export class Line extends DefaultFigure {

  color?: string;
  thickness?: number;

  public parseHeader(header: FigureHeader): void {
    this.color = getFromHeader(header, 0);
    this.thickness = getIntFromHeader(header, 1);
  }

  public exportHeader(): FigureHeader {
    return new FigureHeaderBuilder()
      .add(this.color)
      .add(this.thickness?.toString())
      .build();
  }

  public toSvg(): SVGPathElement {
    let element = document.createElementNS('http://www.w3.org/2000/svg', 'path');
    element.setAttribute('id', this.id);
    element.setAttribute('fill', 'none');

    const path = toSvgPath(this.points);
    element.setAttribute('d', path);

    if (this.color) {
      element.setAttribute('stroke', this.color);
    } else {
      element.setAttribute('stroke', 'black');
    }

    if (this.thickness) {
      element.setAttribute('stroke-width', this.thickness.toString());
    } else {
      element.setAttribute('stroke-width', '1');
    }

    return element;
  }
}