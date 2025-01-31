export const SET_TOOL = 'SET_TOOL';
export const SET_STATE = 'SET_STATE';
export const SET_POSITION = 'SET_POSITION';
export const SET_LINE_COLOR = 'SET_LINE_COLOR';
export const SET_LINE_WIDTH = 'SET_LINE_WIDTH';
export const SET_SCALE = 'SET_SCALE';
export const SET_OFFSET = 'SET_OFFSET';
export const SET_FILL_COLOR = 'SET_FILL_COLOR';

export const setTool = (tool) => ({ type: SET_TOOL, payload: tool });
export const setState = (state) => ({ type: SET_STATE, payload: state });
export const setPosition = (nowX, nowY) => ({
    type: SET_POSITION,
    payload: { nowX, nowY },
});
export const setLineColor = (lineColor) => ({ type: SET_LINE_COLOR, payload: lineColor });
export const setFillColor = (fillColor) => ({ type: SET_FILL_COLOR, payload: fillColor });
export const setLineWidth = (lineWidth) => ({ type: SET_LINE_WIDTH, payload: lineWidth });
export const setScale = (scale) => ({ type: SET_SCALE, payload: scale });
export const setOffset = (offsetX, offsetY) => ({
    type: SET_OFFSET,
    payload: { offsetX, offsetY },
});
