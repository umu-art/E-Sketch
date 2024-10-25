package ru.est.mono.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Controller;
import ru.est.mono.api.UserApi;
import ru.est.mono.domain.UserEntity;
import ru.est.mono.model.AuthDto;
import ru.est.mono.model.RegisterDto;
import ru.est.mono.model.UserDto;
import ru.est.mono.service.UserService;
import ru.est.mono.service.impl.JwtService;

import java.util.List;
import java.util.UUID;

import static ru.est.mono.configuration.auth.JwtAuthFilter.AUTH_COOKIE;

@Controller
@RequiredArgsConstructor
public class UserController implements UserApi {

    private final AuthenticationManager authenticationManager;
    private final PasswordEncoder passwordEncoder;
    private final JwtService jwtService;
    private final UserService userService;

    @Override
    public ResponseEntity<Void> register(RegisterDto registerDto) {
        if (userService.existsByUsername(registerDto.getUsername())) {
            return ResponseEntity.badRequest().build();
        }

        UserEntity user = new UserEntity();
        user.setUsername(registerDto.getUsername());
        user.setPasswordHash(passwordEncoder.encode(registerDto.getPasswordHash()));
        user.setEmail(registerDto.getEmail());

        userService.register(user);

        var authentication = authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(registerDto.getUsername(), registerDto.getPasswordHash()));

        var token = jwtService.generateJwtToken(authentication);

        return ResponseEntity.ok()
                .header("Set-Cookie", AUTH_COOKIE + "=" + token + "; Path=/; HttpOnly; SameSite=Strict")
                .build();
    }

    @Override
    public ResponseEntity<Void> login(AuthDto authDto) {
        var authentication = authenticationManager.authenticate(
                new UsernamePasswordAuthenticationToken(authDto.getUsername(), authDto.getPasswordHash()));

        var token = jwtService.generateJwtToken(authentication);

        return ResponseEntity.ok()
                .header("Set-Cookie", AUTH_COOKIE + "=" + token + "; Path=/; HttpOnly; SameSite=Strict")
                .build();
    }

    @Override
    public ResponseEntity<Void> checkSession() {
        return ResponseEntity.ok().build();
    }

    @Override
    public ResponseEntity<UserDto> getSelf() {
        return ResponseEntity.ok(userService.getSelf());
    }

    @Override
    public ResponseEntity<UserDto> getUserById(UUID userId) {
        return userService.getUserById(userId)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @Override
    public ResponseEntity<List<UserDto>> search(String username) {
        return ResponseEntity.ok(userService.searchByUsername(username));
    }
}
