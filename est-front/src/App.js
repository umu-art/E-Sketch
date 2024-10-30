import './global.scss';
import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import SignUpPage from './Pages/Auth/SignUpPage/SignUpPage';
import SignInPage from './Pages/Auth/SignInPage/SignInPage';


function App() {
    return (
      <Router>
          <Routes>
            <Route path='auth'>
                <Route path="signup" element={<SignUpPage />} />
                <Route path="signin" element={<SignInPage />} />
            </Route>
          </Routes>
      </Router>
    );
  } 

export default App;
