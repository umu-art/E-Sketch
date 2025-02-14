import { decode, encode } from 'coder/dist';

const backUrl = 'wss://e-sketch.ru';

class FigureWebSocket {
    constructor(boardId) {
        this.boardId = boardId;
        this.webSocket = null;
        this.newFigureHandlers = [];
        this.removeFigureHandlers = [];
        this.updateFigureHandlers = [];
        this.messageQueue = [];
        this.isReconnecting = false;
        this.connect();
    }

    connect() {
        this.webSocket = new WebSocket(`${backUrl}/proxy/figure/ws?boardId=${this.boardId}`);

        this.webSocket.onopen = () => {
            console.log('WebSocket connected to figures');
            this.isReconnecting = false;
            this.flushMessageQueue();
            this.getAllFigures();
        };

        this.webSocket.onclose = () => {
            console.log('WebSocket disconnected from figures');
            this.isReconnecting = true;
            this.reconnect();
        };

        this.webSocket.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.webSocket.close();
        };

        this.webSocket.addEventListener('message', this.handleMessage.bind(this));
    }

    reconnect() {
        setTimeout(() => {
            if (this.isReconnecting) {
                this.connect();
            }
        }, 500);
    }

    flushMessageQueue() {
        while (this.messageQueue.length > 0) {
            const message = this.messageQueue.shift();
            this.webSocket.send(message);
        }
    }

    handleMessage(event) {
        if (event.data.length > 36) {
            const actionType = event.data[0];

            if (actionType === '-') {
                const id = event.data.slice(1, 38);
                this.removeFigureHandlers.forEach(handler => handler(id));
            } else if (actionType === '+') {
                const id = event.data.slice(2, 38);
                const newData = event.data.slice(38);
                this.updateFigureHandlers.forEach(handler => handler(id, newData));
            } else {
                this.newFigureHandlers.forEach(handler => handler(decode(event.data)));
            }
        }
    }

    onNewFigure(handler) {
        this.newFigureHandlers.push(handler);
    }

    onRemoveFigure(handler) {
        this.removeFigureHandlers.push(handler);
    }

    onUpdateFigure(handler) {
        this.updateFigureHandlers.push(handler);
    }

    sendMessage(message) {
        if (this.webSocket && this.webSocket.readyState === WebSocket.OPEN) {
            this.webSocket.send(message)
        } else {
            this.messageQueue.push(message);
        }
    }

    createFigure(onNewFigure) {
        this.sendMessage(String.fromCharCode(0));

        const messageHandler = (event) => {
            if (event.data.length === 36) {
                onNewFigure(event.data);
                this.webSocket.removeEventListener('message', messageHandler);
            }
        };

        this.webSocket.addEventListener('message', messageHandler);
    }

    changeFigure(figure) {
        this.sendMessage(String.fromCharCode(1) + encode(figure));
    }

    deleteFigure(figureId) {
        this.sendMessage(String.fromCharCode(2) + String.fromCharCode(0) + figureId);
    }

    updateFigure(newFigure, oldFigure) {
        const oldFigureEncoded = encode(oldFigure);
        const newFigureEncoded = encode(newFigure);
        const newFigurePart = newFigureEncoded.slice(oldFigureEncoded.length);
        this.sendMessage(String.fromCharCode(3) + String.fromCharCode(newFigure.type) + newFigure.id + newFigurePart);
    }

    getAllFigures() {
        this.sendMessage(String.fromCharCode(4));
    }
}

export default FigureWebSocket;