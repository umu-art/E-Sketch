#include "BoardController.hpp"

using namespace est_back::controller;

void BoardController::listByUserId(const HttpRequestPtr& req, Callback callback,
                                   std::string&& userId) {
    LOG_DEBUG << "User: " << userId;
    Json::Value ret;
    ret["userId"] = userId;
    auto resp = HttpResponse::newHttpJsonResponse(ret);
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

void BoardController::getByUuid(const HttpRequestPtr& req, Callback callback,
                                std::string&& boardId) {
}

void BoardController::update(const HttpRequestPtr& req, Callback callback,
                             std::string&& boardId) {
}

void BoardController::deleteBoard(const HttpRequestPtr& req, Callback callback,
                                  std::string&& boardId) {
}

void BoardController::share(const HttpRequestPtr& req, Callback callback,
                            std::string&& boardId) {
}

void BoardController::unshare(const HttpRequestPtr& req, Callback callback,
                              std::string&& boardId) {
}
