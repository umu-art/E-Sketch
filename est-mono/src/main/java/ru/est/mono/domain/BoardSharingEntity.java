package ru.est.mono.domain;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.FetchType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import lombok.Data;
import org.hibernate.annotations.Fetch;
import org.hibernate.annotations.FetchMode;
import ru.est.mono.model.SharingDto;

import java.util.UUID;

@Data
@Entity
@Table(name = "board_sharing")
public class BoardSharingEntity {
    @Id
    private UUID id;

    @ManyToOne(optional = false, fetch = FetchType.EAGER)
    @Fetch(FetchMode.JOIN)
    @JoinColumn(name = "user_id", nullable = false)
    private UserEntity user;

    @ManyToOne(optional = false)
    @JoinColumn(name = "board_id", nullable = false)
    private BoardEntity board;

    @Column(name = "sharing_mode", nullable = false)
    @Enumerated(EnumType.STRING)
    private SharingDto.AccessEnum sharingMode;
}
