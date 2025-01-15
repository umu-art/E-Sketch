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

  private renderFigure(figure: DefaultFigure) {
    let figureElement = findFigureById(this.svgElement, figure.id);

    if (figureElement) {
      this.svgElement.removeChild(figureElement);
    }

    figureElement = figure.toSvg(this.svgElement.ownerDocument);
    this.svgElement.appendChild(figureElement);
  }
}

function findFigureById(board: Element, id: string): Element | undefined {
  for (const figure of board.children) {
    if (figure.id === id) {
      return figure;
    }
  }
}