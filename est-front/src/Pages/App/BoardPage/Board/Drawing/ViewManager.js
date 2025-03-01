import { Point } from 'figures/dist';
import store from '../../../../../redux/store';
import { MIN_SCALE, MAX_SCALE, BASE_OFFSET_X, BASE_OFFSET_Y } from './Constants';

export class ViewManager {
    constructor(board, settings) {
        this.board = board;
        this.settings = settings;

        this.moveState = {
            isMoving: false,
            start: new Point(0, 0),
        }

        store.subscribe(() => {
            const newState = store.getState();
            this.settings = newState;
            this.updateViewBox();
        });

        this.addEventListeners();
    }

    addEventListeners() {
        this.board.addEventListener('mousemove', (e) => this.handleMouseMove(e));
        this.board.addEventListener('wheel', (e) => this.handleWheel(e));
        this.board.addEventListener('mousedown', (e) => this.handleMouseDown(e));
        this.board.addEventListener('mouseup', (e) => this.handleMouseUp(e));
    }

    handleMouseDown(e) {
        e.preventDefault();

        if (e.button !== 2) {
            return;
        }

        this.moveState.isMoving = true;

        this.settings.view.startX = e.clientX;
        this.settings.view.startY = e.clientY;
    }

    handleMouseUp(e) {
        e.preventDefault();

        if (e.button !== 2) {
            return;
        }
        
        this.moveState.isMoving = false;
    }

    handleMouseMove(event) {
        if (!this.moveState.isMoving) {
            return;
        }

        const dx = event.clientX - this.settings.view.startX;
        const dy = event.clientY - this.settings.view.startY;

        this.settings.view.offsetX += dx / this.settings.view.scale;
        this.settings.view.offsetY += dy / this.settings.view.scale;

        this.settings.view.startX = event.clientX;
        this.settings.view.startY = event.clientY;

        this.updateViewBox();
    }

    handleWheel(e) {
        e.preventDefault();

        const scaleChange = this.getScaleChange(e.deltaY);
        const newScale = this.settings.view.scale * scaleChange;

        if (newScale < MIN_SCALE || newScale > MAX_SCALE) {
            return;
        }

        const cursorPosition = this.calculateCursorPosition(e);

        this.updateViewBoxWithScale(cursorPosition, scaleChange);
    }

    getScaleChange(deltaY) {
        return deltaY < 0 ? 1.03 : 0.97;
    }

    calculateCursorPosition(event) {
        const rect = this.board.getBoundingClientRect();

        const cursorX = (event.offsetX + BASE_OFFSET_X - rect.left) / this.settings.view.scale;
        const cursorY = (event.offsetY + BASE_OFFSET_Y - rect.top) / this.settings.view.scale;
        
        return { x: cursorX, y: cursorY };
    }

    updateViewBox() {
        const newWidth = this.board.clientWidth / this.settings.view.scale;
        const newHeight = this.board.clientHeight / this.settings.view.scale;

        this.board.setAttribute('viewBox', `${-this.settings.view.offsetX} ${-this.settings.view.offsetY} ${newWidth} ${newHeight}`);
    }

    updateViewBoxWithScale(cursorPosition, scaleChange = 1) {
        const nowX = cursorPosition.x - this.settings.view.offsetX;
        const nowY = cursorPosition.y - this.settings.view.offsetY;

        this.settings.view.scale *= scaleChange;

        this.settings.view.scale = Math.min(Math.max(this.settings.view.scale, MIN_SCALE), MAX_SCALE);

        this.settings.view.offsetX = -(nowX - cursorPosition.x / scaleChange);
        this.settings.view.offsetY = -(nowY - cursorPosition.y / scaleChange);

        const newWidth = this.board.clientWidth / this.settings.view.scale;
        const newHeight = this.board.clientHeight / this.settings.view.scale;

        this.board.setAttribute('viewBox', `${-this.settings.view.offsetX} ${-this.settings.view.offsetY} ${newWidth} ${newHeight}`);
    }
}