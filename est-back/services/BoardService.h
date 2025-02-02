#pragma once

#include <string>
#include <drogon/HttpAppFramework.h>

#include "../../api/build/est-back-cpp/model/BackBoardListDto.h"
#include "../../api/build/est-back-cpp/model/LinkShareMode.h"
#include "../../api/build/est-back-cpp/model/UpsertBoardDto.h"
#include "../../api/build/est-back-cpp/model/UnshareBoardDto.h"
#include "../../api/build/est-back-cpp/model/BoardIdDto.h"
#include "../../api/build/est-back-cpp/model/RecentBoardIdListDto.h"
#include "../utils/utils.h"
#include "../errors/ServiceException.h"
namespace est_back::service {
    namespace osm = org::openapitools::server::model;
    namespace err = est_back::errors;

    osm::BackBoardListDto getBackBoardListDto(const std::string& userId);

    osm::BackBoardDto createBoard(const osm::UpsertBoardDto& upsertBoardDto, const std::string& userId);

    osm::BackBoardDto getBoard(const std::string& boardId);

    void updateBoard(const osm::UpsertBoardDto& upsertBoardDto, const std::string& boardId);

    void deleteBoard(const std::string& boardId);

    void shareBoard(const osm::BackSharingDto& sharingDto, const std::string& boardId);

    void updateShare(const osm::BackSharingDto& sharingDto, const std::string& boardId);

    void unshareBoard(const osm::UnshareBoardDto& unshareBoardDto, const std::string& boardId);

    void markAsRecent(const std::string& boardId, const std::string& userId);

    std::vector<osm::BackBoardDto> getRecentBoards(const std::string& userId);

    osm::RecentBoardIdListDto getRecentsByMinute(uint32_t minutes);

    osm::BackBoardDto rowToBoardDto(const drogon::orm::Row& row,
                                    const std::vector<osm::BackSharingDto>& sharedWith = {});

    void addBackSharingDtoToMap(std::map<std::string, std::vector<osm::BackSharingDto>>& mp,
                                const drogon::orm::Row& sharingDtoRow);

    std::string linkSharedModeToString(const osm::LinkShareMode& linkShareMode);

    void createBoardInDB(const osm::BackBoardDto& boardDto);

    bool boardExists(const std::string& boardId);

}  // namespace est_back::service