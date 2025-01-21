#pragma once

#include <string>
#include <drogon/HttpAppFramework.h>

#include "../../api/build/est-back-cpp/model/BackBoardListDto.h"
#include "../../api/build/est-back-cpp/model/LinkShareMode.h"
#include "../../api/build/est-back-cpp/model/UpsertBoardDto.h"
#include "../../api/build/est-back-cpp/model/UnshareBoardDto.h"
#include "../../api/build/est-back-cpp/model/BoardIdDto.h"
#include "../../api/build/est-back-cpp/model/RecentBoardIdListDto.h"
namespace est_back::service {
    namespace osm = org::openapitools::server::model;
    std::string toUpper(const std::string& s) {
        std::string res;
        res.reserve(s.size());
        std::transform(s.begin(), s.end(), std::back_inserter(res), [](unsigned char c) { return std::toupper(c); });
        return res;
    }

    std::string toLower(const std::string& s) {
        std::string res;
        res.reserve(s.size());
        std::transform(s.begin(), s.end(), std::back_inserter(res), [](unsigned char c) { return std::tolower(c); });
        return res;
    }

    std::string strVectorToString(const std::vector<std::string>& v) {
        if (v.empty())
            return "";
        std::ostringstream oss;
        for (size_t i = 0; i < v.size(); ++i) {
            oss << "'" << v[i] << "'";
            if (i != v.size() - 1) {
                oss << ",";
            }
        }
        return oss.str();
    }

    osm::BackBoardDto rowToBoardDto(const drogon::orm::Row& row,
                                    const std::vector<osm::BackSharingDto>& sharedWith = {}) {
        osm::BackBoardDto boardDto;
        boardDto.setId(row["id"].as<std::string>());
        boardDto.setName(row["name"].as<std::string>());
        boardDto.setDescription(row["description"].as<std::string>());
        boardDto.setOwnerId(row["owner_id"].as<std::string>());
        boardDto.setSharedWith(sharedWith);
        osm::LinkShareMode linkShareMode;
        nlohmann::json j = toLower(row["link_shared_mode"].as<std::string>());
        osm::from_json(j, linkShareMode);
        boardDto.setLinkSharedMode(linkShareMode);
        return boardDto;
    }

    void addBackSharingDtoToMap(std::map<std::string, std::vector<osm::BackSharingDto>>& mp,
                                const drogon::orm::Row& sharingDtoRow) {
        nlohmann::json backSharingDtoJson = nlohmann::json::object();
        auto boardId = sharingDtoRow["board_id"].as<std::string>();
        backSharingDtoJson["userId"] = sharingDtoRow["user_id"].as<std::string>();
        backSharingDtoJson["access"] = toLower(sharingDtoRow["sharing_mode"].as<std::string>());
        osm::BackSharingDto backSharingDto;
        osm::from_json(backSharingDtoJson, backSharingDto);
        mp[boardId].push_back(backSharingDto);
    }
    std::vector<osm::BackBoardDto> getRecentBoards(const std::string& userId);
    osm::BackBoardListDto getBackBoardListDto(const std::string& userId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto res = clientPtr->execSqlSync(
            "select * "
            "from board "
            "where board.owner_id = $1 "
            "or board.id in (select board_id from board_sharing s where s.user_id = "
            "$1);",
            userId);
        std::vector<std::string> mineBoardId;
        std::vector<std::string> sharedBoardId;
        std::vector<osm::BackBoardDto> mine;
        std::vector<osm::BackBoardDto> shared;
        for (const auto& row : res) {
            if (row["owner_id"].as<std::string>() == userId) {
                mineBoardId.push_back(row["id"].as<std::string>());
            } else {
                sharedBoardId.push_back(row["id"].as<std::string>());
            }
        }
        std::map<std::string, std::vector<osm::BackSharingDto>> mineSharedWith;
        std::map<std::string, std::vector<osm::BackSharingDto>> sharedSharedWith;
        auto mineBoardIdStr = strVectorToString(mineBoardId);
        auto sharedBoardIdStr = strVectorToString(sharedBoardId);
        if (!mineBoardIdStr.empty()) {
            auto mineRes = clientPtr->execSqlSync(
                "select board_id, user_id, sharing_mode "
                "from board_sharing where board_id in(" +
                mineBoardIdStr + ");");
            for (const auto& row : mineRes) {
                addBackSharingDtoToMap(mineSharedWith, row);
            }
        }
        if (!sharedBoardIdStr.empty()) {
            auto sharedRes = clientPtr->execSqlSync(
                "select board_id, user_id, sharing_mode "
                "from board_sharing where board_id in(" +
                sharedBoardIdStr + ");");
            for (const auto& row : sharedRes) {
                addBackSharingDtoToMap(sharedSharedWith, row);
            }
        }
        for (const auto& row : res) {
            osm::BackBoardDto boardDto = rowToBoardDto(row);
            if (row["owner_id"].as<std::string>() == userId) {
                boardDto.setSharedWith(mineSharedWith[row["id"].as<std::string>()]);
                mine.push_back(boardDto);
            } else {
                boardDto.setSharedWith(sharedSharedWith[row["id"].as<std::string>()]);
                shared.push_back(boardDto);
            }
        }
        osm::BackBoardListDto backBoardListDto;
        backBoardListDto.setMine(mine);
        backBoardListDto.setShared(shared);
        backBoardListDto.setRecent(getRecentBoards(userId));
        return backBoardListDto;
    }

    std::string linkSharedModeToString(const osm::LinkShareMode& linkShareMode) {
        nlohmann::json j;
        osm::to_json(j, linkShareMode);
        return j.get<std::string>();
    }

    void createBoardInDB(const osm::BackBoardDto& boardDto) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        clientPtr->execSqlSync(
            "insert into board(id, name, description, owner_id, link_shared_mode) "
            "values($1, $2, $3, $4, $5);",
            boardDto.getId(), boardDto.getName(), boardDto.getDescription(), boardDto.getOwnerId(),
            toUpper(linkSharedModeToString(boardDto.getLinkSharedMode())));
    }

    osm::BackBoardDto createBoard(const osm::UpsertBoardDto& upsertBoardDto, const std::string& userId) {
        osm::BackBoardDto boardDto;
        boardDto.setId(drogon::utils::getUuid());
        boardDto.setName(upsertBoardDto.getName());
        boardDto.setDescription(upsertBoardDto.getDescription());
        boardDto.setOwnerId(userId);
        osm::LinkShareMode linkShareMode;
        if (upsertBoardDto.getLinkSharedMode().getValue() ==
            osm::LinkShareMode::eLinkShareMode::INVALID_VALUE_OPENAPI_GENERATED) {
            linkShareMode.setValue(osm::LinkShareMode::eLinkShareMode::NONE_BY_LINK);
        } else {
            linkShareMode.setValue(upsertBoardDto.getLinkSharedMode().getValue());
        }
        boardDto.setLinkSharedMode(linkShareMode);
        createBoardInDB(boardDto);
        return boardDto;
    }

    osm::BackBoardDto getBoard(const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto res = clientPtr->execSqlSync("select * from board where id = $1;", boardId);
        auto shareWithRes = clientPtr->execSqlSync(
            "select board_id, user_id, sharing_mode "
            "from board_sharing where board_id ='" +
            boardId + "';");
        std::vector<osm::BackSharingDto> shareWith;
        for (const auto& row : shareWithRes) {
            nlohmann::json sharingDtoJson = nlohmann::json::object();
            sharingDtoJson["userId"] = row["user_id"].as<std::string>();
            sharingDtoJson["access"] = toLower(row["sharing_mode"].as<std::string>());
            osm::BackSharingDto sharingDto;
            osm::from_json(sharingDtoJson, sharingDto);
            shareWith.push_back(sharingDto);
        }
        auto row = res.front();
        osm::BackBoardDto boardDto = rowToBoardDto(row, shareWith);

        return boardDto;
    }

    void updateBoard(const osm::UpsertBoardDto& upsertBoardDto, const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        osm::LinkShareMode linkShareMode;
        if (upsertBoardDto.getLinkSharedMode().getValue() ==
            osm::LinkShareMode::eLinkShareMode::INVALID_VALUE_OPENAPI_GENERATED) {
            linkShareMode.setValue(osm::LinkShareMode::eLinkShareMode::NONE_BY_LINK);
        } else {
            linkShareMode.setValue(upsertBoardDto.getLinkSharedMode().getValue());
        }
        clientPtr->execSqlSync("update board set name = $1, description = $2, link_shared_mode = $3 where id = $4;",
                               upsertBoardDto.getName(), upsertBoardDto.getDescription(),
                               toUpper(linkSharedModeToString(linkShareMode)), boardId);
    }

    void deleteBoard(const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");

        clientPtr->execSqlSync(
            "WITH delete_sharing AS ( "
            "DELETE FROM board_sharing WHERE board_id = $1 "
            ") "
            "DELETE FROM board WHERE id = $1;",
            boardId);
    }

    void shareBoard(const osm::BackSharingDto& sharingDto, const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto sharingDtoJson = nlohmann::json::object();
        osm::to_json(sharingDtoJson, sharingDto);
        clientPtr->execSqlSync("insert into board_sharing(id, board_id, user_id, sharing_mode) values($1, $2, $3, $4);",
                               drogon::utils::getUuid(), boardId, sharingDto.getUserId(),
                               toUpper(sharingDtoJson["access"].get<std::string>()));
    }

    void updateShare(const osm::BackSharingDto& sharingDto, const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto sharingDtoJson = nlohmann::json::object();
        osm::to_json(sharingDtoJson, sharingDto);
        clientPtr->execSqlSync("update board_sharing set sharing_mode = $1 where user_id = $2 and board_id = $3;",
                               toUpper(sharingDtoJson["access"].get<std::string>()), sharingDto.getUserId(), boardId);
    }

    void unshareBoard(const osm::UnshareBoardDto& unshareBoardDto, const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        clientPtr->execSqlSync("delete from board_sharing where user_id = $1 and board_id = $2;",
                               unshareBoardDto.getUserId(), boardId);
    }

    void markAsRecent(const std::string& boardId, const std::string& userId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        std::string recentBoardId = drogon::utils::getUuid();
        clientPtr->execSqlAsyncFuture(
            "insert into recent_board("
            "id, board_id, user_id, last_used) values($1, $2, $3, now()) "
            "on conflict (board_id, user_id) do update set board_id = $2, user_id = $3, last_used = now();",
            recentBoardId, boardId, userId);
    }

    std::vector<osm::BackBoardDto> getRecentBoards(const std::string& userId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto res = clientPtr->execSqlSync(
            "select b.* from recent_board r "
            "join board b on r.board_id = b.id "
            "where r.user_id = $1 "
            "order by r.last_used desc;",
            userId);
        std::vector<osm::BackBoardDto> recentBoards;
        for (const auto& row : res) {
            recentBoards.push_back(rowToBoardDto(row));
        }
        std::map<std::string, std::vector<osm::BackSharingDto>> sharedWith;
        std::vector<std::string> boardIds;
        for (const auto& board : recentBoards) {
            boardIds.push_back(board.getId());
        }
        auto boardIdsStr = strVectorToString(boardIds);
        if (!boardIdsStr.empty()) {
            auto sharedWithRes = clientPtr->execSqlSync(
                "select board_id, user_id, sharing_mode "
                "from board_sharing where board_id in(" +
                boardIdsStr + ");");
            for (const auto& row : sharedWithRes) {
                addBackSharingDtoToMap(sharedWith, row);
            }
        }
        for (auto& board : recentBoards) {
            board.setSharedWith(sharedWith[board.getId()]);
        }
        return recentBoards;
    }

    osm::RecentBoardIdListDto getRecentsByMinute(uint32_t minutes) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        std::string minStr = "'" + std::to_string(minutes) + " minutes';";
        auto res = clientPtr
                       ->execSqlAsyncFuture(
                           "select distinct board_id from recent_board "
                           "where last_used >= now() - interval " +
                           minStr)
                       .get();
        std::vector<osm::BoardIdDto> recentsBoardId;
        for (const auto& row : res) {
            osm::BoardIdDto boardIdDto;
            boardIdDto.setId(row["board_id"].as<std::string>());
            recentsBoardId.push_back(boardIdDto);
        }
        osm::RecentBoardIdListDto recentBoardIdListDto;
        recentBoardIdListDto.setBoards(recentsBoardId);
        return recentBoardIdListDto;
    }

}  // namespace est_back::service