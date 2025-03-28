import React from 'react';
import { render, screen } from '@testing-library/react';
import { Provider } from 'react-redux';
import configureStore from 'redux-mock-store';
import Board from './Board';

jest.mock('paint/dist', () => ({
    Board: jest.fn()
}));

jest.mock('./Drawing/Paint', () => ({
    registerDrawListener: jest.fn()
}));

jest.mock('./Markers', () => ({
    registerMarkersListener: jest.fn()
}));

jest.mock('./Sockets/FigureSocket', () => jest.fn());
jest.mock('./Sockets/MarkerSocket', () => jest.fn());

const PaintBoard = require('paint/dist').Board;
const FigureWebSocket = require('./Sockets/FigureSocket');
const MarkerWebSocket = require('./Sockets/MarkerSocket');

describe('Board Component', () => {
    const mockStore = configureStore([]);
    let store;
    let mockFigureWSInstance;
    let mockMarkerWSInstance;
    let mockBoardInstance;

    beforeEach(() => {
        mockFigureWSInstance = {
            webSocket: {
                addEventListener: jest.fn((event, callback) => {
                    if (event === 'open') callback();
                })
            }
        };
        
        mockMarkerWSInstance = {};
        mockBoardInstance = {};
        
        PaintBoard.mockImplementation(() => mockBoardInstance);
        FigureWebSocket.mockImplementation(() => mockFigureWSInstance);
        MarkerWebSocket.mockImplementation(() => mockMarkerWSInstance);
        
        store = mockStore({
            tool: 'pencil',
        });
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('should initialize BoardController and WebSockets on mount', () => {
        render(
            <Provider store={store}>
                <Board boardId="test-board" />
            </Provider>
        );
        
        expect(PaintBoard).toHaveBeenCalledTimes(1);
        expect(FigureWebSocket).toHaveBeenCalledWith('test-board');
        expect(MarkerWebSocket).toHaveBeenCalledWith('test-board');
        
        expect(mockFigureWSInstance.webSocket.addEventListener)
        .toHaveBeenCalledWith('open', expect.any(Function));
    });

    it('should register listeners after initialization', () => {
        const { container } = render(
            <Provider store={store}>
                <Board boardId="test-board" />
            </Provider>
        );

        const boardElement = container.querySelector('#test-board');
        
        const { registerDrawListener } = require('./Drawing/Paint');
        expect(registerDrawListener).toHaveBeenCalledTimes(1);
        expect(registerDrawListener).toHaveBeenCalledWith(
            boardElement,
            mockBoardInstance,
            expect.objectContaining({ tool: 'pencil' }),
            mockFigureWSInstance
        );
        
        const { registerMarkersListener } = require('./Markers');
        expect(registerMarkersListener).toHaveBeenCalledTimes(1);
        expect(registerMarkersListener).toHaveBeenCalledWith(
            boardElement,
            expect.objectContaining({ tool: 'pencil' }),
            mockMarkerWSInstance
        );
    });

    it('should update cursor when tool changes', () => {
        const { rerender } = render(
            <Provider store={store}>
                <Board boardId="test-board" />
            </Provider>
        );
        
        expect(screen.getByTestId('board-svg'))
        .toHaveStyle('cursor: url(pencil.svg) 0 20, auto');
        
        store = mockStore({ tool: 'eraser' });
        rerender(
            <Provider store={store}>
                <Board boardId="test-board" />
            </Provider>
        );
        
        expect(screen.getByTestId('board-svg'))
        .toHaveStyle('cursor: url(eraser.svg) 0 20, auto');
    });

    it('should not reinitialize on rerender with same props', () => {
        const { rerender } = render(
            <Provider store={store}>
                <Board boardId="test-board" />
            </Provider>
        );
        
        const initialPaintBoardCalls = PaintBoard.mock.calls.length;
        
        rerender(
            <Provider store={store}>
                <Board boardId="test-board" />
            </Provider>
        );
        
        expect(PaintBoard.mock.calls.length).toBe(initialPaintBoardCalls);
    });
});