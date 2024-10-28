package ru.est.mono.service.impl;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.est.mono.domain.UserEntity;
import ru.est.mono.jpa.UserJpa;
import ru.est.mono.mapper.UserMapper;
import ru.est.mono.model.UserDto;
import ru.est.mono.service.UserService;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

import static ru.est.mono.configuration.auth.JwtAuthFilter.getPerformer;

@Slf4j
@Service
@RequiredArgsConstructor
public class UserServiceImpl implements UserService, UserDetailsService {

    private final UserJpa userJpa;
    private final UserMapper userMapper;

    @Override
    @Transactional(readOnly = true)
    public UserEntity loadUserByUsername(String username) throws UsernameNotFoundException {
        return userJpa.findByUsername(username)
                .orElseThrow(() -> new UsernameNotFoundException("User not found"));
    }

    @Override
    @Transactional(readOnly = true)
    public boolean existsByUsername(String username) {
        return userJpa.existsByUsername(username);
    }

    @Override
    public boolean existsByEmail(String email) {
        return userJpa.existsByEmail(email);
    }

    @Override
    @Transactional
    public void register(UserEntity user) {
        user.setId(UUID.randomUUID());
        userJpa.save(user);
    }

    @Override
    @Transactional(readOnly = true)
    public UserDto getSelf() {
        var entity = getPerformer();
        return userMapper.toDto(entity);
    }

    @Override
    @Transactional(readOnly = true)
    public Optional<UserDto> getUserById(UUID id) {
        var entity = userJpa.findById(id);
        return entity.map(userMapper::toDto);
    }

    @Override
    public Optional<UserDto> getByEmail(String email) {
        var entity = userJpa.findByEmail(email);
        return entity.map(userMapper::toDto);
    }

    @Override
    @Transactional(readOnly = true)
    public List<UserDto> searchByUsername(String username) {
        var entities = userJpa.searchByUsernameContainingIgnoreCase(username);
        return userMapper.toDto(entities);
    }
}
