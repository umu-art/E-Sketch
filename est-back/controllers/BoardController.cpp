#include "BoardController.hpp"
using namespace est_back::controller;

void BoardController::listByUserId(const HttpRequestPtr& req, Callback callback, std::string&& userId) {
    auto backBoardListDto = est_back::service::getBackBoardListDto(userId);
    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(k200OK);
    resp->setContentTypeCode(CT_APPLICATION_JSON);
    nlohmann::json j;
    org::openapitools::server::model::to_json(j, backBoardListDto);
    resp->setBody(to_string(j));
    callback(resp);
}

void BoardController::create(const HttpRequestPtr& req, Callback callback, std::string&& userId) {
    try {
        auto body = req->getBody();
        nlohmann::json j = nlohmann::json::parse(body);
        org::openapitools::server::model::UpsertBoardDto upsertBoardDto;
        org::openapitools::server::model::from_json(j, upsertBoardDto);
        auto boardDto = est_back::service::createBoard(upsertBoardDto, userId);
        nlohmann::json respJson = nlohmann::json::object();
        org::openapitools::server::model::to_json(respJson, boardDto);
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_APPLICATION_JSON);
        resp->setBody(respJson.dump());
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

void BoardController::getByUuid(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    try {
        nlohmann::json respJson = nlohmann::json::object();
        auto boardDto = est_back::service::getBoard(boardId);
        org::openapitools::server::model::to_json(respJson, boardDto);
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_APPLICATION_JSON);
        resp->setBody(respJson.dump());
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

void BoardController::update(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    try {
        auto body = req->getBody();
        nlohmann::json j = nlohmann::json::parse(body);
        org::openapitools::server::model::UpsertBoardDto upsertBoardDto;
        org::openapitools::server::model::from_json(j, upsertBoardDto);
        est_back::service::updateBoard(upsertBoardDto, boardId);
        auto boardDto = est_back::service::getBoard(boardId);
        auto respJson = nlohmann::json::object();
        org::openapitools::server::model::to_json(respJson, boardDto);
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_APPLICATION_JSON);
        resp->setBody(respJson.dump());
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

void BoardController::deleteBoard(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    try {
        est_back::service::deleteBoard(boardId);
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

void BoardController::share(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    try {
        auto body = req->getBody();
        auto j = nlohmann::json::parse(body);
        org::openapitools::server::model::BackSharingDto sharingDto;
        org::openapitools::server::model::from_json(j, sharingDto);
        est_back::service::shareBoard(sharingDto, boardId);
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

void BoardController::updateShare(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    try {
        auto body = req->getBody();
        auto j = nlohmann::json::parse(body);
        org::openapitools::server::model::BackSharingDto sharingDto;
        org::openapitools::server::model::from_json(j, sharingDto);
        est_back::service::updateShare(sharingDto, boardId);
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

void BoardController::unshare(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    try {
        auto body = req->getBody();
        auto j = nlohmann::json::parse(body);
        org::openapitools::server::model::UnshareBoardDto unshareBoardDto;
        org::openapitools::server::model::from_json(j, unshareBoardDto);
        est_back::service::unshareBoard(unshareBoardDto, boardId);
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
void BoardController::markAsRecent(const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)>&& callback,
                                   std::string&& userId) {
    try {
        auto body = req->getBody();
        auto j = nlohmann::json::parse(body);
        org::openapitools::server::model::BoardIdDto boardIdDto;
        org::openapitools::server::model::from_json(j, boardIdDto);
        est_back::service::markAsRecent(boardIdDto.getId(), userId);
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
void BoardController::recentsByMinute(const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)>&& callback,
                                      uint32_t minutes) {
    try {
        org::openapitools::server::model::RecentBoardIdListDto recentBoardIdListDto =
            est_back::service::getRecentsByMinute(minutes);
        nlohmann::json j;
        org::openapitools::server::model::to_json(j, recentBoardIdListDto);
        auto resp = HttpResponse::newHttpResponse();
        resp->setStatusCode(k200OK);
        resp->setContentTypeCode(CT_APPLICATION_JSON);
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
