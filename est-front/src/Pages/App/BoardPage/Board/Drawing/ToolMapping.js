import { FigureType, Line, Rectangle, Ellipse } from 'figures/dist';

export const toolToClass = {
    'pencil': Line,
    'rectangle': Rectangle,
    'ellipse': Ellipse,
};

export const toolToFigureType = {
    'pencil': FigureType.LINE,
    'rectangle': FigureType.RECTANGLE,
    'ellipse': FigureType.ELLIPSE,
};
