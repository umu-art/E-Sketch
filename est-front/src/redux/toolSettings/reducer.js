import { ADD_MESSAGE, HIDE_GPT_POPOVER, REMOVE_MESSAGE, SET_FILL_COLOR, SET_GPT_STATUS, SET_LINE_COLOR, SET_LINE_WIDTH, SET_OFFSET, SET_SCALE, SET_STATE, SET_TOOL, SHOW_GPT_POPOVER } from './actions';
import { MIN_SCALE, MAX_SCALE } from '../../Pages/App/BoardPage/Board/Drawing/Constants';

const initialState = {
  state: 'idle',
  tool: 'pencil',
  tools: {
    'pencil': {
      lineColor: '#1677ff',
      lineWidth: 2,
    },
    'eraser': {
      lineWidth: 5,
    },
    'rectangle': {
      lineColor: '#1677ff',
      fillColor: '#ffffff00',
      lineWidth: 2,
    },
    'ellipse': {
      lineColor: '#1677ff',
      fillColor: '#ffffff00',
      lineWidth: 2,
    },
    'gpt': {
      status: null,
    },
  },
  view: {
    scale: 1,
    offsetX: 0,
    offsetY: 0,
    startX: 0,
    startY: 0,
  },
  messages: [
    {
      id: "1",
      title: "GPT",
      content: "Похоже, вы нарисовали уравнение \\(x^2 = 1\\).\n\nРешение: \n\n\\[\nx^2 = 1\n\\]\n\nЭто квадратное уравнение принято в стандартной форме:"
    },
    {
      id: "2",
      title: "GPT",
      content: "Уравнение \\(\\cos^2 x = 1\\) решается следующим образом:\n\n1. Извлечем квадратный корень из обеих частей: \n\\[\n   \\cos x = \\pm 1\n   \\]\n\n2. Решим уравнения:\n   - \\(\\cos x = 1\\), что достигается при \\(x = 2\\pi k\\), где \\(k\\) — любое целое число.\n   - \\(\\cos x = -1\\), что достигается при \\(x = \\pi + 2\\pi k\\), где \\(k\\) — любое целое число.\n\nТаким образом, решения:\n\\[\nx = 2\\pi k \\quad \\text{и} \\quad x = \\pi (1 + 2k), \\quad k \\in \\mathbb{Z}.\n\\]"
    },
  ],
  popover: {
    visible: false,
    request: null,
  },
};

const drawingReducer = (state = initialState, action) => {
  switch (action.type) {
    case SET_TOOL:
      return { ...state, tool: action.payload };
    case SET_STATE:
      return { ...state, state: action.payload };
    case SET_LINE_COLOR: {
      const currentTool = state.tools[state.tool];
      return {
        ...state,
        tools: {
          ...state.tools,
          [state.tool]: {
            ...currentTool,
            lineColor: action.payload,
          },
        },
      };
    }
    case SET_FILL_COLOR: {
      const currentTool = state.tools[state.tool];
      return {
        ...state,
        tools: {
          ...state.tools,
          [state.tool]: {
            ...currentTool,
            fillColor: action.payload,
          },
        },
      };
    }
    case SET_LINE_WIDTH: {
      const currentTool = state.tools[state.tool];
      return {
        ...state,
        tools: {
          ...state.tools,
          [state.tool]: {
            ...currentTool,
            lineWidth: action.payload,
          },
        },
      };
    }
    case SET_SCALE: {
      const boundedScale = Math.min(Math.max(action.payload, MIN_SCALE), MAX_SCALE);

      return {
        ...state,
        view: {
          ...state.view,
          scale: boundedScale,
        },
      };
    }
    case SET_OFFSET:
      return {
        ...state,
        view: {
          ...state.view,
          offsetX: action.payload.offsetX,
          offsetY: action.payload.offsetY,
        },
      };
    case SET_GPT_STATUS:
      return {
        ...state,
        tools: {
          ...state.tools,
          'gpt': {
            status: action.payload,
          },
        },
      };
    case ADD_MESSAGE:
      return {
        ...state,
        messages: [...state.messages, action.payload],
      };
    case REMOVE_MESSAGE:
      return {
        ...state,
        messages: state.messages.filter(message => message.id !== action.payload),
      };
    case SHOW_GPT_POPOVER:
      return {
        ...state,
        popover: {
          visible: true,
          request: action.payload,
        },
      };
    case HIDE_GPT_POPOVER:
      return {
        ...state,
        popover: {
          visible: false,
          request: null,
        },
      };
    default:
      return state;
  }
};

export default drawingReducer;