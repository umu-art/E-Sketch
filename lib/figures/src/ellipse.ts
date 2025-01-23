import { DefaultFigure, FigureHeader, FigureHeaderBuilder, getFromHeader, getIntFromHeader } from './default';
import { Point } from './point';

export class Ellipse extends DefaultFigure {

  color?: string;
  fillColor?: string;
  thickness?: number;

  protected parseHeader(header: FigureHeader): void {
    this.color = getFromHeader(header, 0);
    this.fillColor = getFromHeader(header, 1);
    this.thickness = getIntFromHeader(header, 2);
  }

  public exportHeader(): FigureHeader {
    return new FigureHeaderBuilder()
      .add(this.color)
      .add(this.fillColor)
      .add(this.thickness?.toString())
      .build();
  }

  public set(center: Point, radiusX: number, radiusY: number): void {
    this.points = [center, new Point(radiusX, radiusY)];
  }

  public toSvg(document: Document): SVGEllipseElement {
    let element = document.createElementNS('http://www.w3.org/2000/svg', 'ellipse');
    element.setAttribute('id', this.id);
    element.setAttribute('fill', this.fillColor || 'white');
    element.setAttribute('stroke', this.color || 'black');

    const [center, radius] = this.points;
    element.setAttribute('cx', center.x.toString());
    element.setAttribute('cy', center.y.toString());
    element.setAttribute('rx', radius.x.toString());
    element.setAttribute('ry', radius.y.toString());

    return element;
  }
}
