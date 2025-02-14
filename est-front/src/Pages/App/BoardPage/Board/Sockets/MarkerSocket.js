import { decodePoint, encodePoint } from 'coder/dist';

const backUrl = 'wss://e-sketch.ru';

class MarkerWebSocket {
    constructor(boardId) {
        this.boardId = boardId;
        this.webSocket = null;
        this.markerUpdateHandlers = [];
        this.connect();
    }

    connect() {
        this.webSocket = new WebSocket(`${backUrl}/proxy/marker/ws?boardId=${this.boardId}`);

        this.webSocket.addEventListener('open', () => {
            console.log('WebSocket connected to markers');
        });

        this.webSocket.addEventListener('close', () => {
            console.log('WebSocket disconnected from markers');
            setTimeout(() => this.connect(), 1000);
        });

        this.webSocket.addEventListener('message', this.handleMessage.bind(this));
    }

    handleMessage(event) {
        const point = decodePoint(event.data.slice(0, 16));
        const username = event.data.slice(16);
        this.markerUpdateHandlers.forEach(handler => handler(point, username));
    }

    onMarkerUpdate(handler) {
        this.markerUpdateHandlers.push(handler);
    }

    createMarker(point) {
        if (this.webSocket.readyState === WebSocket.OPEN) {
            this.webSocket.send(encodePoint(point));
        }
    }
}

export default MarkerWebSocket;