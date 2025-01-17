import { BoardApi, UserApi } from 'est_proxy_api';
import { FigureType, Line, Point } from 'figures/dist/index.js';
import { encode } from 'coder/dist/index.js';

let userApi = new UserApi();
let boardApi = new BoardApi();

const BOARD_COUNT = 16;
const USER_COUNT = 16;
const TEST_TIMEOUT = 5 * 60 * 1000;

async function main() {
  let authCookie = await userApi.loginWithHttpInfo({
    authDto: {
      email: 'est-test@mail.ru',
      passwordHash: '123',
    },
  }).then(resp => {
    return resp.response.headers['set-cookie'][0]
      .split('; ')[0] + ';';
  });

  boardApi.apiClient.defaultHeaders = {
    ...boardApi.apiClient.defaultHeaders,
    Cookie: authCookie,
  };

  return Promise.all(Array(BOARD_COUNT).fill(null).map((_, i) => {
    return boardApi.create({
      createRequest: {
        name: 'Test board at ' + i + ' at ' + new Date().toISOString(),
        description: 'This is a test board',
        linkSharedMode: 'none_by_link',
      },
    }).then(board => imitateBoardFrFrFr(board, USER_COUNT, authCookie));
  }));
}

async function imitateBoardFrFrFr(board, usersCount, authCookie) {
  return Promise.all(Array(usersCount).fill(null)
    .map((_) => imitateUserFrFrFr(board, authCookie)));
}

async function imitateUserFrFrFr(board, authCookie) {
  let startTime = new Date().getTime();
  let markerWebSocket = new WebSocket('wss://e-sketch.ru/proxy/marker/ws?boardId=' + board.id, {
    headers: {
      Cookie: authCookie,
    },
  });

  while (new Date().getTime() - startTime < TEST_TIMEOUT) {
    let figureWebSocket = new WebSocket('wss://e-sketch.ru/proxy/figure/ws?boardId=' + board.id, {
      headers: {
        Cookie: authCookie,
      },
    });

    await new Promise((resolve, reject) => {
      figureWebSocket.onopen = resolve;
      figureWebSocket.onerror = reject;
    });

    await createFigure(figureWebSocket);
    figureWebSocket.close();
  }

  markerWebSocket.close();
}

async function createFigure(figureWebSocket) {
  try {
    let figureId = await new Promise((resolve, _) => {
      const messageHandler = (event) => {
        if (event.data.length === 36) {
          figureWebSocket.removeEventListener('message', messageHandler);
          resolve(event.data);
        }
      };

      figureWebSocket.addEventListener('message', messageHandler);
      figureWebSocket.send(String.fromCharCode(0));
    });

    let figure = new Line(FigureType.LINE, figureId, ['red', '3'], []);
    figureWebSocket.send(String.fromCharCode(1) + encode(figure));
    await sleep(1000 / 60);

    for (let i = 0; i < 1000; i++) {
      const oldFigureEncoded = encode(figure);
      figure.points.push(new Point(Math.random() * 1000, Math.random() * 1000));
      const newFigureEncoded = encode(figure);
      const newFigurePart = newFigureEncoded.slice(oldFigureEncoded.length);

      figureWebSocket.send(String.fromCharCode(3) + String.fromCharCode(figure.type) + figure.id + newFigurePart);
      await sleep(1000 / 60);
    }

    console.log('Created figure with ID:', figureId);
  }catch (e) {
    console.log(e);
  }
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

main()
  .catch(console.error);