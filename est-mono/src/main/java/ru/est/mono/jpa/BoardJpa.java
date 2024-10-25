package ru.est.mono.jpa;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;
import ru.est.mono.domain.BoardEntity;

import java.util.List;
import java.util.UUID;

@Repository
public interface BoardJpa extends JpaRepository<BoardEntity, UUID> {

    @Query("""
            select be
            from BoardEntity be
            where be.owner.username = :username
                or be.id in (select s.board.id from BoardSharingEntity s where s.user.username = :username)
            """)
    List<BoardEntity> findAllForUser(@Param("username") String username);
}
