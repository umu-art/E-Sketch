import drawingReducer from './reducer';
import {
    SET_TOOL,
    SET_LINE_COLOR,
    SET_LINE_WIDTH,
    SET_SCALE,
    SET_OFFSET,
    SET_GPT_STATUS,
    ADD_MESSAGE,
    REMOVE_MESSAGE,
    SHOW_GPT_POPOVER,
    HIDE_GPT_POPOVER,
} from './actions';

describe('drawingReducer', () => {
    const initialState = {
        state: 'idle',
        tool: 'pencil',
        tools: {
        pencil: {
            lineColor: '#1677ff',
            lineWidth: 2,
        },
        eraser: {
            lineWidth: 5,
        },
        rectangle: {
            lineColor: '#1677ff',
            fillColor: '#ffffff00',
            lineWidth: 2,
        },
        ellipse: {
            lineColor: '#1677ff',
            fillColor: '#ffffff00',
            lineWidth: 2,
        },
        gpt: {
            status: null,
        },
        },
        view: {
        scale: 1,
        offsetX: 0,
        offsetY: 0,
        startX: 0,
        startY: 0,
        },
        messages: [],
        popover: {
        visible: false,
        request: null,
        },
    };

    it('should return the initial state', () => {
        expect(drawingReducer(undefined, {})).toEqual(initialState);
    });

    it('should handle SET_TOOL', () => {
        const action = { type: SET_TOOL, payload: 'eraser' };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        tool: 'eraser',
        });
    });

    it('should handle SET_LINE_COLOR', () => {
        const action = { type: SET_LINE_COLOR, payload: '#ff0000' };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        tools: {
            ...initialState.tools,
            pencil: {
            ...initialState.tools.pencil,
            lineColor: '#ff0000',
            },
        },
        });
    });

    it('should handle SET_LINE_WIDTH', () => {
        const action = { type: SET_LINE_WIDTH, payload: 5 };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        tools: {
            ...initialState.tools,
            pencil: {
            ...initialState.tools.pencil,
            lineWidth: 5,
            },
        },
        });
    });

    it('should handle SET_SCALE', () => {
        const action = { type: SET_SCALE, payload: 2 };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        view: {
            ...initialState.view,
            scale: 2,
        },
        });
    });

    it('should handle SET_OFFSET', () => {
        const action = { type: SET_OFFSET, payload: { offsetX: 10, offsetY: 20 } };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        view: {
            ...initialState.view,
            offsetX: 10,
            offsetY: 20,
        },
        });
    });

    it('should handle SET_GPT_STATUS', () => {
        const action = { type: SET_GPT_STATUS, payload: 'processing' };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        tools: {
            ...initialState.tools,
            gpt: {
            status: 'processing',
            },
        },
        });
    });

    it('should handle ADD_MESSAGE', () => {
        const message = { id: 1, title: 'Test', content: 'Test message' };
        const action = { type: ADD_MESSAGE, payload: message };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        messages: [message],
        });
    });

    it('should handle REMOVE_MESSAGE', () => {
        const initialStateWithMessages = {
        ...initialState,
        messages: [{ id: 1, title: 'Test', content: 'Test message' }],
        };
        const action = { type: REMOVE_MESSAGE, payload: 1 };
        expect(drawingReducer(initialStateWithMessages, action)).toEqual({
        ...initialStateWithMessages,
        messages: [],
        });
    });

    it('should handle SHOW_GPT_POPOVER', () => {
        const request = { leftUp: { x: 0, y: 0 }, rightDown: { x: 100, y: 100 } };
        const action = { type: SHOW_GPT_POPOVER, payload: request };
        expect(drawingReducer(initialState, action)).toEqual({
        ...initialState,
        popover: {
            visible: true,
            request,
        },
        });
    });

    it('should handle HIDE_GPT_POPOVER', () => {
        const initialStateWithPopover = {
        ...initialState,
        popover: {
            visible: true,
            request: { leftUp: { x: 0, y: 0 }, rightDown: { x: 100, y: 100 } },
        },
        };
        const action = { type: HIDE_GPT_POPOVER };
        expect(drawingReducer(initialStateWithPopover, action)).toEqual({
        ...initialStateWithPopover,
        popover: {
            visible: false,
            request: null,
        },
        });
    });
});