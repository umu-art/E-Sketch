import { DefaultFigure } from 'figures/dist';

export class Board {
  svgElements: Element[];
  figures: DefaultFigure[];

  public constructor(svgElement: Element) {
    this.svgElements = [svgElement];
    this.figures = [];
    this.prepareBoard();
  }

  public addSvgElement(svgElement: Element): void {
    this.svgElements.push(svgElement);
    this.prepareBoard();
    for (const figure of this.figures) {
      this.renderFigure(figure);
    }
  }

  private prepareBoard() {
    for (const svgElement of this.svgElements) {
      svgElement.innerHTML = '';
      svgElement.setAttribute('controlled', 'est-paint');
    }
  }

  public upsertFigure(figure: DefaultFigure): void {
    this.figures = this.figures.filter(f => f.id !== figure.id);

    this.figures.push(figure);
    this.renderFigure(figure);
  }

  public removeFigure(figureId: string): void {
    this.figures = this.figures
      .filter(f => f.id !== figureId);

    for (const svgElement of this.svgElements) {
      const figureElement = findFigureById(svgElement, figureId);
      if (figureElement) {
        svgElement.removeChild(figureElement);
      }
    }
  }

  private renderFigure(figure: DefaultFigure) {
    for (const svgElement of this.svgElements) {
      let figureElement = findFigureById(svgElement, figure.id);

      if (figureElement) {
        svgElement.removeChild(figureElement);
      }

      figureElement = figure.toSvg(svgElement.ownerDocument);
      svgElement.appendChild(figureElement);
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