import React, { useState } from 'react';
import { Button, Card, Divider, Flex } from 'antd';
import Icon, { EditOutlined, EllipsisOutlined } from '@ant-design/icons';

import { ReactComponent as EraserOutlined } from './eraser.svg';

const ToolPanel = ({ onToolChange }) => {
    const [selectedTool, setSelectedTool] = useState('pencil');

    const handleToolChange = (tool) => {
        setSelectedTool(tool);
        onToolChange(tool);
    };

    return (
        <Card size='small' className='shadow'>
            <Flex gap="small" align='center'>
                <Button 
                    type={selectedTool === 'pencil' ? 'primary' : 'default'}
                    icon={<EditOutlined />}
                    onClick={() => handleToolChange('pencil')}
                    key='pencil'
                />
                <Button 
                    type={selectedTool === 'eraser' ? 'primary' : 'default'}
                    icon={<Icon component={EraserOutlined}/>}
                    onClick={() => handleToolChange('eraser')}
                    key='eraser'
                />
                <Divider type='vertical' style={{height: 30}}/>
                <Button 
                    type={'default'}
                    icon={<EllipsisOutlined key="ellipsis" />}
                    key='ellipsis'
                ></Button>
            </Flex>
            
        </Card>
    );
};

export default ToolPanel;