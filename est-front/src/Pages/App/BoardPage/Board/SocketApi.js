import { decode, decodePoint, encode, encodePoint } from 'coder/dist';

const backUrl = 'wss://e-sketch.ru';

export let figureWebSocket;

export function connectToFigures(boardId) {
  figureWebSocket = new WebSocket(backUrl + '/proxy/figure/ws?boardId=' + boardId);

  figureWebSocket.onopen = () => {
    console.log('WebSocket connected to figures');
  };

  figureWebSocket.onclose = () => {
    console.log('WebSocket disconnected from figures');
    connectToFigures(boardId);
  };
}

export function onNewFigure(handler) {
  figureWebSocket.addEventListener('message', (event) => {
    if (event.data.length > 36 &&
      event.data[0] !== '-' &&
      event.data[0] !== '+') {

      handler(decode(event.data));
    }
  });
}

export function onRemoveFigure(handler) {
  figureWebSocket.addEventListener('message', (event) => {
    if (event.data.length > 36 &&
      event.data[0] === '-') {
      const id = event.data.slice(1, 38);

      handler(id);
    }
  });
}

export function onUpdateFigure(handler) {
  figureWebSocket.addEventListener('message', (event) => {
    if (event.data.length > 36 &&
      event.data[0] === '+') {

      const id = event.data.slice(2, 38);
      const newData = event.data.slice(38);
      handler(id, newData);
    }
  });
}

export function createFigure(onNewFigure) {
  figureWebSocket.send(String.fromCharCode(0));

  const messageHandler = (event) => {
    if (event.data.length === 36) {
      onNewFigure(event.data);
      figureWebSocket.removeEventListener('message', messageHandler);
    }
  };

  figureWebSocket.addEventListener('message', messageHandler);
}

export function changeFigure(figure) {
  figureWebSocket.send(String.fromCharCode(1) + encode(figure));
}

export function deleteFigure(figureId) {
  figureWebSocket.send(String.fromCharCode(2) + String.fromCharCode(0) + figureId);
}

export function updateFigure(newFigure, oldFigure) {
  const oldFigureEncoded = encode(oldFigure);
  const newFigureEncoded = encode(newFigure);
  const newFigurePart = newFigureEncoded.slice(oldFigureEncoded.length);
  figureWebSocket.send(String.fromCharCode(3) + String.fromCharCode(newFigure.type) + newFigure.id + newFigurePart);
}

export function getAllFigures() {
  figureWebSocket.send(String.fromCharCode(4));
}

export let markerWebSocket;

export function connectToMarkers(boardId) {
  markerWebSocket = new WebSocket(backUrl + '/proxy/marker/ws?boardId=' + boardId);

  markerWebSocket.addEventListener('open', () => {
    console.log('WebSocket connected to markers');
  });

  markerWebSocket.addEventListener('close', () => {
    console.log('WebSocket disconnected from markers');
    connectToMarkers(boardId);
  });
}

export function onMarkerUpdate(handler) {
  markerWebSocket.addEventListener('message', (event) => {
    const point = decodePoint(event.data.slice(0, 16));
    const username = event.data.slice(16);
    handler(point, username);
  });
}

export function createMarker(point) {
  if (markerWebSocket.readyState === WebSocket.OPEN) {
    markerWebSocket.send(encodePoint(point));
  }
}