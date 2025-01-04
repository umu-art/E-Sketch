import { Board as BoardController } from 'paint/dist';
import { Point } from 'figures/dist/point';
import { FigureType, Line } from 'figures/dist';
import { decode, encode } from 'coder/dist';

const FPS = 60;

export function registerDrawListener(board: Element, boardController: BoardController, webSocket: WebSocket) {
  let drawing = {
    isDrawing: false,
    nowX: 0,
    nowY: 0,
    lineColor: 'red',
    lineWidth: 2,
  };

  let currentFigure: Line;

  board.addEventListener('mousedown', () => {
    requestFigureCreation(webSocket);

    const messageHandler = (event) => {
      if (event.data.length === 36) { // Если получили uuid в ответ
        currentFigure = new Line(FigureType.LINE, event.data, [drawing.lineColor, drawing.lineWidth], []);
        drawing.isDrawing = true;
        webSocket.removeEventListener('message', messageHandler);
      }
    };

    webSocket.addEventListener('message', messageHandler);
  });

  board.addEventListener('mouseup', () => {
    drawing.isDrawing = false;
    boardController.upsertFigure(currentFigure);
    sendFigure(webSocket, currentFigure);
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
      sendFigure(webSocket, currentFigure);
    }
  }, 1000 / FPS);

  webSocket.addEventListener('message', (event) => {
    if (event.data.length !== 36) {
      let figure = decode(event.data);
      if (figure.id !== currentFigure.id) {
        boardController.upsertFigure(figure);
      }
    }
  });
}

function requestFigureCreation(webSocket: WebSocket) {
  webSocket.send(String.fromCharCode(0)); // Запрос на создание фигуры
}

function sendFigure(webSocket: WebSocket, figure: Line) {
  webSocket.send(String.fromCharCode(1) + encode(figure));
}