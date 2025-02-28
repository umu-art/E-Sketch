import React, { useEffect, useRef, useState } from 'react';
import { Board as BoardController } from 'paint/dist';
import { registerDrawListener } from './Drawing/Paint';
import { registerMarkersListener } from './Markers';

import pencilIcon from './pencil.svg';
import eraserIcon from './eraser.svg';
import { useSelector } from 'react-redux';
import FigureWebSocket from './Sockets/FigureSocket';
import MarkerWebSocket from './Sockets/MarkerSocket';

const toolsIcons = {
  'pencil': pencilIcon,
  'eraser': eraserIcon,
};

const Board = ({ className, style, boardId }) => {
  const boardControllerRef = useRef(null);

  const currentTool = useSelector((state) => state.tool);

  const [cursor, setCursor] = useState(currentTool);

  const drawing = useSelector((state) => state);

  useEffect(() => {
    setCursor(currentTool);
  }, [currentTool]);

  useEffect(() => {
    if (boardControllerRef.current)
      return;

    const boardElement = document.getElementById(boardId);
    boardControllerRef.current = new BoardController(boardElement);

    const figureWebSocket = new FigureWebSocket(boardId);
    const markerWebSocket = new MarkerWebSocket(boardId);

    figureWebSocket.webSocket.addEventListener('open', () => {
      registerDrawListener(boardElement, boardControllerRef.current, drawing, figureWebSocket);
    });

    registerMarkersListener(boardElement, drawing, markerWebSocket);
  }, [boardId, drawing]);

  return (<svg
    id={boardId}
    className={`${className} board`}
    style={{
      ...style,
      backgroundColor: 'white',
      overflow: 'hidden',
      zIndex: 5,
      cursor: `url(${toolsIcons[cursor]}) 0 20, auto`,
    }}
  ></svg>);
};

export default Board;