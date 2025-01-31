import {
    SET_TOOL,
    SET_STATE,
    SET_LINE_COLOR,
    SET_LINE_WIDTH,
    SET_SCALE,
    SET_OFFSET,
    SET_FILL_COLOR,
} from './actions';

const initialState = {
    state: 'idle',
    tool: 'pencil',
    tools: {
        'pencil': {
            lineColor: '#1677ff',
            lineWidth: 2,
        },
        'eraser': {
            lineWidth: 5,
        },
        'rectangle': {
            lineColor: '#1677ff',
            fillColor: '#ffffff00',
            lineWidth: 2,
        },
        'ellipse': {
            lineColor: '#1677ff',
            fillColor: '#ffffff00',
            lineWidth: 2,
        },
    },
    view: {
        scale: 1,
        offsetX: 0,
        offsetY: 0,
        startX: 0,
        startY: 0,
    }
};

const minScale = 0.1;
const maxScale = 10;

const drawingReducer = (state = initialState, action) => {
    switch (action.type) {
        case SET_TOOL:
            return { ...state, tool: action.payload };
        case SET_STATE:
            return { ...state, state: action.payload };
        case SET_LINE_COLOR: {
            const currentTool = state.tools[state.tool];
            return {
                ...state,
                tools: {
                    ...state.tools,
                    [state.tool]: {
                        ...currentTool,
                        lineColor: action.payload
                    }
                }
            };
        }
        case SET_FILL_COLOR: {
            const currentTool = state.tools[state.tool];
            return {
                ...state,
                tools: {
                    ...state.tools,
                    [state.tool]: {
                        ...currentTool,
                        fillColor: action.payload
                    }
                }
            };
        }
        case SET_LINE_WIDTH: {
            const currentTool = state.tools[state.tool];
            return {
                ...state,
                tools: {
                    ...state.tools,
                    [state.tool]: {
                        ...currentTool,
                        lineWidth: action.payload
                    }
                }
            };
        }
        case SET_SCALE:
            const boundedScale = Math.min(Math.max(action.payload, minScale), maxScale);

            return {
                ...state,
                view: {
                    ...state.view,
                    scale: boundedScale
                }
            };
        case SET_OFFSET:
            return {
                ...state,
                view: {
                    ...state.view,
                    offsetX: action.payload.offsetX,
                    offsetY: action.payload.offsetY
                }
            };
        default:
            return state;
    }
};

export default drawingReducer;