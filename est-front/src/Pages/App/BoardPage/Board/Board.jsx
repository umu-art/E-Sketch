import React, { useEffect, useRef } from 'react';
import { Board as BoardController } from 'paint/dist';
import { registerDrawListener } from './Paint';
import { connectToFigures, connectToMarkers, figureWebSocket } from './SocketApi';
import { registerMarkersListener } from './Markers';

const Board = ({ className, style, boardId }) => {
  const width = '100vw';
  const height = '100vh';

  const boardControllerRef = useRef(null);

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
    width={width}
    height={height}
    className={className}
    style={{ ...style, backgroundColor: 'white', overflow: 'hidden', zIndex: 5 }}
  ></svg>);
};

export default Board;