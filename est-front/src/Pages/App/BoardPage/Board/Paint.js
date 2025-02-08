import { Point } from 'figures/dist/point';
import { Ellipse, FigureType, Line, Rectangle } from 'figures/dist';
import { changeFigure, createFigure, deleteFigure, getAllFigures, onNewFigure, onRemoveFigure, onUpdateFigure, updateFigure } from './SocketApi';
import { decode, encode } from 'coder/dist';

import * as d3 from 'd3';
import store from '../../../../redux/store';

const FPS = 60;
export const BASE_OFFSET_X = 0;
export const BASE_OFFSET_Y = 20;

export const MIN_SCALE = 0.1;
export const MAX_SCALE = 10;

export const DrawingStates = {
  IDLE: 'idle',
  CREATING: 'creating',
  DRAWING: 'drawing',
};

export function registerDrawListener(board, boardController, initialDrawing) {
  let settings = initialDrawing;

  let drawing = {
    nowX: 0,
    nowY: 0,
  };

  let currentFigure;
  let oldCurrentFigure;
  let isMoving = false;

  let isMouseDown = false;

  store.subscribe(() => {
    const newState = store.getState();

    settings = newState;

    updateViewBox();
  });

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

      if (settings.tool !== 'eraser' || !isMouseDown)
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
  onRemoveFigure(handleRemoveFigure);

  getAllFigures();

  const toolToClass = {
    'pencil': Line,
    'rectangle': Rectangle,
    'ellipse': Ellipse,
  };

  const toolToFigureType = {
    'pencil': FigureType.LINE,
    'rectangle': FigureType.RECTANGLE,
    'ellipse': FigureType.ELLIPSE,
  };

  function handleMouseDown(e) {
    e.preventDefault();

    if (e.button === 0 && settings.tool !== 'eraser') {
      settings.state = DrawingStates.CREATING;

      currentFigure = toolToClass[settings.tool].startProcess(toolToFigureType[settings.tool], 'waiting', Object.values(settings.tools[settings.tool]), new Point(drawing.nowX, drawing.nowY));
      oldCurrentFigure = null;

      createFigure((uuid) => {
        currentFigure.id = uuid;
        settings.state = DrawingStates.DRAWING;
      });
    } else if (e.button === 2) {
      isMoving = true;

      settings.view.startX = e.clientX;
      settings.view.startY = e.clientY;
    }
  }

  function finishDrawing() {
    settings.state = DrawingStates.IDLE;
    boardController.upsertFigure(currentFigure);

    triggerUpdateFigure(currentFigure, oldCurrentFigure);
    oldCurrentFigure = currentFigure.clone();
  }

  function waitForDrawingState() {
    if (settings.state === DrawingStates.DRAWING) {
      finishDrawing();
    } else {
      setTimeout(waitForDrawingState, 50);
    }
  }

  function handleMouseUp(e) {
    e.preventDefault();

    if (e.button === 0 && (settings.tool !== 'eraser')) {
      switch (settings.state) {
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
    drawing.nowX = (event.offsetX + BASE_OFFSET_X - rect.left) / settings.view.scale - settings.view.offsetX;
    drawing.nowY = (event.offsetY + BASE_OFFSET_Y - rect.top) / settings.view.scale - settings.view.offsetY;

    if (isMoving) {
      const dx = event.clientX - settings.view.startX;
      const dy = event.clientY - settings.view.startY;

      settings.view.offsetX += dx / settings.view.scale;
      settings.view.offsetY += dy / settings.view.scale;

      settings.view.startX = event.clientX;
      settings.view.startY = event.clientY;

      const newWidth = board.clientWidth / settings.view.scale;
      const newHeight = board.clientHeight / settings.view.scale;

      board.setAttribute('viewBox', `${-settings.view.offsetX} ${-settings.view.offsetY} ${newWidth} ${newHeight}`);
    }
  }

  function handleWheel(event) {
    event.preventDefault();

    const scaleChange = getScaleChange(event.deltaY);
    const cursorPosition = calculateCursorPosition(event);

    updateDrawingPosition(cursorPosition, scaleChange);
    updateViewBox();
  }

  function getScaleChange(deltaY) {
    return deltaY < 0 ? 1.03 : 0.97;
  }

  function calculateCursorPosition(event) {
    const rect = board.getBoundingClientRect();
    const cursorX = (event.offsetX + BASE_OFFSET_X - rect.left) / settings.view.scale;
    const cursorY = (event.offsetY + BASE_OFFSET_Y - rect.top) / settings.view.scale;
    return { x: cursorX, y: cursorY };
  }

  function updateDrawingPosition(cursorPosition, scaleChange) {
    drawing.nowX = cursorPosition.x - settings.view.offsetX;
    drawing.nowY = cursorPosition.y - settings.view.offsetY;

    settings.view.scale *= scaleChange;

    settings.view.scale = Math.min(Math.max(settings.view.scale, MIN_SCALE), MAX_SCALE);

    settings.view.offsetX = -(drawing.nowX - cursorPosition.x / scaleChange);
    settings.view.offsetY = -(drawing.nowY - cursorPosition.y / scaleChange);
  }

  function updateViewBox() {
    const newWidth = board.clientWidth / settings.view.scale;
    const newHeight = board.clientHeight / settings.view.scale;

    board.setAttribute('viewBox', `${-settings.view.offsetX} ${-settings.view.offsetY} ${newWidth} ${newHeight}`);
  }

  function updateDrawing() {
    if (settings.state === DrawingStates.IDLE) {
      return;
    }

    currentFigure.process(new Point(drawing.nowX, drawing.nowY));
    boardController.upsertFigure(currentFigure);

    // if (settings.tool === "pencil") {
    //   currentFigure.points.push(new Point(drawing.nowX, drawing.nowY));
    //   boardController.upsertFigure(currentFigure);
    // } else if (settings.tool === "rectangle") {
    //   currentFigure.points[1].x = drawing.nowX;
    //   currentFigure.points[1].y = drawing.nowY;

    //   boardController.upsertFigure(currentFigure);
    // }


    if (settings.state === DrawingStates.DRAWING) {
      triggerUpdateFigure(currentFigure, oldCurrentFigure);
      oldCurrentFigure = currentFigure.clone();
    }
  }

  function handleNewFigure(figure) {
    if (!currentFigure || figure.id !== currentFigure.id) {
      boardController.upsertFigure(figure);
    }
  }

  function handleRemoveFigure(id) {
    boardController.removeFigure(id);
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