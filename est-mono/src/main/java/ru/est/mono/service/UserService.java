package ru.est.mono.service;

import ru.est.mono.domain.UserEntity;
import ru.est.mono.model.UserDto;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface UserService {

    boolean existsByUsername(String username);

    void register(UserEntity user);

    UserDto getSelf();

    Optional<UserDto> getUserById(UUID id);

    List<UserDto> searchByUsername(String username);
}
