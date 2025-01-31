import { DefaultFigure, FigureHeader, FigureHeaderBuilder, FigureType, getFromHeader, getIntFromHeader } from './default';
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

  static startProcess(type: FigureType, id: string, header: FigureHeader, point: Point) {
    return new Ellipse(type, id, header, [point, point]);
  }
  
  process(cursor: Point): void {
    this.points[1] = cursor;
  };

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
    const centerX = (this.points[0].x + this.points[1].x) / 2;
    const centerY = (this.points[0].y + this.points[1].y) / 2;
    const radiusX = Math.abs(this.points[1].x - this.points[0].x) / 2;
    const radiusY = Math.abs(this.points[1].y - this.points[0].y) / 2;

    let element = document.createElementNS('http://www.w3.org/2000/svg', 'ellipse');
    element.setAttribute('id', this.id);
    element.setAttribute('fill', this.fillColor || 'white');
    element.setAttribute('stroke', this.color || 'black');
    if (this.thickness) {
      element.setAttribute('stroke-width', this.thickness.toString());
    }

    element.setAttribute('cx', centerX.toString());
    element.setAttribute('cy', centerY.toString());
    element.setAttribute('rx', radiusX.toString());
    element.setAttribute('ry', radiusY.toString());

    return element;
  }
}
