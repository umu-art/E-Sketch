import { BoardApi, UserApi } from 'est_proxy_api';
import { FigureType, Line, Point } from 'figures/dist/index.js';
import { encode, encodePoint } from 'coder/dist/index.js';
import { WebSocket } from 'ws';

let userApi = new UserApi();
let boardApi = new BoardApi();

const MAX_BOARD_COUNT = 16;
const MAX_USER_COUNT = 16;

async function main() {
  const email = process.env.EST_TEST_EMAIL || 'est-test@mail.ru';
  const passwordHash = process.env.EST_TEST_PASSWORD_HASH;

  if (!passwordHash) {
    console.error('EST_TEST_PASSWORD_HASH environment variable is not set');
    process.exit(1);
  }

  global.authCookie = await userApi.loginWithHttpInfo({
    authDto: {
      email: email,
      passwordHash: passwordHash,
    },
  }).then(resp => {
    return resp.response.headers['set-cookie'][0]
      .split('; ')[0] + ';';
  });

  boardApi.apiClient.defaultHeaders = {
    ...boardApi.apiClient.defaultHeaders,
    Cookie: authCookie,
  };

  for (let boardCount = 1; boardCount <= MAX_BOARD_COUNT; boardCount *= 2) {
    for (let userCount = 1; userCount <= MAX_USER_COUNT; userCount *= 2) {
      console.log('Starting with ' + boardCount + ' boards and ' + userCount + ' users');
      await Promise.all(Array(boardCount).fill(null).map((_, i) => {
        return boardApi.create({
          createRequest: {
            name: 'Test board at ' + i + ' at ' + new Date().toISOString(),
            description: 'This is a test board',
            linkSharedMode: 'none_by_link',
          },
        }).then(board => imitateBoardFrFrFr(board, userCount));
      }));
    }
  }
}

async function imitateBoardFrFrFr(board, usersCount) {
  const startTime = new Date().getTime();

  await Promise.all(Array(usersCount).fill(null)
    .map((_) => imitateUserFrFrFr(board, authCookie)));

  const endTime = new Date().getTime();
  console.log('#!! Board ' + board.id + ' with ' + usersCount + ' users imitated in ' + (endTime - startTime) + ' ms !!#');
}

async function imitateUserFrFrFr(board) {
  let markerWebSocket = await connect('wss://e-sketch.ru/proxy/marker/ws?boardId=' + board.id);
  let figureWebSocket = await connect('wss://e-sketch.ru/proxy/figure/ws?boardId=' + board.id);

  for (let i = 0; i < 30; i++) {
    await createFigure(figureWebSocket, markerWebSocket);
  }

  figureWebSocket.close();
  markerWebSocket.close();
}

async function createFigure(figureWebSocket, markerWebSocket) {
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

    let figure = new Line(FigureType.LINE, figureId, ['blue', 3], []);

    await safeSend(figureWebSocket, String.fromCharCode(1) + encode(figure));
    await sleep(1000 / 60);

    for (let i = 0; i < 120; i++) {
      const oldFigureEncoded = encode(figure);
      figure.points.push(new Point(Math.random() * 1000, Math.random() * 1000));
      const newFigureEncoded = encode(figure);
      const newFigurePart = newFigureEncoded.slice(oldFigureEncoded.length);

      if (i % 2 === 0) {
        await Promise.all([
          safeSend(figureWebSocket, String.fromCharCode(3) + String.fromCharCode(figure.type) + figure.id + newFigurePart),
          safeSend(markerWebSocket, encodePoint(new Point(Math.random() * 1000, Math.random() * 1000))),
          sleep(1000 / 60),
        ]);
      } else {
        await Promise.all([
          safeSend(figureWebSocket, String.fromCharCode(3) + String.fromCharCode(figure.type) + figure.id + newFigurePart),
          sleep(1000 / 60),
        ]);
      }
    }
  } catch (e) {
    console.log(e);
  }
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function safeSend(webSocket, message) {
  if (webSocket.readyState !== WebSocket.OPEN) {
    console.log('reopening');
    webSocket = connect(webSocket.url);
  }
  webSocket.send(message);
}

async function connect(addr) {
  const maxRetries = 5;

  for (let i = 0; i < maxRetries; i++) {
    try {
      let webSocket = await tryConnect(addr);
      if (webSocket) {
        return webSocket;
      }
    } catch (error) {
      console.log(`Connection attempt ${i + 1} failed: ${error.message}`);
      if (i === maxRetries - 1) {
        throw new Error('Unable to connect to WebSocket server after multiple attempts');
      }
    }
  }
}

async function tryConnect(addr) {
  return new Promise((resolve, reject) => {
    let webSocket = new WebSocket(addr, {
      headers: {
        Cookie: authCookie,
      },
    });

    webSocket.addEventListener('open', () => {
      webSocket.addEventListener('close', () => {
      });
      resolve(webSocket);
    });

    webSocket.addEventListener('error', (error) => {
      if (error.error && error.error.code === 'ETIMEDOUT') {
        reject(new Error('Connection timed out'));
      } else {
        reject(error);
      }
    });
  });
}

main()
  .catch(e => console.log(e));