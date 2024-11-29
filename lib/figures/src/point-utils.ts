import { Point } from './point';

/**
 * Converts an array of points into an SVG path string.
 *
 * This function generates an SVG path string representation of a series of points,
 * using a combination of straight lines and quadratic Bezier curves to create a smooth path.
 *
 * @param points - An array of Point objects representing the path to be converted.
 *                 Each Point object should have x and y coordinates.
 *
 * @returns A string representing the SVG path. The string will be empty if no points are provided.
 *          For a single point, it returns a "move to" command.
 *          For two points, it returns a simple line.
 *          For three or more points, it creates a path with smooth transitions using Bezier curves.
 */
export function toSvgPath(points: Point[]): string {
  if (points.length === 0) {
    return '';
  }

  let result = 'M' + serialize(points[0]);

  if (points.length === 1) {
    return result;
  }

  if (points.length === 2) {
    return result + ' L' + serialize(points[1]);
  }

  // Первая треть отрезка:
  let position = gfm(points[0], points[1]);
  result += ' L' + serialize(position);

  for (let i = 1; i < points.length - 1; i++) {
    // Шаг по кривой Безье:
    const nextPosition = gsm(points[i], points[i + 1]);
    result += ' Q' + serialize(points[i]) + ' ' + serialize(nextPosition);
  }

  // Последняя треть отрезка:
  result += ' L' + serialize(points[points.length - 1]);

  return result;
}

/**
 * Calculates a point at 1/3 of the distance from point1 to point2.
 *
 * @param point1 - The starting point.
 * @param point2 - The ending point.
 * @returns A new Point object representing the calculated point at 1/3 of the distance from point1 to point2.
 */
function gfm(point1: Point, point2: Point): Point {
  return new Point(
    point1.x + (point2.x - point1.x) / 3,
    point1.y + (point2.y - point1.y) / 3,
  );
}

/**
 * Calculates a point at 2/3 of the distance from point1 to point2.
 *
 * @param point1 - The starting point.
 * @param point2 - The ending point.
 * @returns A new Point object representing the calculated point at 2/3 of the distance from point1 to point2.
 */
function gsm(point1: Point, point2: Point): Point {
  return new Point(
    point1.x + (point2.x - point1.x) * 2 / 3,
    point1.y + (point2.y - point1.y) * 2 / 3,
  );
}

/**
 * Serializes a Point object into a string representation.
 *
 * @param point - The Point object to be serialized.
 * @returns A string representation of the point's coordinates, with x and y values
 *          formatted to 6 decimal places and separated by a comma.
 */
function serialize(point: Point) {
  return point.x.toFixed(6) + ',' + point.y.toFixed(6);
}
