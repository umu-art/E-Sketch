package ru.est.mono.configuration.auth;

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

    @ExceptionHandler
    public ResponseEntity<String> handle(Exception ex) {
        log.error("Ошибка при обработке запроса: ", ex);

        int errorCode = 500;
        String errorMessage = "Внутренняя ошибка сервера";

        if (ex instanceof MonoException) {
            errorCode = ((MonoException) ex).getErrorCode();
            errorMessage = ((MonoException) ex).getErrorMessage();
        }

        if (ex instanceof BindException || ex instanceof HttpMessageConversionException || ex instanceof TypeMismatchException) {
            errorCode = 400;
            errorMessage = ex.getLocalizedMessage();
        }

        if (ex instanceof BadCredentialsException) {
            errorCode = 401;
            errorMessage = "Неверный email или пароль";
        }

        return ResponseEntity
                .status(errorCode)
                .body(errorMessage);
    }
}
