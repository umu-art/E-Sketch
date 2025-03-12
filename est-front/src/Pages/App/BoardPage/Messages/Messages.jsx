import React, { useEffect, useRef, useState, useCallback } from 'react';
import { Card, Button, Flex } from 'antd';
import { useDispatch, useSelector } from 'react-redux';
import { removeMessage } from '../../../../redux/toolSettings/actions';
import 'katex/dist/katex.min.css';
import {
    CloseOutlined,
    UpOutlined,
    DownOutlined,
    LeftOutlined,
    RightOutlined,
    ExpandOutlined,
    CompressOutlined
} from '@ant-design/icons';
import renderMarkdownWithLatex from './renderMarkdownWithLatex';
import PropTypes from 'prop-types';

const Messages = () => {
    const dispatch = useDispatch();
    const messages = useSelector(state => [...state.messages]);
    const [activeIndex, setActiveIndex] = useState(0);
    const [isCollapsed, setIsCollapsed] = useState(false);
    const [isExpanded, setIsExpanded] = useState(false);

    const prevMessagesLengthRef = useRef(messages.length);

    useEffect(() => {
        if (messages.length > prevMessagesLengthRef.current) {
            setActiveIndex(messages.length - 1);
        }
        prevMessagesLengthRef.current = messages.length;
    }, [messages.length]);

    const handleClose = useCallback((id) => {
        dispatch(removeMessage(id));
        if (activeIndex >= messages.length - 1) {
            setActiveIndex(Math.max(0, messages.length - 2));
        }
    }, [dispatch, activeIndex, messages.length]);

    const handleNext = useCallback(() => {
        setActiveIndex((prevIndex) => (prevIndex + 1) % messages.length);
    }, [messages.length]);

    const handlePrev = useCallback(() => {
        setActiveIndex((prevIndex) => (prevIndex - 1 + messages.length) % messages.length);
    }, [messages.length]);

    const toggleCollapse = useCallback(() => {
        setIsCollapsed(prev => !prev);
    }, []);

    const toggleExpand = useCallback(() => {
        setIsExpanded(prev => !prev);
    }, []);

    return (
        <Flex style={{ zIndex: 15, position: "absolute", right: 20, width: !isCollapsed && isExpanded ? 600 : 300 }} vertical gap="small" align="center">
            {messages.length > 0 && (
                <MessageCard
                    message={messages[activeIndex]}
                    index={activeIndex}
                    total={messages.length}
                    isCollapsed={isCollapsed}
                    isExpanded={isExpanded}
                    onClose={handleClose}
                    onToggleCollapse={toggleCollapse}
                    onToggleExpand={toggleExpand}
                />
            )}
            {messages.length > 1 && !isCollapsed && (
                <NavigationButtons onPrev={handlePrev} onNext={handleNext} />
            )}
        </Flex>
    );
};

const MessageCard = React.memo(({ message, index, total, isCollapsed, isExpanded, onClose, onToggleCollapse, onToggleExpand }) => (
    <Card
        className='shadow'
        key={message.id}
        title={`(${index + 1}/${total}) ${message.title}`}
        extra={
            <Flex gap="small">
                {!isCollapsed && [
                    <Button type="link" icon={<CloseOutlined />} onClick={() => onClose(message.id)} key="close" />,
                    <Button type="link" icon={isExpanded ? <CompressOutlined /> : <ExpandOutlined />} onClick={onToggleExpand} key="expand" />,
                ]}
                <Button type="link" icon={isCollapsed ? <DownOutlined /> : <UpOutlined />} onClick={onToggleCollapse} />
            </Flex>
        }
        style={{ width: !isCollapsed && isExpanded ? 600 : 300 }}
        styles={{
            body: { padding: isCollapsed ? 0 : undefined }
        }}
    >
        <Flex vertical style={{ 
            height: (isCollapsed ? 0 : 1) * (isExpanded ? 500 : 250), 
            overflowY: 'auto' 
        }}>
            {renderMarkdownWithLatex(message.content)}
        </Flex>
    </Card>
));

MessageCard.propTypes = {
    message: PropTypes.object.isRequired,
    index: PropTypes.number.isRequired,
    total: PropTypes.number.isRequired,
    isCollapsed: PropTypes.bool.isRequired,
    isExpanded: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onToggleCollapse: PropTypes.func.isRequired,
    onToggleExpand: PropTypes.func.isRequired,
};

const NavigationButtons = React.memo(({ onPrev, onNext }) => (
    <Flex gap="small" justify="center" className="w100p">
        <Button icon={<LeftOutlined />} onClick={onPrev} style={{ width: '-webkit-fill-available' }} />
        <Button icon={<RightOutlined />} onClick={onNext} style={{ width: '-webkit-fill-available' }} />
    </Flex>
));

NavigationButtons.propTypes = {
    onPrev: PropTypes.func.isRequired,
    onNext: PropTypes.func.isRequired,
};

export default Messages;