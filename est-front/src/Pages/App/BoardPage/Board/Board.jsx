import React, { useEffect, useRef, useState } from 'react';
import { Board as BoardController } from 'paint/dist';
import { registerDrawListener } from './Paint';
import { connectToFigures, connectToMarkers, figureWebSocket } from './SocketApi';
import { registerMarkersListener } from './Markers';

import pencilIcon from './pencil.svg';
import eraserIcon from './eraser.svg';
import { useSelector } from 'react-redux';

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

    connectToFigures(boardId);
    connectToMarkers(boardId);

    figureWebSocket.addEventListener('open', () => {
      registerDrawListener(boardElement, boardControllerRef.current, drawing);
    });

    registerMarkersListener(boardElement, drawing);
  }, [boardId, drawing]);

  return (<svg
    id={boardId}
    className={`${className} board`}
    style={{
      ...style,
      backgroundColor: 'white',
      overflow: 'hidden',
      zIndex: 5,
      cursor: `url(${toolsIcons[cursor]}), auto`,
    }}
  ></svg>);
};

export default Board;