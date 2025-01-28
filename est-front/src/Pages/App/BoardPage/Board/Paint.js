import { Point } from 'figures/dist/point';
import { FigureType, Line } from 'figures/dist';
import { changeFigure, createFigure, deleteFigure, getAllFigures, onNewFigure, onUpdateFigure, updateFigure } from './SocketApi';
import { decode, encode } from 'coder/dist';

import * as d3 from 'd3';

const FPS = 60;
export const BASE_OFFSET_X = 0;
export const BASE_OFFSET_Y = 20;

export const DrawingStates = {
  IDLE: 'idle',
  CREATING: 'creating',
  DRAWING: 'drawing',
};


export const drawing = {
  tool: 'pencil',
  state: DrawingStates.IDLE,
  nowX: 0,
  nowY: 0,
  lineColor: 'red',
  lineWidth: 2,
  scale: 1,
  offsetX: 0,
  offsetY: 0,
};

export function registerDrawListener(board, boardController) {
  let currentFigure;
  let oldCurrentFigure;
  let isMoving = false;

  let isMouseDown = false;

  d3.select('.board')
    .on('mousedown', function(event) {
      if (event.button === 0) {
        isMouseDown = true;
      }
    })
    .on('mouseup', function() {
      isMouseDown = false;
    })
    .on('mousemove', function(event) {
      const target = event.target;

      if (target.tagName !== 'path')
        return;

      if (drawing.tool !== 'eraser' || !isMouseDown)
        return;

      const pathId = d3.select(target).attr('id');

      boardController.removeFigure(pathId);
      deleteFigure(pathId);
    });

  board.addEventListener('mousedown', handleMouseDown);
  board.addEventListener('mouseup', handleMouseUp);
  board.addEventListener('mousemove', handleMouseMove);
  board.addEventListener('wheel', handleWheel);
  board.addEventListener('contextmenu', (e) => e.preventDefault());

  setInterval(updateDrawing, 1000 / FPS);

  onNewFigure(handleNewFigure);
  onUpdateFigure(handleUpdateFigure);

  getAllFigures();

  function handleMouseDown(e) {
    e.preventDefault();

    if (e.button === 0 && drawing.tool === 'pencil') {
      drawing.state = DrawingStates.CREATING;
      currentFigure = new Line(FigureType.LINE, 'waiting', [drawing.lineColor, drawing.lineWidth], []);
      oldCurrentFigure = null;

      createFigure((uuid) => {
        currentFigure.id = uuid;
        drawing.state = DrawingStates.DRAWING;
      });
    } else if (e.button === 0 && drawing.tool === 'eraser') {
    } else if (e.button === 2) {
      isMoving = true;

      drawing.startX = e.clientX;
      drawing.startY = e.clientY;
    }
  }

  function finishDrawing() {
    drawing.state = DrawingStates.IDLE;
    boardController.upsertFigure(currentFigure);

    triggerUpdateFigure(currentFigure, oldCurrentFigure);
    oldCurrentFigure = currentFigure.clone();
  }

  function waitForDrawingState() {
    if (drawing.state === DrawingStates.DRAWING) {
      finishDrawing();
    } else {
      setTimeout(waitForDrawingState, 50);
    }
  }

  function handleMouseUp(e) {
    e.preventDefault();

    if (e.button === 0 && drawing.tool === 'pencil') {
      switch (drawing.state) {
        case DrawingStates.CREATING:
          waitForDrawingState();
          break;
        case DrawingStates.DRAWING:
          finishDrawing();
          break;
        default:
          break;
      }
    } else if (e.button === 2) {
      isMoving = false;
    }
  }

  function handleMouseMove(event) {
    const rect = board.getBoundingClientRect();
    drawing.nowX = (event.offsetX + BASE_OFFSET_X - rect.left) / drawing.scale - drawing.offsetX;
    drawing.nowY = (event.offsetY + BASE_OFFSET_Y - rect.top) / drawing.scale - drawing.offsetY;

    if (isMoving) {
      const dx = event.clientX - drawing.startX;
      const dy = event.clientY - drawing.startY;

      drawing.offsetX += dx / drawing.scale;
      drawing.offsetY += dy / drawing.scale;

      drawing.startX = event.clientX;
      drawing.startY = event.clientY;

      const newWidth = board.clientWidth / drawing.scale;
      const newHeight = board.clientHeight / drawing.scale;

      board.setAttribute('viewBox', `${-drawing.offsetX} ${-drawing.offsetY} ${newWidth} ${newHeight}`);
    }
  }

  function handleWheel(event) {
    event.preventDefault();

    const scaleChange = event.deltaY < 0 ? 1.01 : 0.99;

    const rect = board.getBoundingClientRect();
    const cursorX = (event.offsetX + BASE_OFFSET_X - rect.left) / drawing.scale;
    const cursorY = (event.offsetY + BASE_OFFSET_Y - rect.top) / drawing.scale;

    drawing.nowX = cursorX - drawing.offsetX;
    drawing.nowY = cursorY - drawing.offsetY;

    drawing.scale *= scaleChange;

    const newWidth = board.clientWidth / drawing.scale;
    const newHeight = board.clientHeight / drawing.scale;

    drawing.offsetX = -(drawing.nowX - cursorX / scaleChange);
    drawing.offsetY = -(drawing.nowY - cursorY / scaleChange);

    board.setAttribute('viewBox', `${-drawing.offsetX} ${-drawing.offsetY} ${newWidth} ${newHeight}`);
  }

  function updateDrawing() {
    if (drawing.state === DrawingStates.IDLE) {
      return;
    }

    currentFigure.points.push(new Point(drawing.nowX, drawing.nowY));
    boardController.upsertFigure(currentFigure);

    if (drawing.state === DrawingStates.DRAWING) {
      triggerUpdateFigure(currentFigure, oldCurrentFigure);
      oldCurrentFigure = currentFigure.clone();
    }
  }

  function handleNewFigure(figure) {
    if (!currentFigure || figure.id !== currentFigure.id) {
      boardController.upsertFigure(figure);
    }
  }

  function handleUpdateFigure(id, data) {
    if (!currentFigure || id !== currentFigure.id) {
      const updatableFigure = boardController.figures.find((figure) => figure.id === id);

      if (updatableFigure) {
        const encoded = encode(updatableFigure);
        const decoded = decode(encoded + data);

        boardController.upsertFigure(decoded);
      }
    }
  }
}

function triggerUpdateFigure(newFigure, oldFigure) {
  if (!oldFigure) {
    changeFigure(newFigure);
  } else {
    updateFigure(newFigure, oldFigure);
  }
}