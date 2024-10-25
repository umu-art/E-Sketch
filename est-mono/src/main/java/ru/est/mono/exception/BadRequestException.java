package ru.est.mono.exception;

import org.springframework.http.HttpStatus;

public class BadRequestException extends MonoException {
    public BadRequestException(String message) {
        super(message, HttpStatus.BAD_REQUEST);
    }
}
