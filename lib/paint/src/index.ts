import { DefaultFigure } from 'figures/dist';

export class Board {
  svgElement: Element;
  figures: DefaultFigure[];

  public constructor(svgObject: Element) {
    this.svgElement = svgObject;
    this.figures = [];
    this.prepareBoard();
  }

  private prepareBoard() {
    this.svgElement.innerHTML = '';
    this.svgElement.setAttribute('controlled', 'est-paint');
  }

  public upsertFigure(figure: DefaultFigure): void {
    this.figures = this.figures.filter(f => f.id !== figure.id);

    this.figures.push(figure);
    this.renderFigure(figure);
  }

  public removeFigure(figureId: string): void {
    this.figures = this.figures
      .filter(f => f.id !== figureId);

    const figureElement = findFigureById(this.svgElement, figureId);
    if (figureElement) {
      this.svgElement.removeChild(figureElement);
    }
  }

  private comparePaths(path1: any, path2: any): any {
    const d1 = path1.getAttribute('d');
    const d2 = path2.getAttribute('d');

    if (d1 === d2) {
        return true;
    } else {
        return false;
    }
}

  private renderFigure(figure: DefaultFigure) {
    const figureElement = findFigureById(this.svgElement, figure.id);

    const newFigureElement = figure.toSvg(this.svgElement.ownerDocument);

    if (figureElement) {
        if (!this.comparePaths(figureElement, newFigureElement)) {
          this.svgElement.appendChild(newFigureElement);
          
          this.svgElement.removeChild(figureElement);
        }
    } else {
        this.svgElement.appendChild(newFigureElement);
    }
  }
}

function findFigureById(board: Element, id: string): Element | undefined {
  for (const figure of board.children) {
    if (figure.id === id) {
      return figure;
    }
  }
}