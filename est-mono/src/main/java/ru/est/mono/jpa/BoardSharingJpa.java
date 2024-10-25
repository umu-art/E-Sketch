package ru.est.mono.jpa;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ru.est.mono.domain.BoardSharingEntity;

import java.util.UUID;

@Repository
public interface BoardSharingJpa extends JpaRepository<BoardSharingEntity, UUID> {

}
