import { DefaultFigure, FigureHeader, FigureHeaderBuilder, FigureType, getFromHeader, getIntFromHeader } from './default';
import { Point } from './point';

export class Rectangle extends DefaultFigure {

  color?: string;
  fillColor?: string;
  thickness?: number;

  static startProcess(type: FigureType, id: string, header: FigureHeader, point: Point) {
    return new Rectangle(type, id, header, [point, point]);
  }
    
  process(cursor: Point): void {
    this.points[1] = cursor;
  };

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

  public set(topLeft: Point, bottomRight: Point): void {
    this.points = [topLeft, bottomRight];
  }

  public toSvg(document: Document): SVGPathElement {
    const [first, second] = this.points;
    const [topLeft, bottomRight] = [new Point(Math.min(first.x, second.x), Math.min(first.y, second.y)), new Point(Math.max(first.x, second.x), Math.max(first.y, second.y))];

    let element = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
    
    element.setAttribute('id', this.id);
    element.setAttribute('fill', this.fillColor || 'white');
    element.setAttribute('stroke', this.color || 'black');
    if (this.thickness) {
      element.setAttribute('stroke-width', this.thickness.toString());
    }

    element.setAttribute('x', topLeft.x.toString());
    element.setAttribute('y', topLeft.y.toString());
    element.setAttribute('width', (bottomRight.x - topLeft.x).toString());
    element.setAttribute('height', (bottomRight.y - topLeft.y).toString());
    
    return element;
  }
}
