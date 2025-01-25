import React, { useCallback, useEffect, useState } from 'react';
import { boardControllerRef } from '../Board/Board';
import { drawing } from '../Board/Paint';

const scale = 3;

export const Minimap = ({ style, className, boardId }) => {
  const [pointerBox, setPointerBox] = useState({ x: 0, y: 0, width: 0, height: 0 });

  const addToBoardController = useCallback((boardElement) => {
    if (!boardControllerRef.current) {
      setTimeout(() => addToBoardController(boardElement), 100);
      return;
    }
    boardControllerRef.current.addSvgElement(boardElement);
  }, []);

  const updatePointerBox = useCallback((minimap) => {
    setPointerBox({
      width: minimap.clientWidth / scale,
      height: minimap.clientHeight / scale,
      x: (minimap.clientWidth - (minimap.clientWidth / scale)) / 2,
      y: (minimap.clientHeight - (minimap.clientHeight / scale)) / 2,
    });
  }, []);

  const updateMinimap = useCallback((minimap) => {
    let refBoard = document.getElementById(boardId);
    const refWidth = refBoard.clientWidth / drawing.scale;
    const refHeight = refBoard.clientHeight / drawing.scale;

    const newWidth = refWidth * scale;
    const newHeight = refHeight * scale;
    const newX = -drawing.offsetX - (newWidth - refWidth) / 2;
    const newY = -drawing.offsetY - (newHeight - refHeight) / 2;

    minimap.setAttribute('viewBox', `${newX} ${newY} ${newWidth} ${newHeight}`);
  }, [boardId]);

  useEffect(() => {
    const minimapElement = document.getElementById(boardId + '-minimap');
    addToBoardController(minimapElement);
    updatePointerBox(minimapElement);

    const updateMinimapInterval = setInterval(() => updateMinimap(minimapElement), 1000 / 60);
    return () => clearInterval(updateMinimapInterval);
  }, [addToBoardController, updateMinimap, updatePointerBox, boardId]);

  return (
    <div style={{ ...style }} className={className}>
      <svg
        id={boardId + '-minimap'}
        style={{
          backgroundColor: 'rgba(222,222,222,0.73)',
          overflow: 'hidden',
          zIndex: 6,
          pointerEvents: 'none',
          width: '100%',
          height: '100%',
          position: 'absolute',
          top: 0,
          left: 0,
        }}
      ></svg>
      <div
        style={{
          border: '1px solid rgba(0,0,0,0.5)',
          position: 'absolute',
          top: pointerBox.y,
          left: pointerBox.x,
          pointerEvents: 'none',
          width: pointerBox.width,
          height: pointerBox.height,
          zIndex: 7,
        }}>
      </div>
    </div>
  );
};