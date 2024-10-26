#include "BoardController.hpp"
#include "../services/BoardService.h"
using namespace est_back::controller;

void BoardController::listByUserId(const HttpRequestPtr& req, Callback callback, std::string&& userId) {
    auto backBoardListDto = est_back::service::getBackBoardListDto(userId);

    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(k200OK);
    resp->setContentTypeCode(CT_APPLICATION_JSON);
    resp->setBody(backBoardListDto.toJsonString());
    callback(resp);
}

void BoardController::create(const HttpRequestPtr& req, Callback callback) {
    auto body = req->getJsonObject();
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
