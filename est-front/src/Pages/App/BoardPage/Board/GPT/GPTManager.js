import { Point } from "figures/dist";
import { DrawingStates, FPS } from "../Drawing/Constants";
import { toolToClass, toolToFigureType } from "../Drawing/ToolMapping";
import { setGPTStatus, showGPTPopover } from "../../../../../redux/toolSettings/actions";
import store from "../../../../../redux/store";

export class GPTManager {
    constructor(drawingManager) {
        this.drawingManager = drawingManager;
        this.addEventListeners();
    }

    addEventListeners() {
        this.drawingManager.board.addEventListener('mousedown', (e) => this.handleMouseDown(e));
        this.drawingManager.board.addEventListener('mouseup', (e) => this.handleMouseUp(e));

        setInterval(() => this.updateDrawing(), 1000 / FPS);
    }

    handleMouseDown(e) {
        if (e.button !== 0 || this.drawingManager.settings.tool !== 'gpt' || this.drawingManager.settings.state !== DrawingStates.IDLE) {
            return;
        }

        this.drawingManager.settings.state = DrawingStates.SELECTING;

        this.drawingManager.currentFigure = toolToClass['rectangle'].startProcess(
            toolToFigureType['rectangle'],
            'gpt-selection',
            Object.values({
                lineColor: '#1677ff',
                fillColor: '#1677ff15',
                lineWidth: 2,
            }),
            new Point(this.drawingManager.drawing.nowX, this.drawingManager.drawing.nowY)
        );
    }

    updateDrawing() {
        if (this.drawingManager.settings.state !== DrawingStates.SELECTING) {
            return;
        }

        this.drawingManager.currentFigure.process(new Point(this.drawingManager.drawing.nowX, this.drawingManager.drawing.nowY));
        this.drawingManager.boardController.upsertFigure(this.drawingManager.currentFigure);
    }

    handleMouseUp(e) {
        if (e.button !== 0 || this.drawingManager.settings.tool !== 'gpt' || !this.drawingManager.currentFigure) {
            return;
        }

        this.drawingManager.settings.state = DrawingStates.IDLE;

        store.dispatch(setGPTStatus("processing"));

        store.dispatch(showGPTPopover({
            leftUp: this.drawingManager.currentFigure.points[0],
            rightDown: this.drawingManager.currentFigure.points[1],
        }));

        this.drawingManager.boardController.removeFigure(this.drawingManager.currentFigure.id);
        this.drawingManager.figureWebSocket.deleteFigure(this.drawingManager.currentFigure.id);

        this.drawingManager.currentFigure = null;
    }
}