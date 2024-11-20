import { Flex } from 'antd';
import React from 'react';

import { LoadingOutlined } from '@ant-design/icons';

const LoadingPage = () => {
    return (
        <Flex className='w100vw h100vh' style={{ fontSize: 50 }} justify='center' align='center'>
            <Flex style={{ width: 100, height: 100 }} justify='center' align='center'>
                <LoadingOutlined/>
            </Flex>
        </Flex>
    );
};

export default LoadingPage;