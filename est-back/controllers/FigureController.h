#pragma once

#include <drogon/HttpController.h>
#include <nlohmann/json.hpp>
#include "services/FigureService.h"

using namespace drogon;

namespace est_back::controller {
    class FigureController : public drogon::HttpController<FigureController> {
    public:
        METHOD_LIST_BEGIN
        ADD_METHOD_TO(FigureController::listByBoardId, "/back/figure/list/{boardId}", Get);
        ADD_METHOD_TO(FigureController::createFigure, "/back/figure/create/{boardId}", Post);
        ADD_METHOD_TO(FigureController::updateFigure, "/back/figure/{figureId}", Patch);
        ADD_METHOD_TO(FigureController::deleteFigure, "/back/figure/{figureId}", Delete);
        METHOD_LIST_END
    private:
        using Callback = std::function<void(const HttpResponsePtr&)>&&;
        void listByBoardId(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void createFigure(const HttpRequestPtr& req, Callback callback, std::string&& boardId);
        void updateFigure(const HttpRequestPtr& req, Callback callback, std::string&& figureId);
        void deleteFigure(const HttpRequestPtr& req, Callback callback, std::string&& figureId);
    };
}  // namespace est_back::controller
