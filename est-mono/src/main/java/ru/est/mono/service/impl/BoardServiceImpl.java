package ru.est.mono.service.impl;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.est.mono.domain.BoardEntity;
import ru.est.mono.domain.BoardSharingEntity;
import ru.est.mono.domain.UserEntity;
import ru.est.mono.exception.NotAllowedException;
import ru.est.mono.exception.NotFoundException;
import ru.est.mono.jpa.BoardJpa;
import ru.est.mono.jpa.BoardSharingJpa;
import ru.est.mono.jpa.UserJpa;
import ru.est.mono.mapper.BoardMapper;
import ru.est.mono.model.BoardDto;
import ru.est.mono.model.BoardListDto;
import ru.est.mono.model.CreateRequest;
import ru.est.mono.model.ShareBoardDto;
import ru.est.mono.model.SharingDto;
import ru.est.mono.model.UnshareRequest;
import ru.est.mono.service.BoardService;

import java.util.UUID;

import static ru.est.mono.configuration.auth.JwtAuthFilter.getPerformer;

@Service
@RequiredArgsConstructor
public class BoardServiceImpl implements BoardService {

    private final UserJpa userJpa;
    private final BoardJpa boardJpa;
    private final BoardSharingJpa boardSharingJpa;
    private final BoardMapper boardMapper;

    @Override
    @Transactional(readOnly = true)
    public BoardListDto getBoards() {
        var executorUsername = getPerformer().getUsername();
        var entities = boardJpa.findAllForUser(executorUsername);
        var boardListDto = new BoardListDto();

        entities.forEach(entity -> {
            var boardDto = boardMapper.toBoardDto(entity);
            if (entity.getOwner().getUsername().equals(executorUsername)) {
                boardListDto.addMineItem(boardDto);
            } else {
                boardListDto.addSharedItem(boardDto);
            }
        });

        return boardListDto;
    }

    @Override
    @Transactional(readOnly = true)
    public BoardDto getByUuid(UUID boardId) {
        checkAdminAccess(boardId);
        var entity = boardJpa.findById(boardId)
                .orElseThrow(() -> new NotFoundException(BoardEntity.class, boardId));
        return boardMapper.toBoardDto(entity);
    }

    @Override
    @Transactional
    public BoardDto create(CreateRequest createRequest) {
        var entity = boardMapper.toBoardEntity(createRequest);
        entity.setId(UUID.randomUUID());
        entity.setLinkSharedMode(BoardDto.LinkSharedModeEnum.NONE_BY_LINK);
        entity.setOwner(getPerformer());
        entity = boardJpa.save(entity);

        return boardMapper.toBoardDto(entity);
    }

    @Override
    @Transactional
    public BoardDto update(UUID boardId, CreateRequest createRequest) {
        checkAdminAccess(boardId);
        var entity = boardJpa.findById(boardId)
                .orElseThrow(() -> new NotFoundException(BoardEntity.class, boardId));

        boardMapper.updateBoardEntity(entity, createRequest);
        entity = boardJpa.save(entity);

        return boardMapper.toBoardDto(entity);
    }

    @Override
    @Transactional
    public void delete(UUID boardId) {
        checkAdminAccess(boardId);
        boardJpa.deleteById(boardId);
    }

    @Override
    @Transactional
    public void share(UUID boardId, ShareBoardDto shareBoardDto) {
        checkAdminAccess(boardId);
        var board = boardJpa.findById(boardId)
                .orElseThrow(() -> new NotFoundException(BoardEntity.class, boardId));

        var user = userJpa.findById(shareBoardDto.getUserId())
                .orElseThrow(() -> new NotFoundException(UserEntity.class, shareBoardDto.getUserId()));

        if (board.getBoardSharings()
                .stream()
                .anyMatch(sharing -> sharing.getUser().equals(user))) {
            return;
        }

        var boardSharingEntity = new BoardSharingEntity();
        boardSharingEntity.setId(UUID.randomUUID());
        boardSharingEntity.setUser(user);
        boardSharingEntity.setBoard(board);
        boardSharingEntity.setSharingMode(boardMapper.toSharingMode(shareBoardDto.getAccess()));

        boardSharingJpa.save(boardSharingEntity);
    }

    @Override
    @Transactional
    public void unshare(UUID boardId, UnshareRequest unshareRequest) {
        checkAdminAccess(boardId);
        var board = boardJpa.findById(boardId)
                .orElseThrow(() -> new NotFoundException(BoardEntity.class, boardId));

        var user = userJpa.findById(unshareRequest.getUserId())
                .orElseThrow(() -> new NotFoundException(UserEntity.class, unshareRequest.getUserId()));

        var boardSharingEntity = board.getBoardSharings().stream()
                .filter(sharing -> sharing.getUser().equals(user))
                .findFirst()
                .orElseThrow(() -> new NotFoundException(BoardSharingEntity.class, unshareRequest.getUserId()));

        boardSharingJpa.delete(boardSharingEntity);
    }

    @Override
    public void changeAccess(UUID boardId, ShareBoardDto shareBoardDto) {
        checkAdminAccess(boardId);
        var board = boardJpa.findById(boardId)
                .orElseThrow(() -> new NotFoundException(BoardEntity.class, boardId));

        var user = userJpa.findById(shareBoardDto.getUserId())
                .orElseThrow(() -> new NotFoundException(UserEntity.class, shareBoardDto.getUserId()));

        var boardSharingEntity = board.getBoardSharings().stream()
                .filter(sharing -> sharing.getUser().equals(user))
                .findFirst()
                .orElseThrow(() -> new NotFoundException(BoardSharingEntity.class, shareBoardDto.getUserId()));

        boardSharingEntity.setSharingMode(boardMapper.toSharingMode(shareBoardDto.getAccess()));
        boardSharingJpa.save(boardSharingEntity);
    }

    private void checkAdminAccess(UUID boardId) {
        var executorUsername = getPerformer().getUsername();

        var board = boardJpa.findById(boardId)
                .orElseThrow(() -> new NotFoundException(BoardEntity.class, boardId));

        if (executorUsername.equals(board.getOwner().getUsername())) {
            return;
        }

        if (board.getBoardSharings()
                .stream()
                .anyMatch(boardSharingEntity ->
                        executorUsername.equals(boardSharingEntity.getUser().getUsername()) &&
                                boardSharingEntity.getSharingMode() == SharingDto.AccessEnum.ADMIN)) {
            return;
        }

        throw new NotAllowedException();
    }
}
