import { createMarker, onMarkerUpdate } from './SocketApi';
import { Point } from 'figures/dist';

import { drawing, BASE_OFFSET_X, BASE_OFFSET_Y } from './Paint';
import { UserApi } from 'est_proxy_api';

const apiInstance = new UserApi();

const timeoutDuration = 10000;

export function registerMarkersListener(board) {
  let markersMap = new Map();
  let timersMap = new Map();

  let me;

  apiInstance.getSelf().then((data) => {
    me = data;
  }).catch((error) => {
    console.log(error);
  });

  board.addEventListener('mousemove', (event) => {
    const rect = board.getBoundingClientRect();

    const x = (event.offsetX + BASE_OFFSET_X - rect.left) / drawing.scale - drawing.offsetX;
    const y = (event.offsetY + BASE_OFFSET_Y - rect.top) / drawing.scale - drawing.offsetY;

    const point = new Point(x, y);

    createMarker(point);
  });

  onMarkerUpdate((point, username) => {
    if (username === me.username) {
      return;
    }
    
    if (markersMap[username]) {
      document.body.removeChild(markersMap[username]);
    }

    const windowWidth = window.innerWidth;
    const windowHeight = window.innerHeight;

    const marker = document.createElement('p');
    marker.textContent = username;
    marker.style.pointerEvents = 'none';
    marker.style.position = 'absolute';
    marker.style.visibility = 'hidden';

    document.body.appendChild(marker);

    const markerWidth = marker.offsetWidth;
    const markerHeight = marker.offsetHeight;

    document.body.removeChild(marker);

    const x = (point.x + drawing.offsetX) * drawing.scale;
    const y = (point.y + drawing.offsetY) * drawing.scale;

    const safeX = Math.min(Math.max(x + 10, 10), windowWidth - (markerWidth + 10));
    const safeY = Math.min(Math.max(y - 10, 10), windowHeight - (markerHeight + 10));

    marker.style.left = `${safeX}px`;
    marker.style.top = `${safeY}px`;
    marker.style.backgroundColor = 'yellow';
    marker.style.userSelect = 'none';
    marker.style.visibility = 'visible';
    marker.style.zIndex = 10;

    document.body.appendChild(marker);

    markersMap[username] = marker;
    document.body.appendChild(marker);

    if (timersMap[username]) {
      clearTimeout(timersMap[username]);
    }

    timersMap[username] = setTimeout(() => {
      if (markersMap[username]) {
          document.body.removeChild(markersMap[username]);
          delete markersMap[username];
      }

      delete timersMap[username];
    }, timeoutDuration);
  });
}