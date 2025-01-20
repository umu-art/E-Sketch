import React from 'react';

import { Navigate, Route, Routes } from 'react-router-dom';

import { Flex } from 'antd';

import Menu from './Menu/Menu.jsx';
import MyBoardsPage from './Pages/MyBoardsPage/MyBoardsPage.jsx';
import SharedBoardsPage from './Pages/SharedBoardsPage/SharedBoardsPage.jsx';


const HomePage = () => {
    return (
        <Flex className="mh100vh" style={{ padding: "50px 100px" }}>
            <Flex className='h100p'>
                <Menu/>
            </Flex>
            <Routes>
                <Route path='my' element={<MyBoardsPage />}/>
                <Route path='shared' element={<SharedBoardsPage />}/>
                <Route path='*' element={<Navigate to="my" />}/>
                <Route path='' element={<Navigate to="my" />}/>
            </Routes>
        </Flex>
    );
};

export default HomePage;