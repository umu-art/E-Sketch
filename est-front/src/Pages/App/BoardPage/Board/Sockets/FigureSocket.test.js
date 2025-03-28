import FigureWebSocket from './FigureSocket';


global.WebSocket = jest.fn(() => ({
    onopen: jest.fn(),
    onclose: jest.fn(),
    onerror: jest.fn(),
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    send: jest.fn(),
    close: jest.fn(),
    readyState: 1,
}));

jest.mock('coder/dist', () => ({
    encode: jest.fn((data) => `encoded-${JSON.stringify(data)}`),
    decode: jest.fn((data) => JSON.parse(data.replace('encoded-', ''))),
}));

describe('FigureWebSocket', () => {
    let wsInstance;
    const mockBoardId = 'test-board-id';
    const mockWebSocket = {
        onopen: jest.fn(),
        onclose: jest.fn(),
        onerror: jest.fn(),
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
        send: jest.fn(),
        close: jest.fn(),
        readyState: 1,
    };

    beforeEach(() => {
        jest.clearAllMocks();
        WebSocket.mockImplementation(() => mockWebSocket);
        wsInstance = new FigureWebSocket(mockBoardId);
    });

    it('should initialize with correct properties', () => {
        expect(wsInstance.boardId).toBe(mockBoardId);
        expect(wsInstance.newFigureHandlers).toEqual([]);
        expect(wsInstance.removeFigureHandlers).toEqual([]);
        expect(wsInstance.updateFigureHandlers).toEqual([]);
        expect(wsInstance.messageQueue).toEqual([]);
        expect(wsInstance.isReconnecting).toBe(false);
    });

    it('should connect to correct WebSocket URL', () => {
        expect(WebSocket).toHaveBeenCalledWith(
            'wss://e-sketch.ru/proxy/figure/ws?boardId=test-board-id'
        );
    });

    it('should set up WebSocket event handlers', () => {
        expect(mockWebSocket.onopen).toBeDefined();
        expect(mockWebSocket.onclose).toBeDefined();
        expect(mockWebSocket.onerror).toBeDefined();
        expect(mockWebSocket.addEventListener).toHaveBeenCalledWith(
            'message',
            expect.any(Function)
        );
    });

    describe('onopen', () => {
        it('should flush message queue and get all figures when connection opens', () => {
            wsInstance.messageQueue = ['test-message'];
            wsInstance.getAllFigures = jest.fn();
            
            mockWebSocket.onopen();
            
            expect(mockWebSocket.send).toHaveBeenCalledWith('test-message');
            expect(wsInstance.messageQueue).toEqual([]);
            expect(wsInstance.getAllFigures).toHaveBeenCalled();
            expect(wsInstance.isReconnecting).toBe(false);
        });
    });

    describe('onclose', () => {
        it('should set isReconnecting and attempt to reconnect', () => {
            jest.useFakeTimers();
            mockWebSocket.onclose();
            
            expect(wsInstance.isReconnecting).toBe(true);
            jest.advanceTimersByTime(500);
            expect(WebSocket).toHaveBeenCalledTimes(2);
            
            jest.useRealTimers();
        });
    });

    describe('onerror', () => {
        it('should close WebSocket on error', () => {
            const mockError = new Error('Test error');
            mockWebSocket.onerror(mockError);
            
            expect(mockWebSocket.close).toHaveBeenCalled();
        });
    });

    describe('handleMessage', () => {
        const mockHandler = jest.fn();
        
        it('should call removeFigureHandlers for "-" action', () => {
            wsInstance.removeFigureHandlers = [mockHandler];
            const testId = '123456789012345678901234567890123456';
            
            wsInstance.handleMessage({ data: `-${testId}` });
            
            expect(mockHandler).toHaveBeenCalledWith(testId);
        });
        
        it('should call updateFigureHandlers for "+" action', () => {
            wsInstance.updateFigureHandlers = [mockHandler];
            const testId = '123456789012345678901234567890123456';
            const testData = 'test-data';
            
            wsInstance.handleMessage({ data: `+0${testId}${testData}` });
            
            expect(mockHandler).toHaveBeenCalledWith(testId, testData);
        });
    });

    describe('event handlers', () => {
        it('should add new figure handler', () => {
            const handler = jest.fn();
            wsInstance.onNewFigure(handler);
            expect(wsInstance.newFigureHandlers).toContain(handler);
        });
        
        it('should add remove figure handler', () => {
            const handler = jest.fn();
            wsInstance.onRemoveFigure(handler);
            expect(wsInstance.removeFigureHandlers).toContain(handler);
        });
        
        it('should add update figure handler', () => {
            const handler = jest.fn();
            wsInstance.onUpdateFigure(handler);
            expect(wsInstance.updateFigureHandlers).toContain(handler);
        });
    });

    describe('sendMessage', () => {
        it('should send message immediately if WebSocket is open', () => {
            mockWebSocket.readyState = WebSocket.OPEN;
            wsInstance.sendMessage('test-message');
            expect(mockWebSocket.send).toHaveBeenCalledWith('test-message');
        });
        
        it('should add message to queue if WebSocket is not open', () => {
            mockWebSocket.readyState = 0; // CONNECTING
            wsInstance.sendMessage('test-message');
            expect(mockWebSocket.send).not.toHaveBeenCalled();
            expect(wsInstance.messageQueue).toContain('test-message');
        });
    });

    describe('createFigure', () => {
        it('should send create message and set up temporary handler', () => {
            const callback = jest.fn();
            mockWebSocket.readyState = WebSocket.OPEN;
            
            wsInstance.createFigure(callback);
            
            expect(mockWebSocket.send).toHaveBeenCalledWith('\x00');
            expect(mockWebSocket.addEventListener).toHaveBeenCalledWith(
                'message',
                expect.any(Function)
            );
        });
    });

    describe('changeFigure', () => {
        it('should send encoded figure with type 1', () => {
            mockWebSocket.readyState = WebSocket.OPEN;
            const testFigure = { id: 'test-id', type: 1 };
            
            require('coder/dist').encode.mockReturnValueOnce('encoded-figure-data');
            
            wsInstance.changeFigure(testFigure);
            
            expect(mockWebSocket.send).toHaveBeenCalledWith(
                '\x01encoded-figure-data'
            );
        });
    });

    describe('deleteFigure', () => {
        it('should send delete message with figure ID', () => {
            const testId = 'test-figure-id';
            wsInstance.deleteFigure(testId);
            
            expect(mockWebSocket.send).toHaveBeenCalledWith(
                String.fromCharCode(2) + String.fromCharCode(0) + testId
            );
        });
    });

    describe('updateFigure', () => {
        it('should send update message with diff', () => {
            mockWebSocket.readyState = WebSocket.OPEN;
            const oldFigure = { id: 'test-id', type: 1, prop: 'old' };
            const newFigure = { id: 'test-id', type: 1, prop: 'new' };
            
            require('coder/dist').encode
                .mockReturnValueOnce('encoded-old-figure')
                .mockReturnValueOnce('encoded-new-figure');
            
            wsInstance.updateFigure(newFigure, oldFigure);
            
            expect(mockWebSocket.send).toHaveBeenCalledWith(
                expect.stringContaining("test-id")
            );
        });
    });

    describe('getAllFigures', () => {
        it('should send message with type 4', () => {
            wsInstance.getAllFigures();
            expect(mockWebSocket.send).toHaveBeenCalledWith(String.fromCharCode(4));
        });
    });
});