package ru.est.mono.service;

import ru.est.mono.model.BoardDto;
import ru.est.mono.model.BoardListDto;
import ru.est.mono.model.CreateRequest;
import ru.est.mono.model.ShareBoardDto;
import ru.est.mono.model.UnshareRequest;

import java.util.UUID;

public interface BoardService {

    BoardListDto getBoards();

    BoardDto getByUuid(UUID boardId);

    BoardDto create(CreateRequest createRequest);

    BoardDto update(UUID boardId, CreateRequest createRequest);

    void delete(UUID boardId);

    void share(UUID boardId, ShareBoardDto shareBoardDto);

    void unshare(UUID boardId, UnshareRequest unshareRequest);

    void changeAccess(UUID boardId, ShareBoardDto shareBoardDto);
}
