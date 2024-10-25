package ru.est.mono.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.MappingTarget;
import ru.est.mono.domain.BoardEntity;
import ru.est.mono.domain.BoardSharingEntity;
import ru.est.mono.model.BoardDto;
import ru.est.mono.model.CreateRequest;
import ru.est.mono.model.ShareBoardDto;
import ru.est.mono.model.SharingDto;

@Mapper(componentModel = "spring")
public interface BoardMapper {

    @Mapping(target = "userInfo", source = "user")
    @Mapping(target = "access", source = "sharingMode")
    SharingDto toSharingDto(BoardSharingEntity entity);

    @Mapping(target = "ownerInfo", source = "owner")
    @Mapping(target = "sharedWith", source = "boardSharings")
    BoardDto toBoardDto(BoardEntity entity);

    BoardEntity toBoardEntity(CreateRequest createRequest);

    void updateBoardEntity(@MappingTarget BoardEntity entity, CreateRequest createRequest);

    SharingDto.AccessEnum toSharingMode(ShareBoardDto.AccessEnum access);
}
