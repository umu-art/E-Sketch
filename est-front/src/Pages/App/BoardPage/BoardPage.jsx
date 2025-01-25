import { Flex } from 'antd';
import React, { useEffect, useState } from 'react';
import Board from './Board/Board';
import { useNavigate, useParams } from 'react-router-dom';
import ToolPanel from './ToolPanel/ToolPanel';
import { BoardApi } from 'est_proxy_api';
import LoadingPage from '../../LoadingPage/LoadingPage';
import HeadMenu from './HeadMenu/HeadMenu';

import { drawing } from './Board/Paint';
import ErrorPage from '../../ErrorPages/ErrorPage';
import ScaleElement from './Scale/Scale';

const apiInstance = new BoardApi();

const BoardPage = () => {
  const { boardId } = useParams();
  const [data, setData] = useState(null);
  const [tool, setTool] = useState(drawing.tool);
  const [err, setErr] = useState(null);

  const navigate = useNavigate();

  const updateData = (newData) => {
    apiInstance.update(boardId, { 'createRequest': newData }).then((respData) => {
      setData(respData);
    }).catch((error) => {
      console.log(error);
    });
  };

  const refreshBoardData = () => {
    apiInstance.getByUuid(boardId).then((data) => {
      setData(data);
    }).catch((error) => {
      if (error.code === 401) {
        navigate('/auth/signin');
      } else {
        setErr(error);
      }
    });
  };


  useEffect(() => {
    if (!data && !err) {
      refreshBoardData();
    }
  });

  if (err) {
    return (
      <ErrorPage status={err.status === 400 ? 404 : err.status} />
    );
  }

  if (!data) {
    return (
      <LoadingPage />
    );
  }

  return (
    <>
      <Board className="h100vh w100vw" style={{ position: 'absolute' }} boardId={boardId} currentTool={tool} />
      <ScaleElement style={{ position: 'absolute', right: '0', bottom: '0', zIndex: '6' }}></ScaleElement>
      { /* Menu wrap */}
      <Flex className="h100vh w100vw" style={{ padding: '20px 20px', position: 'absolute' }} vertical
            align="center" justify="space-between">
        { /* Top */}
        <HeadMenu data={data} updateData={updateData} refreshData={refreshBoardData} />
        { /* Bottom */}
        <Flex className="w100p" justify="center">
          <ToolPanel onToolChange={(tool) => {
            console.log(tool);
            drawing.tool = tool;
            setTool(tool);
          }} />
        </Flex>
      </Flex>
    </>

  );
};

export default BoardPage;