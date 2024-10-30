#pragma once

#include <string>
#include <drogon/HttpAppFramework.h>

#include "../../api/build/est-back-cpp/model/BackBoardListDto.h"
#include "../../api/build/est-back-cpp/model/LinkShareMode.h"
#include "../../api/build//est-back-cpp/model/UpsertBoardDto.h"
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

    osm::BackBoardListDto getBackBoardListDto(const std::string& userId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto res = clientPtr->execSqlAsyncFuture(
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
        auto resGet = res.get();
        for (const auto& row : resGet) {
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
            auto mineRes = clientPtr
                               ->execSqlAsyncFuture(
                                   "select board_id, user_id, sharing_mode "
                                   "from board_sharing where board_id in(" +
                                   mineBoardIdStr + ");")
                               .get();
            for (const auto& row : mineRes) {
                addBackSharingDtoToMap(mineSharedWith, row);
            }
        }
        if (!sharedBoardIdStr.empty()) {
            auto sharedRes = clientPtr
                                 ->execSqlAsyncFuture(
                                     "select board_id, user_id, sharing_mode "
                                     "from board_sharing where board_id in(" +
                                     sharedBoardIdStr + ");")
                                 .get();
            for (const auto& row : sharedRes) {
                addBackSharingDtoToMap(sharedSharedWith, row);
            }
        }
        for (const auto& row : resGet) {
            osm::BackBoardDto boardDto;
            boardDto.setId(row["id"].as<std::string>());
            boardDto.setName(row["name"].as<std::string>());
            boardDto.setDescription(row["description"].as<std::string>());
            boardDto.setOwnerId(row["owner_id"].as<std::string>());
            osm::LinkShareMode linkShareMode;
            nlohmann::json j = toLower(row["link_shared_mode"].as<std::string>());
            osm::from_json(j, linkShareMode);
            boardDto.setLinkSharedMode(linkShareMode);
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
        return backBoardListDto;
    }

    std::string linkSharedModeToString(const osm::LinkShareMode& linkShareMode) {
        nlohmann::json j;
        osm::to_json(j, linkShareMode);
        return j.get<std::string>();
    }

    void createBoardInDB(const osm::BackBoardDto& boardDto) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        clientPtr->execSqlAsyncFuture(
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
            org::openapitools::server::model::LinkShareMode::eLinkShareMode::INVALID_VALUE_OPENAPI_GENERATED) {
            linkShareMode.setValue(org::openapitools::server::model::LinkShareMode::eLinkShareMode::NONE_BY_LINK);
        } else {
            linkShareMode.setValue(upsertBoardDto.getLinkSharedMode().getValue());
        }
        boardDto.setLinkSharedMode(linkShareMode);
        createBoardInDB(boardDto);
        return boardDto;
    }

    osm::BackBoardDto getBoard(const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto res = clientPtr->execSqlAsyncFuture(
            "select * from board where id = $1;",
            boardId);
        auto row = res.get().front();
        osm::BackBoardDto boardDto;
        boardDto.setId(row["id"].as<std::string>());
        boardDto.setName(row["name"].as<std::string>());
        boardDto.setDescription(row["description"].as<std::string>());
        boardDto.setOwnerId(row["owner_id"].as<std::string>());
        osm::LinkShareMode linkShareMode;
        nlohmann::json j = toLower(row["link_shared_mode"].as<std::string>());
        osm::from_json(j, linkShareMode);
        boardDto.setLinkSharedMode(linkShareMode);
        return boardDto;
    }
}  // namespace est_back::service