import { createMarker, onMarkerUpdate } from './SocketApi';
import { Point } from 'figures/dist';

export function registerMarkersListener(board: Element) {
  let markersMap = new Map();

  board.addEventListener('mousemove', (event) => {
    const rect = board.getBoundingClientRect();
    const x = event.offsetX - rect.left;
    const y = event.offsetY - rect.top;
    const point = new Point(x, y);
    createMarker(point);
  });

  onMarkerUpdate((point, username) => {
    if (markersMap['username']) {
      document.body.removeChild(markersMap['username']);
    }

    markersMap['username'] = document.createElement('p');
    markersMap['username'].style.position = 'absolute';
    markersMap['username'].style.left = `${point.x + 10}px`;
    markersMap['username'].style.top = `${point.y + 10}px`;
    markersMap['username'].style.fontSize = '12px';
    markersMap['username'].style.color = 'black';
    markersMap['username'].textContent = username;
    markersMap['username'].style.backgroundColor = 'yellow';
    markersMap['username'].style.zIndex = 1000;
    document.body.appendChild(markersMap['username']);
  });
}