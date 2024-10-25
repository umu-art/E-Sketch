package ru.est.mono.mapper;

import org.mapstruct.Mapper;
import ru.est.mono.domain.UserEntity;
import ru.est.mono.model.UserDto;

import java.util.List;
import java.util.Optional;

@Mapper(componentModel = "spring")
public interface UserMapper {

    UserDto toDto(UserEntity entity);

    List<UserDto> toDto(List<UserEntity> entity);
}
