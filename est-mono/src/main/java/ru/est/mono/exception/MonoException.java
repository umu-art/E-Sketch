package ru.est.mono.exception;

import lombok.Getter;
import org.springframework.http.HttpStatus;

@Getter
public class MonoException extends RuntimeException {

    private final String errorMessage;
    private final int errorCode;

    public MonoException(String message, HttpStatus code) {
        super(message);
        this.errorMessage = message;
        this.errorCode = code.value();
    }
}
