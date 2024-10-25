package ru.est.mono.configuration.auth;

import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.TypeMismatchException;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageConversionException;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.validation.BindException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import ru.est.mono.exception.MonoException;

@Slf4j
@ControllerAdvice
public class ExceptionControllerAdvice {

    @Data
    public static class ErrorDto {
        private int errorCode;
        private String errorMessage;
    }

    @ExceptionHandler
    public ResponseEntity<ErrorDto> handle(Exception ex) {
        log.error("Ошибка при обработке запроса: ", ex);

        var error = new ErrorDto();
        error.setErrorCode(500);
        error.setErrorMessage("Внутренняя ошибка сервера");

        if (ex instanceof MonoException) {
            error.setErrorCode(((MonoException) ex).getErrorCode());
            error.setErrorMessage(((MonoException) ex).getErrorMessage());
        }

        if (ex instanceof BindException || ex instanceof HttpMessageConversionException || ex instanceof TypeMismatchException) {
            error.setErrorCode(400);
            error.setErrorMessage(ex.getLocalizedMessage());
        }

        if (ex instanceof BadCredentialsException) {
            error.setErrorCode(401);
            error.setErrorMessage("Неверный логин или пароль");
        }

        return ResponseEntity
                .status(error.getErrorCode())
                .body(error);
    }
}
