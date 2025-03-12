import { decode, encode } from 'coder/dist';
import { BASE_OFFSET_X, BASE_OFFSET_Y, DrawingStates, FPS } from './Constants';
import { toolToClass, toolToFigureType } from './ToolMapping';
import { ViewManager } from './ViewManager';
import { Point } from 'figures/dist/point';
import store from '../../../../../redux/store';

export class DrawingManager {
    constructor(boardController, figureWebSocket, settings, drawing) {
        this.boardController = boardController;
        this.board = boardController.svgElement;

        this.viewManager = new ViewManager(this.board, settings);

        this.figureWebSocket = figureWebSocket;
        
        this.settings = settings;
        this.drawing = drawing;
        
        this.currentFigure = null;
        this.oldCurrentFigure = null;
        
        this.history = [];

        store.subscribe(() => {
            const newState = store.getState();
            this.settings = newState;
        });

        this.addEventListeners();
    }

    addEventListeners() {
        this.board.addEventListener('mousedown', (e) => this.handleMouseDown(e));
        this.board.addEventListener('mouseup', (e) => this.handleMouseUp(e));
        this.board.addEventListener('mousemove', (e) => this.handleMouseMove(e));

        this.figureWebSocket.onNewFigure((figure) => this.handleNewFigure(figure));
        this.figureWebSocket.onUpdateFigure((id, data) => this.handleUpdateFigure(id, data));
        this.figureWebSocket.onRemoveFigure((id) => this.handleRemoveFigure(id));

        window.addEventListener('keydown', (e) => this.keyPressHandler(e));
        setInterval(() => this.updateDrawing(), 1000 / FPS);
    }

    handleMouseMove(event) {
        const rect = this.board.getBoundingClientRect();
        
        this.drawing.nowX = (event.offsetX + BASE_OFFSET_X - rect.left) / this.settings.view.scale - this.settings.view.offsetX;
        this.drawing.nowY = (event.offsetY + BASE_OFFSET_Y - rect.top) / this.settings.view.scale - this.settings.view.offsetY;
    }

    handleMouseDown(e) {
        e.preventDefault();

        if (e.button !== 0 || this.settings.tool === 'eraser' || this.settings.tool === "gpt" || this.settings.state !== DrawingStates.IDLE) {
            return;
        }

        this.settings.state = DrawingStates.CREATING;

        this.currentFigure = toolToClass[this.settings.tool].startProcess(
            toolToFigureType[this.settings.tool],
            'waiting',
            Object.values(this.settings.tools[this.settings.tool]),
            new Point(this.drawing.nowX, this.drawing.nowY)
        );
        this.oldCurrentFigure = null;

        this.figureWebSocket.createFigure((uuid) => {
            this.boardController.removeFigure("waiting");

            this.currentFigure.id = uuid;

            if (this.settings.state === DrawingStates.CREATING) {
                this.settings.state = DrawingStates.DRAWING;
            } else {
                this.settings.state = DrawingStates.FINISHING;
            }

            this.history.push(() => {
                this.boardController.removeFigure(uuid);
                this.figureWebSocket.deleteFigure(uuid);
            });
        });
    }

    finishDrawing() {
        this.boardController.upsertFigure(this.currentFigure);

        this.triggerUpdateFigure(this.currentFigure, this.oldCurrentFigure);
        this.oldCurrentFigure = this.currentFigure.clone();

        this.settings.state = DrawingStates.IDLE;
    }

    waitForDrawingState() {
        if (this.settings.state === DrawingStates.DRAWING || this.settings.state === DrawingStates.FINISHING) {
            this.finishDrawing();
        } else {
            this.settings.state = DrawingStates.WAITING;
            setTimeout(() => this.waitForDrawingState(), 50);
        }
    }

    handleMouseUp(e) {
        e.preventDefault();

        if (e.button !== 0 || this.settings.tool === 'eraser') {
            return;
        }

        switch (this.settings.state) {
            case DrawingStates.CREATING:
                this.waitForDrawingState();
                break;
            case DrawingStates.DRAWING:
                this.finishDrawing();
                break;
            default:
                break;
        }
    }

    updateDrawing() {
        if (this.settings.state !== DrawingStates.DRAWING && this.settings.state !== DrawingStates.CREATING) {
            return;
        }

        this.currentFigure.process(new Point(this.drawing.nowX, this.drawing.nowY));
        this.boardController.upsertFigure(this.currentFigure);

        if (this.settings.state === DrawingStates.DRAWING) {
            this.triggerUpdateFigure(this.currentFigure, this.settings.tool === 'pencil' ? this.oldCurrentFigure : null);
            this.oldCurrentFigure = this.currentFigure.clone();
        }
    }

    handleNewFigure(figure) {
        if (!figure.id) {
            return;
        }

        if (!this.currentFigure || figure.id !== this.currentFigure.id) {
            this.boardController.upsertFigure(figure);
        }
    }

    handleRemoveFigure(id) {
        this.boardController.removeFigure(id);
    }

    handleUpdateFigure(id, data) {
        if (!this.currentFigure || id !== this.currentFigure.id) {
            const updatableFigure = this.boardController.figures.find((figure) => figure.id === id);

            if (updatableFigure) {
                const encoded = encode(updatableFigure);
                const decoded = decode(encoded + data);

                this.boardController.upsertFigure(decoded);
            }
        }
    }

    keyPressHandler(e) {
        if (e.ctrlKey && e.keyCode === 90) {
            if (this.history.length > 0) {
                const undoFunction = this.history.pop();
                undoFunction();
            }
        }
    }

    triggerUpdateFigure(newFigure, oldFigure) {
        if (!oldFigure) {
            this.figureWebSocket.changeFigure(newFigure);
        } else {
            this.figureWebSocket.updateFigure(newFigure, oldFigure);
        }
    }
}