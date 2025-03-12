import React from 'react';
import { Button, Card, Divider, Flex } from 'antd';
import Icon, { EditOutlined, EllipsisOutlined, OpenAIOutlined } from '@ant-design/icons';

import { ReactComponent as EraserOutlined } from './eraser.svg';
import { ReactComponent as RectangleOutlined } from './rectangle.svg';
import { ReactComponent as EllipseOutlined } from './ellipse.svg';

import ToolButton from './tools/ToolButton';

const ToolPanel = () => {
  return (
    <Card size="small" className="shadow" style={{ zIndex: 15 }}>
      <Flex gap="small" align="center">
        <ToolButton
          tool="pencil"
          icon={<EditOutlined />}
          showColorChange
          showWidthChange
        />
        <ToolButton
          tool="eraser"
          icon={<Icon component={EraserOutlined} />}
          showWidthChange
        />
        <ToolButton
          tool="rectangle"
          icon={<Icon component={RectangleOutlined} />}
          showWidthChange
          showColorChange
          showFillColorChange
        />
        <ToolButton
          tool="ellipse"
          icon={<Icon component={EllipseOutlined} />}
          showWidthChange
          showColorChange
          showFillColorChange
        />
        <ToolButton
          tool="gpt"
          icon={<Icon component={OpenAIOutlined} />}  
        />
        <Divider type="vertical" style={{ height: 30 }} />
        <Button 
          icon={<EllipsisOutlined />}
        />
      </Flex>
    </Card>
  );
};

export default ToolPanel;