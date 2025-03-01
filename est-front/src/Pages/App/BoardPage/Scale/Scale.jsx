import React, { useCallback, useEffect, useRef, useState } from 'react';

import { useDispatch, useSelector } from 'react-redux';
import store from '../../../../redux/store';
import { setOffset, setScale } from '../../../../redux/toolSettings/actions';

import { MAX_SCALE, MIN_SCALE } from '../Board/Drawing/Constants';

const round = Math.round;

const baseC = 1.1;

const TextAndRect = ({ angle, scale, onChange }) => {
  const active = scale < MAX_SCALE * 100 && scale > MIN_SCALE  * 100;
  
  return <g transform={`rotate(${angle} 170 170)`} style={{pointerEvents: 'all', cursor: active ? "pointer" : "auto"}} onClick={() => {onChange(scale)}}>
    <text fill={active ? "black" : "gray"} x="63" y="175" fontFamily="Arial" fontSize="15">
      {scale}%
    </text>
    <rect x="115" y="170" width="7" height="0.5" rx="1" fill="black" stroke="black"/>
  </g>
};

const ScaleElement = ({ style }) => {
  const [currentScale, setCurrentScale] = useState(100);
  const [referenceScale, setReferenceScale] = useState(100);
  const [rotate, setRotate] = useState(0);

  const dispatch = useDispatch();
  const drawingScale = useSelector((state) => state.view.scale);

  const angles = [105, 85, 65, 45, 25, 5, -15];
  const scales = [
    round(referenceScale * baseC ** 3),
    round(referenceScale * baseC ** 2),
    round(referenceScale * baseC ** 1),
    round(referenceScale),
    round(referenceScale / baseC ** 1),
    round(referenceScale / baseC ** 2),
    round(referenceScale / baseC ** 3)
  ];

  const changeScale = (scale) => {
    if (scale > MAX_SCALE * 100 || scale < MIN_SCALE * 100) {
      return;
    }

    const state = store.getState()

    const scaleChange = state.view.scale / (scale / 100);

    const centerX = (window.innerWidth / 2) / state.view.scale;
    const centerY = (window.innerHeight / 2) / state.view.scale;

    const nowX = centerX - state.view.offsetX;
    const nowY = centerY - state.view.offsetY;

    const newOffsetX = -(nowX - centerX / scaleChange);
    const newOffsetY = -(nowY - centerY / scaleChange);

    dispatch(setScale(scale / 100));
    dispatch(setOffset(newOffsetX, newOffsetY));
  }

  const recalculateDefs = useCallback((scale) => {
    let lReferenceScale = 100;
    while (lReferenceScale < scale) {
      lReferenceScale *= baseC;
    }
    while (lReferenceScale > scale) {
      lReferenceScale /= baseC;
    }

    let uReferenceScale = lReferenceScale * baseC;
    let newRotate = (uReferenceScale - scale) / (uReferenceScale - lReferenceScale) * 20;

    return { rotate: newRotate, referenceScale: uReferenceScale };
  }, []);

  const updateScale = useCallback(() => {
    const newScale = drawingScale * 100;
    setCurrentScale(newScale);
    const { rotate: newRotate, referenceScale: newReferenceScale } = recalculateDefs(newScale);
    setRotate(newRotate);
    setReferenceScale(newReferenceScale);
  }, [recalculateDefs, drawingScale]);

  const [position, setPosition] = useState({ x: 0, y: 0 });
  const lastScaleChangeTime = useRef(Date.now());
  const lastScale = useRef(drawingScale);

  const updatePosition = useCallback(() => {
    const newScale = drawingScale * 100;
    if (newScale !== lastScale.current) {
      lastScaleChangeTime.current = Date.now();
      lastScale.current = newScale;
    }

    if (Date.now() - lastScaleChangeTime.current > 3000) {
      setPosition({ x: 40, y: 40 });
    } else {
      setPosition({ x: 0, y: 0 });
    }
  }, [drawingScale]);

  useEffect(() => {
    updateScale();
    updatePosition();

    const updateScaleIntervalId = setInterval(updateScale, 10);
    const updatePositionIntervalId = setInterval(updatePosition, 100);

    return () => {
      clearInterval(updateScaleIntervalId);
      clearInterval(updatePositionIntervalId);
    };
  }, [updateScale, updatePosition]);

  return (
    <div 
      style={{
        ...style,
        height: '170px',
        width: '170px',
        userSelect: 'none',
        pointerEvents: 'none',
        transition: 'transform 0.2s ease-in-out',
        transform: `translate(${position.x}px, ${position.y}px)`,
      }}
    >
      <svg width="170" height="170" viewBox="0 0 170 170" xmlns="http://www.w3.org/2000/svg">
        <svg style={{pointerEvents: 'all'}} x="51" y="51" width="250" height="250" viewBox="0 0 115 115"
          filter="drop-shadow(0 2px 3px rgba(0, 0, 0, 0.2))"
          xmlns="http://www.w3.org/2000/svg">
          <path fillRule="evenodd" clipRule="evenodd"
                d="M57.5 115C89.2563 115 115 89.2563 115 57.5C115 25.7437 89.2563 0 57.5 0C25.7437 0 0 25.7437 0 57.5C0 89.2563 25.7437 115 57.5 115ZM58 81C71.2549 81 82 70.2549 82 57C82 43.7451 71.2549 33 58 33C44.7451 33 34 43.7451 34 57C34 70.2549 44.7451 81 58 81Z"
                fill="white" />
        </svg>  
        <g transform={`rotate(${rotate} 170 170)`}>
          {angles.map((angle, index) => <TextAndRect angle={angle} scale={scales[index]} onChange={changeScale} key={"key-" + index}/>)}
        </g>
        <g transform="rotate(45 170 170)">
          <rect x="55" y="159" width="60" height="22" rx="4.5" fill="white" stroke="black" />
          <text fill="black" x="63" y="175" fontFamily="Arial" fontSize="15">{round(currentScale)}%</text>
          <rect x="115" y="169" width="10" height="1.5" rx="1" fill="black" stroke="black" />
        </g>
      </svg>
    </div>
  );
};

export default ScaleElement;