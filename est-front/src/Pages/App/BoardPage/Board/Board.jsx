import React, { useEffect, useRef } from 'react';
import { Board as BoardController } from 'paint/dist';
import { registerDrawListener } from './Paint';
import { decode, encode } from 'coder/dist';
import { FigureType, Line, Point } from 'figures/dist';

const Board = ({ className, style, boardId }) => {
  const width = '100vw';
  const height = '100vh';

  const boardControllerRef = useRef(null);

  useEffect(() => {
    if (boardControllerRef.current)
      return;
    const webSocket = new WebSocket('wss://' + window.location.host + '/proxy/ws?boardId=' + boardId);

    webSocket.addEventListener('message', (event) => {
      if (event.data.length === 36) { // uuid - ответ на запрос на создание фигуры
        const figure = new Line(
          FigureType.LINE,
          event.data,
          [],
          [new Point(12, 21)],
        );
        webSocket.send(encode(figure)); // Отправляем на сервер фигуру
        return;
      }

      const figure = decode(event.data);
      console.log('Figure from server:', figure); // Чет получили
    });

    webSocket.addEventListener('open', () => {
      console.log("WebSocket connected");
      webSocket.send(String.fromCharCode(0)); // Запрос на создание фигуры
    });


    const boardElement = document.getElementById(boardId);
    boardControllerRef.current = new BoardController(boardElement);
    registerDrawListener(boardElement, boardControllerRef.current);
  }, [boardId]);

  return (<svg
    id={boardId}
    width={width}
    height={height}
    className={className}
    style={{ ...style, backgroundColor: 'white', overflow: 'hidden', zIndex: 5 }}
  ></svg>);
};

export default Board;