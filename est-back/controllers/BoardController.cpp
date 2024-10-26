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

void BoardController::create(const HttpRequestPtr& req, Callback callback) {
    auto body = req->getBody();
    nlohmann::json j = nlohmann::json::parse(body);
    org::openapitools::server::model::UpsertBoardDto upsertBoardDto;
    org::openapitools::server::model::from_json(j, upsertBoardDto);
    //    createBoard(upsertBoardDto);
    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(k200OK);
    resp->setContentTypeCode(CT_TEXT_PLAIN);
    resp->setBody("OK");
    callback(resp);
}

void BoardController::getByUuid(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
}

void BoardController::update(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
}

void BoardController::deleteBoard(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
}

void BoardController::share(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
}

void BoardController::unshare(const HttpRequestPtr& req, Callback callback, std::string&& boardId) {
}
