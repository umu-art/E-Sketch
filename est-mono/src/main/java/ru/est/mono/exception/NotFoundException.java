package ru.est.mono.exception;

import org.springframework.http.HttpStatus;

import java.util.UUID;

public class NotFoundException extends MonoException {

    public NotFoundException(Class<?> clazz, UUID id) {
        super(String.format("Entity %s with id %s not found", clazz.getSimpleName(), id), HttpStatus.NOT_FOUND);
    }

}
