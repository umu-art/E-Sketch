#pragma once

#include "../../api/build/est-back-cpp/model/FigureListDto.h"
#include "../../api/build/est-back-cpp/model/FigureIdDto.h"

namespace est_back::service {
    namespace osm = org::openapitools::server::model;

    osm::FigureDto rowToFigureDto(const drogon::orm::Row& row) {
        osm::FigureDto figureDto;
        figureDto.setData(row["figure_data"].as<std::string>());
        return figureDto;
    }

    osm::FigureListDto getFigureListDto(const std::string& boardId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
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
        auto figureId = drogon::utils::getUuid();
        clientPtr->execSqlSync("insert into figure(id, board_id) values($1, $2);", figureId, boardId);
        osm::FigureIdDto figureIdDto;
        figureIdDto.setId(figureId);
        return figureIdDto;
    }

    void updateFigure(const osm::FigureDto& figureDto, const std::string& figureId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        auto figureData = figureDto.getData();
        clientPtr->execSqlSync("update figure set figure_data = $1 where id = $2;", figureData, figureId);
    }

    void deleteFigure(const std::string& figureId) {
        auto clientPtr = drogon::app().getDbClient("est-data");
        clientPtr->execSqlSync("delete from figure where id = $1", figureId);
    }
}  // namespace est_back::service
