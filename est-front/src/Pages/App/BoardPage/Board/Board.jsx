import React from 'react';

const Board = ({ className, style, gridStep = 100, boardId, ...props }) => {
    const lines = [];
    const width = "100vw";
    const height = "100vh";

    // Генерация вертикальных линий
    for (let x = gridStep; x < window.innerWidth; x += gridStep) {
        lines.push(<line key={`v-${x}`} x1={x} y1={0} x2={x} y2={window.innerHeight} stroke="lightgray"
                         strokeWidth="1"/>);
    }

    // Генерация горизонтальных линий
    for (let y = gridStep; y < window.innerHeight; y += gridStep) {
        lines.push(<line key={`h-${y}`} x1={0} y1={y} x2={window.innerWidth} y2={y} stroke="lightgray"
                         strokeWidth="1"/>);
    }

    // Определяем размеры шрифта и рассчитываем центр
    const fontSize = 20;
    const text = "Board " + boardId;
    const xCenter = window.innerWidth / 2;
    const yCenter = window.innerHeight / 2 + fontSize / 2; // Сдвинуть вниз на половину высоты шрифта

    return (
        <svg
            width={width}
            height={height}
            className={className}
            style={{ ...style, backgroundColor: 'white', overflow: 'hidden' }}
        >
            {lines}
            <text x={xCenter} y={yCenter - 50} fill="black" fontSize={fontSize} textAnchor="middle"
                  dominantBaseline="middle">
                {text}
            </text>
        </svg>
    );
};

export default Board;