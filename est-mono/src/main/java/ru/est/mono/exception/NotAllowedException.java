package ru.est.mono.exception;

import org.springframework.http.HttpStatus;

public class NotAllowedException extends MonoException {
    public NotAllowedException() {
        super("Доступ запрещен", HttpStatus.FORBIDDEN);
    }
}
