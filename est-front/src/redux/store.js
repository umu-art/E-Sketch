import { createStore } from 'redux';
import drawingReducer from './toolSettings/reducer';

const store = createStore(drawingReducer);

export default store;
