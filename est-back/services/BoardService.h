#pragma once

#include <string>
#include <drogon/HttpAppFramework.h>

#include "../../api/build/est-back-cpp/model/BackBoardListDto.h"
#include "../../api/build/est-back-cpp/model/LinkShareMode.h"
namespace est_back::service {

    std::string toUpper(std::string s) {
        std::string res;
        for (auto c : res) {
            if ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
                res += toupper(c);
            }
        }
        return res;
    }

    org::openapitools::server::model::BackBoardListDto getBackBoardListDto(const std::string& userId) {
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
        std::vector<org::openapitools::server::model::BackBoardDto> mine;
        std::vector<org::openapitools::server::model::BackBoardDto> shared;
        for (const auto& row : res.get()) {
            if (row["id"].as<std::string>() == userId) {
                mineBoardId.push_back(row["id"].as<std::string>());
            } else {
                sharedBoardId.push_back(row["id"].as<std::string>());
            }
        }
        std::map<std::string, std::vector<org::openapitools::server::model::BackSharingDto>> mineSharedWith;
        std::map<std::string, std::vector<org::openapitools::server::model::BackSharingDto>> sharedSharedWith;
//        auto mineRes = clientPtr->execSqlAsyncFuture(
//            "select * "
//            "from board "
//            "where board.owner_id = $1 "
//            "or board.id in $2;", userId, mineBoardId);
        for (const auto& row : res.get()) {
            mineBoardId.push_back(row["id"].as<std::string>());
            if (row["owner_id"].as<std::string>() == userId) {
                org::openapitools::server::model::BackBoardDto boardDto;
                boardDto.setId(row["id"].as<std::string>());
                boardDto.setName(row["name"].as<std::string>());
                boardDto.setDescription(row["description"].as<std::string>());
                boardDto.setOwnerId(row["owner_id"].as<std::string>());
                boardDto.setSharedWith({org::openapitools::server::model::BackSharingDto{}});
                org::openapitools::server::model::LinkShareMode linkShareMode;
                nlohmann::json j = row["link_shared_mode"].as<std::string>();
                org::openapitools::server::model::from_json(j, linkShareMode);
                boardDto.setLinkSharedMode(linkShareMode);
                mine.push_back(boardDto);
            }
        }
        org::openapitools::server::model::BackBoardListDto backBoardListDto;
        backBoardListDto.setMine(mine);
        backBoardListDto.setShared(shared);
        return backBoardListDto;
    }


}  // namespace est_back::service