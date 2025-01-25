import React, { useEffect, useRef, useState } from 'react';
import { Board as BoardController } from 'paint/dist';
import { registerDrawListener } from './Paint';
import { connectToFigures, connectToMarkers, figureWebSocket } from './SocketApi';
import { registerMarkersListener } from './Markers';

import pencilIcon from './pencil.svg';
import eraserIcon from './eraser.svg';

const toolsIcons = {
  'pencil': pencilIcon,
  'eraser': eraserIcon,
};

const Board = ({ className, style, boardId, currentTool }) => {
  const boardControllerRef = useRef(null);
  const [cursor, setCursor] = useState(currentTool);

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
      registerDrawListener(boardElement, boardControllerRef.current);
    });

    registerMarkersListener(boardElement);
  }, [boardId]);

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