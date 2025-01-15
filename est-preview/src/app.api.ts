import { DefaultFigure } from 'figures/dist';
import { decode } from 'coder/dist';

const backUrl = process.env.BACK_URL || 'http://est-back.e-sketch.svc.cluster.local';

export async function getAllFigures(boardId: string): Promise<DefaultFigure[]> {
  const url = `${backUrl}/back/figure/list/${boardId}`;

  return fetch(url)
    .then(response => response.json())
    .then(data => data.figures
      .filter(raw => raw.data.length > 0)
      .map(raw => decode(raw.data))
    );
}