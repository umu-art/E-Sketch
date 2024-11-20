import './global.scss';
import React from 'react';
import { BrowserRouter as Router, Navigate, Route, Routes } from 'react-router-dom';

import SignUpPage from './Pages/Auth/SignUpPage/SignUpPage';
import SignInPage from './Pages/Auth/SignInPage/SignInPage';
import HomePage from './Pages/App/HomePage/HomePage';
import BoardPage from './Pages/App/BoardPage/BoardPage';


function App() {
    return (
      <Router>
          <Routes>
            <Route path='auth'>
                <Route path="signup" element={<SignUpPage />} />
                <Route path="signin" element={<SignInPage />} />
            </Route>
            <Route path="app">
                <Route path="home/*" element={<HomePage />}/>
                <Route path="board/:boardId" element={<BoardPage />}/>
                <Route path='*' element={<Navigate to="home" />}/>
                <Route path='' element={<Navigate to="home" />}/>
            </Route>
            <Route path='*' element={<Navigate to="app" />} />
            <Route path='' element={<Navigate to="app" />} />
          </Routes>
      </Router>
    );
  } 

export default App;
