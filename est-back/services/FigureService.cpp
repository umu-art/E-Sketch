#include "FigureService.h"
#include "BoardService.h"

namespace est_back::service {
    namespace osm = org::openapitools::server::model;
    namespace err = est_back::errors;

    osm::FigureDto rowToFigureDto(const drogon::orm::Row& row) {
        osm::FigureDto figureDto;
        figureDto.setData(row["figure_data"].as<std::string>());
        return figureDto;
    }

    osm::FigureListDto getFigureListDto(const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        if (!est_back::service::boardExists(boardId)) {
            throw err::ServiceException(err::ServiceError::NOT_FOUND, "Board not found");
        }
        auto res = clientPtr->execSqlSync("select figure_data from figure where board_id = $1;", boardId);
        osm::FigureListDto figureListDto;
        std::vector<osm::FigureDto> figureList;
        for (const auto& row : res) {
            figureList.push_back(rowToFigureDto(row));
        }
        figureListDto.setFigures(figureList);
        return figureListDto;
    }

    osm::FigureIdDto createFigure(const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        if (!est_back::service::boardExists(boardId)) {
            throw err::ServiceException(err::ServiceError::NOT_FOUND, "Board not found");
        }
        auto figureId = drogon::utils::getUuid();
        clientPtr->execSqlSync("insert into figure(id, board_id) values($1, $2);", figureId, boardId);
        osm::FigureIdDto figureIdDto;
        figureIdDto.setId(figureId);
        return figureIdDto;
    }

    osm::FigureDto getFigure(const std::string& figureId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        if (!figureExists(figureId)) {
            throw err::ServiceException(err::ServiceError::NOT_FOUND, "Figure not found");
        }
        auto res = clientPtr->execSqlSync("select figure_data from figure where id = $1;", figureId);
        return rowToFigureDto(res[0]);
    }

    void updateFigure(const osm::FigureDto& figureDto, const std::string& figureId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        if (!figureExists(figureId)) {
            throw err::ServiceException(err::ServiceError::NOT_FOUND, "Figure not found");
        }
        auto figureData = figureDto.getData();
        std::vector<char> figureDataVec(figureData.begin(), figureData.end());
        clientPtr->execSqlSync("update figure set figure_data = $1 where id = $2;", figureDataVec, figureId);
    }

    void deleteFigure(const std::string& figureId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        if (!figureExists(figureId)) {
            throw err::ServiceException(err::ServiceError::NOT_FOUND, "Figure not found");
        }
        clientPtr->execSqlSync("delete from figure where id = $1", figureId);
    }

    bool figureExists(const std::string& figureId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto existsRes = clientPtr->execSqlSync("select exists(select 1 from figure where id = $1);", figureId);
        return existsRes[0][0].as<bool>();
    }

}  // namespace est_back::service