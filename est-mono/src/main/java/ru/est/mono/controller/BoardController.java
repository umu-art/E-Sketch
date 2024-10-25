package ru.est.mono.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import ru.est.mono.api.BoardApi;
import ru.est.mono.model.BoardDto;
import ru.est.mono.model.BoardListDto;
import ru.est.mono.model.CreateRequest;
import ru.est.mono.model.ShareBoardDto;
import ru.est.mono.model.UnshareRequest;
import ru.est.mono.service.BoardService;

import java.util.UUID;

@Controller
@RequiredArgsConstructor
public class BoardController implements BoardApi {

    private final BoardService boardService;

    @Override
    public ResponseEntity<BoardListDto> callList() {
        return ResponseEntity.ok(boardService.getBoards());
    }

    @Override
    public ResponseEntity<BoardDto> getByUuid(UUID boardId) {
        return ResponseEntity.ok(boardService.getByUuid(boardId));
    }

    @Override
    public ResponseEntity<BoardDto> create(CreateRequest createRequest) {
        return ResponseEntity.ok(boardService.create(createRequest));
    }

    @Override
    public ResponseEntity<BoardDto> update(UUID boardId, CreateRequest createRequest) {
        return ResponseEntity.ok(boardService.update(boardId, createRequest));
    }

    @Override
    public ResponseEntity<Void> deleteBoard(UUID boardId) {
        boardService.delete(boardId);
        return ResponseEntity.ok().build();
    }

    @Override
    public ResponseEntity<Void> share(UUID boardId, ShareBoardDto shareBoardDto) {
        boardService.share(boardId, shareBoardDto);
        return ResponseEntity.ok().build();
    }

    @Override
    public ResponseEntity<Void> changeAccess(UUID boardId, ShareBoardDto shareBoardDto) {
        boardService.changeAccess(boardId, shareBoardDto);
        return ResponseEntity.ok().build();
    }

    @Override
    public ResponseEntity<Void> unshare(UUID boardId, UnshareRequest unshareRequest) {
        boardService.unshare(boardId, unshareRequest);
        return ResponseEntity.ok().build();
    }

    @Override
    public ResponseEntity<Void> connect(UUID boardId) {
        return ResponseEntity.status(101).build(); // TODO: заглушка
    }
}
