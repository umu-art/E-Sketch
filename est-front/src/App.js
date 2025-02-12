import './global.scss';
import React from 'react';
import { BrowserRouter as Router, Navigate, Route, Routes } from 'react-router-dom';

import { Provider } from 'react-redux';

import SignUpPage from './Pages/Auth/SignUpPage/SignUpPage';
import SignInPage from './Pages/Auth/SignInPage/SignInPage';
import HomePage from './Pages/App/HomePage/HomePage';
import BoardPage from './Pages/App/BoardPage/BoardPage';
import store from './redux/store';
import EmailConfirmPage from './Pages/Auth/EmailConfirmPage/EmailConfirmPage';


function App() {
    return (
        <Provider store={store}>
        <Router>
            <Routes>
                <Route path='auth'>
                    <Route path="signup" element={<SignUpPage/>}/>
                    <Route path="signin" element={<SignInPage/>}/>
                    <Route path='confirm' element={<EmailConfirmPage />}/>
                    <Route path='*' element={<Navigate to="signin"/>}/>
                    <Route path='' element={<Navigate to="signin"/>}/>
                </Route>
                <Route path="app">
                    <Route path="home/*" element={<HomePage/>}/>
                    <Route path="board/:boardId" element={<BoardPage/>}/>
                    <Route path='*' element={<Navigate to="home"/>}/>
                    <Route path='' element={<Navigate to="home"/>}/>
                </Route>
                <Route path='*' element={<Navigate to="app"/>}/>
                <Route path='' element={<Navigate to="app"/>}/>
            </Routes>
        </Router>
        </Provider>
    );
}

export default App;
