package ru.est.mono.service.impl;

import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.io.Decoders;
import io.jsonwebtoken.security.Keys;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import javax.crypto.SecretKey;
import java.util.Date;

@Slf4j
@Service
public class JwtService {

    @Value("${est.auth.jwtSecret}")
    private String jwtSecret;

    @Value("${est.auth.jwtExpirationHours}")
    private int jwtExpirationHours;

    private SecretKey getSigningKey() {
        return Keys.hmacShaKeyFor(Decoders.BASE64.decode(jwtSecret));
    }

    public String generateJwtToken(Authentication authentication) {
        UserDetails userPrincipal = (UserDetails) authentication.getPrincipal();

        var issuedDate = new Date();

        return Jwts.builder()
                .subject(userPrincipal.getUsername())
                .issuedAt(issuedDate)
                .expiration(new Date(issuedDate.getTime() + jwtExpirationHours * 3600000L))
                .signWith(getSigningKey())
                .compact();

    }

    public boolean validateJwtToken(String authToken) {
        try {
            Jwts.parser()
                    .verifyWith(getSigningKey())
                    .build()
                    .parseSignedClaims(authToken);

            return true;
        } catch (Exception e) {
            log.error("invalid token", e);
        }

        return false;
    }


    public String getUserNameFromJwtToken(String token) {
        return Jwts.parser()
                .verifyWith(getSigningKey())
                .build()
                .parseSignedClaims(token)
                .getPayload()
                .getSubject();
    }
}
