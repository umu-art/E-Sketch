package ru.est.mono.configuration.auth;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.tomcat.util.codec.binary.StringUtils;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;
import ru.est.mono.domain.UserEntity;
import ru.est.mono.service.impl.JwtService;

import java.io.IOException;
import java.util.Arrays;
import java.util.List;

import static java.util.Objects.isNull;

@Slf4j
@Component
@RequiredArgsConstructor
public class JwtAuthFilter extends OncePerRequestFilter {

    public static final List<String> IGNORE_URL = List.of("/proxy/user/login", "/proxy/user/register", "/error");
    public static final String AUTH_COOKIE = "eja";
    private final UserDetailsService userDetailsService;
    private final JwtService jwtService;

    public static UserEntity getPerformer() {
        return (UserEntity) SecurityContextHolder.getContext().getAuthentication().getPrincipal();
    }

    @Override
    protected void doFilterInternal(@NonNull HttpServletRequest request, @NonNull HttpServletResponse response, @NonNull FilterChain filterChain) throws ServletException, IOException {
        if (IGNORE_URL.contains(request.getRequestURI())) {
            filterChain.doFilter(request, response);
            return;
        }

        String jwt = parseJwtHeader(request);

        if (jwt == null || !jwtService.validateJwtToken(jwt)) {
            response.setContentType("application/json");
            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
            response.getOutputStream().write(StringUtils.getBytesUtf8("Требуется авторизация"));
            return;
        }

        String username = jwtService.getUserNameFromJwtToken(jwt);

        UserDetails userDetails = userDetailsService.loadUserByUsername(username);
        UsernamePasswordAuthenticationToken authentication = new UsernamePasswordAuthenticationToken(userDetails, null, userDetails.getAuthorities());
        authentication.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));

        SecurityContextHolder.getContext()
                .setAuthentication(authentication);

        filterChain.doFilter(request, response);
    }

    private String parseJwtHeader(HttpServletRequest request) {
        if (isNull(request.getCookies())) {
            return null;
        }

        var authCookie = Arrays.stream(request.getCookies())
                .filter(cookie -> AUTH_COOKIE.equals(cookie.getName()))
                .findFirst();

        return authCookie
                .map(Cookie::getValue)
                .orElse(null);
    }
}