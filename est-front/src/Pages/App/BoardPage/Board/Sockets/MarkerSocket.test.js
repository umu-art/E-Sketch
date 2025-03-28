// MarkerSocket.test.js
import MarkerWebSocket from './MarkerSocket';

global.WebSocket = jest.fn(() => ({
    addEventListener: jest.fn(),
    readyState: 1,
    send: jest.fn(),
}));

jest.mock('coder/dist', () => ({
    decodePoint: jest.fn((data) => ({ x: data.charCodeAt(0), y: data.charCodeAt(1) })),
    encodePoint: jest.fn((point) => `${String.fromCharCode(point.x)}${String.fromCharCode(point.y)}`),
}));

describe('MarkerWebSocket', () => {
    let wsInstance;
    const mockBoardId = 'test-board-id';
    const mockWebSocket = {
        addEventListener: jest.fn(),
        readyState: 1,
        send: jest.fn(),
    };

    beforeEach(() => {
        jest.clearAllMocks();
        WebSocket.mockImplementation(() => mockWebSocket);
        wsInstance = new MarkerWebSocket(mockBoardId);
    });

    it('should initialize with correct properties', () => {
        expect(wsInstance.boardId).toBe(mockBoardId);
        expect(wsInstance.markerUpdateHandlers).toEqual([]);
    });

    it('should connect to correct WebSocket URL', () => {
        expect(WebSocket).toHaveBeenCalledWith(
            'wss://e-sketch.ru/proxy/marker/ws?boardId=test-board-id'
        );
    });

    it('should set up WebSocket event handlers', () => {
        expect(mockWebSocket.addEventListener).toHaveBeenCalledWith(
            'open',
            expect.any(Function)
        );
        expect(mockWebSocket.addEventListener).toHaveBeenCalledWith(
            'close',
            expect.any(Function)
        );
        expect(mockWebSocket.addEventListener).toHaveBeenCalledWith(
            'message',
            expect.any(Function)
        );
    });

    describe('on open event', () => {
        it('should log connection message', () => {
            console.log = jest.fn();
            const openHandler = mockWebSocket.addEventListener.mock.calls.find(
                call => call[0] === 'open'
            )[1];
            
            openHandler();
            
            expect(console.log).toHaveBeenCalledWith('WebSocket connected to markers');
        });
    });

    describe('on close event', () => {
        it('should attempt to reconnect after delay', () => {
            jest.useFakeTimers();
            const closeHandler = mockWebSocket.addEventListener.mock.calls.find(
                call => call[0] === 'close'
            )[1];
            
            closeHandler();
            
            expect(console.log).toHaveBeenCalledWith('WebSocket disconnected from markers');
            jest.advanceTimersByTime(1000);
            expect(WebSocket).toHaveBeenCalledTimes(2);
            
            jest.useRealTimers();
        });
    });

    describe('handleMessage', () => {
        const mockHandler = jest.fn();
        
        it('should call markerUpdateHandlers with decoded point and username', () => {
            wsInstance.markerUpdateHandlers = [mockHandler];
            const testPoint = { x: 10, y: 20 };
            const testUsername = 'test-user';
            
            // Создаем mock данных: 2 байта для точки + username в виде строки
            const pointData = `${String.fromCharCode(testPoint.x)}${String.fromCharCode(testPoint.y)}`;
            const mockData = `${pointData}00000000000000${testUsername}`;
            
            require('coder/dist').decodePoint.mockReturnValue(testPoint);
            
            wsInstance.handleMessage({ data: mockData });
            
            expect(mockHandler).toHaveBeenCalledWith(testPoint, testUsername);
        });
    });

    describe('onMarkerUpdate', () => {
        it('should add handler to markerUpdateHandlers', () => {
            const handler = jest.fn();
            wsInstance.onMarkerUpdate(handler);
            expect(wsInstance.markerUpdateHandlers).toContain(handler);
        });
    });

    describe('createMarker', () => {
        it('should send encoded point if WebSocket is open', () => {
            mockWebSocket.readyState = WebSocket.OPEN;
            const testPoint = { x: 10, y: 20 };
            
            require('coder/dist').encodePoint.mockReturnValueOnce('encoded-point');
            
            wsInstance.createMarker(testPoint);
            
            expect(mockWebSocket.send).toHaveBeenCalledWith('encoded-point');
        })
    });
});