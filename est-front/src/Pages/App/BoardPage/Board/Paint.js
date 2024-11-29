import { Board as BoardController } from 'paint/dist';
import { Point } from 'figures/dist/point';
import { FigureType, Line } from 'figures/dist';
import { v4 } from 'uuid';

const FPS = 60;

export function registerDrawListener(board: Element, boardController: BoardController) {
  let drawing = {
    isDrawing: false,
    nowX: 0,
    nowY: 0,
    lineColor: 'red',
    lineWidth: 2,
  };

  let currentFigure: Line;

  board.addEventListener('mousedown', () => {
    drawing.isDrawing = true;
    currentFigure = new Line(FigureType.LINE, v4(), [drawing.lineColor, drawing.lineWidth], []);
  });

  board.addEventListener('mouseup', () => {
    drawing.isDrawing = false;
    boardController.upsertFigure(currentFigure);
  });

  board.addEventListener('mousemove', (event) => {
    const rect = board.getBoundingClientRect();
    drawing.nowX = event.offsetX - rect.left;
    drawing.nowY = event.offsetY - rect.top;
  });

  setInterval(() => {
    if (drawing.isDrawing) {
      currentFigure.points.push(new Point(drawing.nowX, drawing.nowY));
      boardController.upsertFigure(currentFigure);
    }
  }, 1000 / FPS);
}