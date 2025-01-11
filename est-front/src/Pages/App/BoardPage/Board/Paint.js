import { Board as BoardController } from 'paint/dist';
import { Point } from 'figures/dist/point';
import { FigureType, Line } from 'figures/dist';
import { changeFigure, createFigure, getAllFigures, onNewFigure, onUpdateFigure, updateFigure } from './SocketApi';
import { decode, encode } from 'coder/dist';


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
  let oldCurrentFigure: Line;

  board.addEventListener('mousedown', () => {
    createFigure((uuid) => {
      currentFigure = new Line(FigureType.LINE, uuid, [drawing.lineColor, drawing.lineWidth], []);
      oldCurrentFigure = null;
      drawing.isDrawing = true;
    });
  });

  board.addEventListener('mouseup', () => {
    if (!currentFigure) {
      return;
    }

    drawing.isDrawing = false;
    boardController.upsertFigure(currentFigure);

    triggerUpdateFigure(currentFigure, oldCurrentFigure);
    oldCurrentFigure = cloneLine(currentFigure);
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

      triggerUpdateFigure(currentFigure, oldCurrentFigure);
      oldCurrentFigure = cloneLine(currentFigure);
    }
  }, 1000 / FPS);

  onNewFigure((figure) => {
    if (!currentFigure || figure.id !== currentFigure.id) {
      boardController.upsertFigure(figure);
    }
  });

  onUpdateFigure((id, data) => {
    if (!currentFigure || id !== currentFigure.id) {
      const updatableFigure = boardController.figures.find((figure) => figure.id === id);
      if (updatableFigure) {
        const encoded = encode(updatableFigure);
        const decoded = decode(encoded + data);
        boardController.upsertFigure(decoded);
      }
    }
  });

  getAllFigures();
}

function triggerUpdateFigure(newFigure, oldFigure) {
  if (!oldFigure) {
    changeFigure(newFigure);
  } else {
    updateFigure(newFigure, oldFigure);
  }
}

function cloneLine(line: Line): Line {
  const newPoints = line.points.map((point) => new Point(point.x, point.y));
  return new Line(line.type, line.id, [line.color, line.thickness], newPoints);
}