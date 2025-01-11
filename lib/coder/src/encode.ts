import { DefaultFigure } from 'figures/dist';
import { Point } from 'figures/dist/point';

/**
 * Encodes a figure object into a byte string.
 *
 * The encoding process includes:
 * 1. The figure type
 * 2. The figure id
 * 3. The encoded header
 * 4. All points of the figure
 *
 * @param figure - Any figure, extended from DefaultFigure.
 *
 * @returns A string where each character represents a byte of the encoded figure.
 */
export function encode(figure: DefaultFigure): string {
  let res = '';

  res += String.fromCharCode(figure.type);
  res += figure.id;

  let header = figure.exportHeader();
  res += encodeHeader(header);

  for (const point of figure.points) {
    res += encodePoint(point);
  }

  return res;
}

/**
 * Encodes the header.
 *
 * @param header - An array of strings representing the header information of the figure.
 *
 * @returns A string containing the encoded header information. The string starts with
 *          a character representing the length of the header array, followed by the
 *          header elements joined with a '|' character.
 */
function encodeHeader(header: string[]) {
  let res = '';
  res += String.fromCharCode(header.length);
  res += header.join('|');
  res += '|';
  return res;
}


/**
 * Encodes a point object.
 *
 * @param point - The point object to be encoded.
 *
 * @returns A string containing the binary representation of the point's x and y coordinates.
 */
export function encodePoint(point: Point): string {
  return floatToBinary(point.x) + floatToBinary(point.y);
}


/**
 * Converts a floating-point number to its binary representation as a string.
 *
 * This function takes a JavaScript number (float) and converts it to a binary
 * representation using a 64-bit float (double precision) format. The resulting
 * binary data is then converted to a string where each byte is represented by
 * its corresponding ASCII character.
 *
 * @param float - The floating-point number to be converted to binary.
 *
 * @returns A string where each character represents a byte of the binary
 *          representation of the input float. The string has a length of 8
 *          characters, corresponding to the 8 bytes of a 64-bit float.
 */
function floatToBinary(float: number): string {
  const buffer = new ArrayBuffer(8);
  const view = new DataView(buffer);

  view.setFloat64(0, float);

  let binaryStr = '';
  for (let i = 0; i < 8; i++) {
    binaryStr += String.fromCharCode(view.getUint8(i));
  }

  return binaryStr;
}