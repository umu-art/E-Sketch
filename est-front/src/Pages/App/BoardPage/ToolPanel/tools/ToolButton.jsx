import { Badge, Button, ColorPicker, Flex, InputNumber, Popover } from 'antd';
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setFillColor, setLineColor, setLineWidth, setTool } from '../../../../../redux/toolSettings/actions';
import PropTypes from 'prop-types';


const basicColors = [
    { label: 'Black', color: '#000000' },
    { label: 'White', color: '#FFFFFF' },
    { label: 'Blue', color: '#1677ff' },
    { label: 'Red', color: '#FF0000' },
    { label: 'Green', color: '#00FF00' },
    { label: 'Yellow', color: '#FFFF00' },
];

const presets = [{
    label: 'Basic Colors',
    colors: basicColors.map(color => color.color),
    key: 'basicColors',
}]

const ToolButton = ({ 
    tool,
    icon,
    showColorChange = false,
    showWidthChange = false,
    showFillColorChange = false,
}) => {
    const dispatch = useDispatch();
    const selectedTool = useSelector((state) => state.tool);
    const selectedToolStatus = useSelector((state) => state.tools[tool].status);
    const selectedColor = useSelector((state) => state.tools[tool].lineColor);
    const selectedWidth = useSelector((state) => state.tools[tool].lineWidth);
    const selectedFillColor = useSelector((state) => state.tools[tool].fillColor);

    const changeTool = (newTool) => {
        dispatch(setTool(newTool));
    };

    const changeColor = (newColor) => {
        dispatch(setLineColor(newColor));
    };

    const changeFillColor = (newColor) => {
        console.log(newColor);  
        dispatch(setFillColor(newColor));
    };

    const changeWidth = (newWidth) => {
        dispatch(setLineWidth(newWidth));
    };

    const handleClick = () => {
        changeTool(tool);
    };

    return (
        <Popover
            content={
                <Flex gap="small">
                    {
                        showColorChange ? 
                        <ColorPicker disabledAlpha disabledFormat presets={presets} value={selectedColor} onChangeComplete={(color) => changeColor(color.toHexString())} />
                        :
                        null
                    }
                    {
                        showFillColorChange ? 
                        <ColorPicker disabledFormat presets={presets} value={selectedFillColor} onChangeComplete={(color) => changeFillColor(color.toHexString())} />
                        :
                        null
                    }
                    {
                        showWidthChange ?
                        <InputNumber
                            defaultValue={selectedWidth} 
                            changeOnWheel
                            min={1}
                            max={100}
                            onMouseLeave={(e) => e.target.blur()}
                            onChange={(value) => changeWidth(value)}
                            parser={(value) => {
                                if (isNaN(value)) {
                                    return 2;
                                }
                                
                                return Math.floor(value);
                            }}
                            style={{
                                width: '65px',
                            }}
                        />
                        :
                        null
                    }
                </Flex>
            }
            trigger={(showColorChange || showFillColorChange || showWidthChange) && "click"}
        >
            <Badge status={selectedToolStatus} dot={selectedToolStatus}>
                <Button
                    type={selectedTool === tool ? 'primary' : 'default'}
                    icon={icon}
                    onClick={handleClick}
                    key="pencil"
                />
            </Badge>
        </Popover>
    );
};

ToolButton.propTypes = {
    tool: PropTypes.string.isRequired,
    icon: PropTypes.node.isRequired,
    showColorChange: PropTypes.bool,
    showWidthChange: PropTypes.bool,
    showFillColorChange: PropTypes.bool,
};

ToolButton.defaultProps = {
    showColorChange: false,
    showWidthChange: false,
    showFillColorChange: false,
};

export default ToolButton;