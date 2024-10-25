package ru.est.mono.domain;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.FetchType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.OneToMany;
import jakarta.persistence.OneToOne;
import jakarta.persistence.Table;
import lombok.Data;
import org.hibernate.annotations.Fetch;
import org.hibernate.annotations.FetchMode;
import ru.est.mono.model.BoardDto;

import java.util.List;
import java.util.UUID;

@Data
@Entity
@Table(name = "board")
public class BoardEntity {
    @Id
    private UUID id;

    @Column(name = "name", nullable = false)
    private String name;

    @Column(name = "description")
    private String description;

    @OneToOne(optional = false, fetch = FetchType.EAGER)
    @Fetch(FetchMode.JOIN)
    @JoinColumn(name = "owner_id", nullable = false)
    private UserEntity owner;

    @OneToMany(mappedBy = "board", fetch = FetchType.EAGER)
    @Fetch(FetchMode.SUBSELECT)
    private List<BoardSharingEntity> boardSharings;

    @Column(name = "link_shared_mode", nullable = false)
    @Enumerated(EnumType.STRING)
    private BoardDto.LinkSharedModeEnum linkSharedMode;
}
