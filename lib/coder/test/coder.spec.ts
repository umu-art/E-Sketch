import { describe, it } from 'node:test';
import { DefaultFigure, FigureType, Line, Point } from 'figures/dist';
import { v4 } from 'uuid';
import { decode, encode } from '../src';
import { strict as assert } from 'node:assert';

describe('test coder', () => {
  it('should encode and decode figures', () => {
    let testFigure = createTestFigure();

    let encodedFigure = encode(testFigure);
    let decodedFigure = decode(encodedFigure);

    assert.deepStrictEqual(decodedFigure, testFigure);
  });

  it('should correctly insert id', () => {
    let testFigure = createTestFigure();

    let encodedFigure = encode(testFigure);
    let id = encodedFigure.slice(1, 37);

    assert.strictEqual(id, testFigure.id);
  });
});

function createTestFigure(): DefaultFigure {
  let testHeaders = ['red', '3'];
  let testPoints = [new Point(0, 0), new Point(10, 10)];
  return new Line(FigureType.LINE, v4(), testHeaders, testPoints);
}