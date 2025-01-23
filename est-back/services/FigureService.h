#pragma once

#include <drogon/orm/Row.h>
#include <drogon/HttpAppFramework.h>

#include "../../api/build/est-back-cpp/model/FigureListDto.h"
#include "../../api/build/est-back-cpp/model/FigureIdDto.h"
#include "../errors/ServiceException.h"
#include "../utils/utils.h"

namespace est_back::service {
    namespace osm = org::openapitools::server::model;
    namespace err = est_back::errors;

    osm::FigureDto rowToFigureDto(const drogon::orm::Row& row);

    osm::FigureListDto getFigureListDto(const std::string& boardId);

    osm::FigureIdDto createFigure(const std::string& boardId);

    osm::FigureDto getFigure(const std::string& figureId);

    void updateFigure(const osm::FigureDto& figureDto, const std::string& figureId);

    void deleteFigure(const std::string& figureId);
}  // namespace est_back::service
