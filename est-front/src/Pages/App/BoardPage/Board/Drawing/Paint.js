import store from '../../../../../redux/store';
import { DrawingManager } from './DrawingManager';

import * as d3 from 'd3';
import { GPTManager } from '../GPT/GPTManager';

export function registerDrawListener(board, boardController, initialDrawing, figureWebSocket) {
    let settings = initialDrawing;
    let drawing = { nowX: 0, nowY: 0 };

    const drawingManager = new DrawingManager(boardController, figureWebSocket, settings, drawing);
    const gptManager = new GPTManager(drawingManager);

    console.log(gptManager);

    store.subscribe(() => {
        const newState = store.getState();
        settings = newState;
    });

    d3.select('.board')
        .on('mousedown', function(event) {
            if (event.button === 0) {
                drawingManager.isMouseDown = true;
            }
        })
        .on('mouseup', function() {
            drawingManager.isMouseDown = false;
        })
        .on('mousemove', function(event) {
            const target = event.target;

            if (target.tagName !== 'path' && target.tagName !== 'ellipse' && target.tagName !== 'rect')
                return;

            if (settings.tool !== 'eraser' || !drawingManager.isMouseDown)
                return;

            const pathId = d3.select(target).attr('id');

            const figure = boardController.figures.find(figure => figure.id === pathId);
            figure.id = "waiting"

            boardController.removeFigure(pathId);
            figureWebSocket.deleteFigure(pathId);

            drawingManager.history.push(() => {
                figureWebSocket.createFigure((uuid) => {
                    figure.id = uuid;
                
                    boardController.upsertFigure(figure);
                    figureWebSocket.changeFigure(figure);
                });
            });
        });

    board.addEventListener('contextmenu', (e) => e.preventDefault());
}