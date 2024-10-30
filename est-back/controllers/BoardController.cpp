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
}

void BoardController::getByUuid(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    nlohmann::json respJson = nlohmann::json::object();
    auto boardDto = est_back::service::getBoard(boardId);
    org::openapitools::server::model::to_json(respJson, boardDto);
    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(k200OK);
    resp->setContentTypeCode(CT_APPLICATION_JSON);
    resp->setBody(respJson.dump());
    callback(resp);
}

void BoardController::update(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
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
}

void BoardController::deleteBoard(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
    est_back::service::deleteBoard(boardId);
    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(k200OK);
    resp->setContentTypeCode(CT_TEXT_PLAIN);
    resp->setBody("OK");
    callback(resp);
}

void BoardController::share(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
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
}

void BoardController::updateShare(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
}

void BoardController::unshare(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
}
