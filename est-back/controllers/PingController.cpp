#include "PingController.hpp"

using namespace est_back::controller;

void PingController::get(const HttpRequestPtr& req, Callback callback) const {
    auto clientPtr = drogon::app().getDbClient("est-data");

    std::cout << clientPtr->hasAvailableConnections() << std::endl;

    auto f = clientPtr->execSqlAsyncFuture("select 1");
    try {
        auto result = f.get();
        std::cout << result.size() << " rows selected!" << std::endl;
    } catch (const drogon::orm::DrogonDbException& e) {
        std::cerr << "error:" << e.base().what() << std::endl;
    }

    auto resp = HttpResponse::newHttpResponse();
    resp->setStatusCode(k200OK);
    resp->setContentTypeCode(CT_TEXT_PLAIN);
    resp->setBody("OK");
    callback(resp);
}
