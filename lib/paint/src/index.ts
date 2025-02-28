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

  private compare(path1: any, path2: any): any {
    return path1 === path2;
  }

  private renderFigure(figure: DefaultFigure) {
    const figureElement = findFigureById(this.svgElement, figure.id);

    const newFigureElement = figure.toSvg(this.svgElement.ownerDocument);

    if (figureElement) {
      if (!this.compare(figureElement, newFigureElement)) {
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