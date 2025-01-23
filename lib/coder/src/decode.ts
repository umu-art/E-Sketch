import { Ellipse, DefaultFigure, FigureType, Line, Rectangle } from 'figures/dist';
import { Point } from 'figures/dist/point';

/**
 * Decodes a figure from byte-string.
 *
 * The decoding process includes:
 * 1. Extracting the figure type
 * 2. Extracting the figure id
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
    index += 16;
  }

  switch (type) {
    case FigureType.RECTANGLE:
      return new Rectangle(type, id, header, points);
    case FigureType.LINE:
      return new Line(type, id, header, points);
    case FigureType.ELLIPSE:
      return new Ellipse(type, id, header, points);
    default:
      throw new TypeError('Invalid figure type');
  }
}

/**
 * Extracts figure id.
 *
 * @param encodedString - The encoded string representation of a figure.
 * @param startIndex - The index in the encodedString where the id starts.
 *
 * @returns The extracted id as a string.
 */
function extractId(encodedString: string, startIndex: number): string {
  const UUID_LENGTH = 36;
  return encodedString.slice(startIndex, startIndex + UUID_LENGTH);
}

/**
 * Decodes the header.
 *
 * @param encodedHeader - The encoded string representation of the header.
 *
 * @returns A tuple containing the decoded header array and its length in the encoded string.
 */
function decodeHeader(encodedHeader: string): [string[], number] {
  const headerElemsCount = encodedHeader.charCodeAt(0);
  const headerEnd = findHeaderEndIndex(encodedHeader, headerElemsCount);
  const headerString = encodedHeader.slice(1, headerEnd);
  const header = headerString.split('|');
  return [header, headerEnd + 1];
}

/**
 * Finds the end index of the header in the encoded string.
 *
 * @param encodedString - The encoded string containing the header.
 * @param headerElemsCount - The number of elements in the header.
 * @returns The index in the encoded string where the header ends.
 */
function findHeaderEndIndex(encodedString: string, headerElemsCount: number): number {
  let result = 1;
  while (headerElemsCount !== 0) {
    if (encodedString.charAt(result) === '|') {
      headerElemsCount--;
    }
    result++;
  }
  return result - 1;
}

/**
 * Decodes a point from its binary string.
 *
 * @param encodedPoint - The encoded string representation of a point.
 *
 * @returns A Point object with x and y coordinates.
 */
export function decodePoint(encodedPoint: string): Point {
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