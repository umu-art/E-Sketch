#include "FigureController.h"
#include "errors/HttpResponseFactory.h"

using namespace est_back::controller;

void FigureController::listByBoardId(const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)>&& callback,
                                     std::string&& boardId) {
    try {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_APPLICATION_JSON);
        nlohmann::json j;
        auto figureListDto = est_back::service::getFigureListDto(boardId);
        org::openapitools::server::model::to_json(j, figureListDto);
        resp->setBody(j.dump());
        callback(resp);
    } catch (const est_back::errors::ServiceException& e) {
        auto resp = est_back::controller::createErrorHttpResponse(e);
        callback(resp);
        return;
    } catch (const std::exception& e) {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k500InternalServerError);
        resp->setContentTypeCode(CT_TEXT_PLAIN);
        resp->setBody("Internal server error");
        callback(resp);
        return;
    }
}

void FigureController::createFigure(const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)>&& callback,
                                    std::string&& boardId) {
    try {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_APPLICATION_JSON);
        nlohmann::json j;
        auto figureIdDto = est_back::service::createFigure(boardId);
        org::openapitools::server::model::to_json(j, figureIdDto);
        resp->setBody(j.dump());
        callback(resp);
    } catch (const est_back::errors::ServiceException& e) {
        auto resp = est_back::controller::createErrorHttpResponse(e);
        callback(resp);
        return;
    } catch (const std::exception& e) {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k500InternalServerError);
        resp->setContentTypeCode(CT_TEXT_PLAIN);
        resp->setBody("Internal server error");
        callback(resp);
        return;
    }
}

void FigureController::getFigure(const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)>&& callback,
                                 std::string&& figureId) {
    try {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_APPLICATION_JSON);
        nlohmann::json j;
        auto figureDto = est_back::service::getFigure(figureId);
        org::openapitools::server::model::to_json(j, figureDto);
        resp->setBody(j.dump());
        callback(resp);
    } catch (const est_back::errors::ServiceException& e) {
        auto resp = est_back::controller::createErrorHttpResponse(e);
        callback(resp);
        return;
    } catch (const std::exception& e) {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k500InternalServerError);
        resp->setContentTypeCode(CT_TEXT_PLAIN);
        resp->setBody("Internal server error");
        callback(resp);
        return;
    }
}

void FigureController::updateFigure(const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)>&& callback,
                                    std::string&& figureId) {
    try {
        auto body = req->getBody();
        nlohmann::json updFigureJson = nlohmann::json::parse(body);
        org::openapitools::server::model::FigureDto updFigureDto;
        org::openapitools::server::model::from_json(updFigureJson, updFigureDto);
        est_back::service::updateFigure(updFigureDto, figureId);
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_TEXT_PLAIN);
        resp->setBody("OK");
        callback(resp);
    } catch (const est_back::errors::ServiceException& e) {
        auto resp = est_back::controller::createErrorHttpResponse(e);
        callback(resp);
        return;
    } catch (const std::exception& e) {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k500InternalServerError);
        resp->setContentTypeCode(CT_TEXT_PLAIN);
        resp->setBody("Internal server error");
        callback(resp);
        return;
    }
}
void FigureController::deleteFigure(const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)>&& callback,
                                    std::string&& figureId) {
    try {
        est_back::service::deleteFigure(figureId);
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_TEXT_PLAIN);
        resp->setBody("OK");
        callback(resp);
    } catch (const est_back::errors::ServiceException& e) {
        auto resp = est_back::controller::createErrorHttpResponse(e);
        callback(resp);
        return;
    } catch (const std::exception& e) {
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k500InternalServerError);
        resp->setContentTypeCode(CT_TEXT_PLAIN);
        resp->setBody("Internal server error");
        callback(resp);
        return;
    }
}
