import { Circle, DefaultFigure, FigureType, Line, Rectangle } from 'figures/dist';
import { Point } from 'figures/dist/point';

/**
 * Decodes a byte string representation into a DefaultFigure object.
 *
 * This function takes a string where each character represents a byte and
 * converts it back into a DefaultFigure object. The decoding process includes:
 * 1. Extracting the figure type
 * 2. Extracting the figure ID
 * 3. Decoding the header
 * 4. Decoding all points of the figure
 *
 * @param encodedString - The encoded string representation of a figure.
 *
 * @returns A DefaultFigure object reconstructed from the encoded string.
 */
export function decode(encodedString: string): DefaultFigure {
  let index = 0;

  const type = encodedString.charCodeAt(index++) as FigureType;

  const id = extractId(encodedString, index);
  index += id.length;

  const [header, headerLength] = decodeHeader(encodedString.slice(index));
  index += headerLength;

  const points: Point[] = [];
  while (index < encodedString.length) {
    const point = decodePoint(encodedString.slice(index));
    points.push(point);
    index += 16; // 8 bytes for x + 8 bytes for y
  }

  switch (type) {
    case FigureType.RECTANGLE:
      return new Rectangle(type, id, header, points);
    case FigureType.LINE:
      return new Line(type, id, header, points);
    case FigureType.CIRCLE:
      return new Circle(type, id, header, points);
    default:
      throw new TypeError('Invalid figure type');
  }
}

/**
 * Extracts the ID from an encoded string representation of a figure.
 *
 * This function reads characters from the encoded string starting at the given index
 * until it encounters a null character (ASCII code 0). The extracted characters
 * form the ID of the figure.
 *
 * @param encodedString - The encoded string representation of a figure.
 * @param startIndex - The index in the encodedString where the ID starts.
 *
 * @returns The extracted ID as a string.
 */
function extractId(encodedString: string, startIndex: number): string {
  let id = '';
  while (encodedString.charCodeAt(startIndex) !== 0) {
    id += encodedString[startIndex++];
  }
  return id;
}

/**
 * Decodes the header from the encoded string representation.
 *
 * @param encodedHeader - The encoded string representation of the header.
 *
 * @returns A tuple containing the decoded header array and its length in the encoded string.
 */
function decodeHeader(encodedHeader: string): [string[], number] {
  const headerLength = encodedHeader.charCodeAt(0);
  const headerString = encodedHeader.slice(1, headerLength + 1);
  const header = headerString.split('|');
  return [header, headerLength + 1];
}

/**
 * Decodes a point from its binary string representation.
 *
 * @param encodedPoint - The encoded string representation of a point.
 *
 * @returns A Point object with x and y coordinates.
 */
function decodePoint(encodedPoint: string): Point {
  const x = binaryToFloat(encodedPoint.slice(0, 8));
  const y = binaryToFloat(encodedPoint.slice(8, 16));
  return new Point(x, y);
}

/**
 * Converts a binary string representation back to a floating-point number.
 *
 * @param binaryStr - The binary string representation of a float.
 *
 * @returns The decoded floating-point number.
 */
function binaryToFloat(binaryStr: string): number {
  const buffer = new ArrayBuffer(8);
  const view = new DataView(buffer);

  for (let i = 0; i < 8; i++) {
    view.setUint8(i, binaryStr.charCodeAt(i));
  }

  return view.getFloat64(0);
}