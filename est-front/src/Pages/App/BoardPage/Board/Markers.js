import { Point } from 'figures/dist';

import { BASE_OFFSET_X, BASE_OFFSET_Y } from './Drawing/Constants'
import { UserApi } from 'est_proxy_api';
import store from '../../../../redux/store';

const apiInstance = new UserApi();

const timeoutDuration = 10000;

const colorPalette = [
  "#FF5733", // Красный
  "#33FF57", // Зеленый
  "#3357FF", // Синий
  "#F3FF33", // Желтый
  "#FF33A1", // Розовый
  "#33FFF7", // Бирюзовый
  "#B833FF", // Фиолетовый
  "#FFC300", // Золотой
  "#DAF7A6", // Лаймовый
  "#FF3377"  // Ярко-розовый
];

const colorCache = {};

function generateColorFromUsername(username) {
  if (colorCache[username]) {
    return colorCache[username];
  }

  let hash = 0;
  for (let i = 0; i < username.length; i++) {
      hash = username.charCodeAt(i) + ((hash << 5) - hash);
  }

  const colorIndex = Math.abs(hash) % colorPalette.length;

  colorCache[username] = colorPalette[colorIndex];
  
  return colorPalette[colorIndex];
}

export function registerMarkersListener(board, initialDrawing, markerWebSocket) {
  let markersMap = new Map();
  let timersMap = new Map();

  let me;

  let drawing = initialDrawing;

  store.subscribe(() => {
    const newState = store.getState()

    drawing = newState;
  })

  apiInstance.getSelf().then((data) => {
    me = data;
  }).catch((error) => {
    console.log(error);
  });

  board.addEventListener('mousemove', (event) => {
    const rect = board.getBoundingClientRect();

    const x = (event.offsetX + BASE_OFFSET_X - rect.left) / drawing.view.scale - drawing.view.offsetX;
    const y = (event.offsetY + BASE_OFFSET_Y - rect.top) / drawing.view.scale - drawing.view.offsetY;

    const point = new Point(x, y);

    markerWebSocket.createMarker(point);
  });

  markerWebSocket.onMarkerUpdate((point, username) => {
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

    marker.style.padding = '5px 10px';
    marker.style.borderRadius = '10px 10px 10px 0px';
    marker.style.boxShadow = '2px 2px 5px rgba(0, 0, 0, 0.3)';

    document.body.appendChild(marker);

    const markerWidth = marker.offsetWidth;
    const markerHeight = marker.offsetHeight;

    document.body.removeChild(marker);

    const x = (point.x + drawing.view.offsetX) * drawing.view.scale;
    const y = (point.y + drawing.view.offsetY) * drawing.view.scale;

    const safeX = Math.min(Math.max(x, 0), windowWidth - markerWidth);
    const safeY = Math.min(Math.max(y - markerHeight, 0), windowHeight - markerHeight - 1);

    marker.style.left = `${safeX}px`;
    marker.style.top = `${safeY}px`;
    marker.style.backgroundColor = generateColorFromUsername(username);
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