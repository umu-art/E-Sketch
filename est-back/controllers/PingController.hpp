#pragma once

#include <drogon/HttpController.h>

using namespace drogon;

namespace est_back::controller {
    class PingController : public drogon::HttpController<PingController> {
    public:
        METHOD_LIST_BEGIN
        ADD_METHOD_TO(PingController::get, "/ping", Get);
        METHOD_LIST_END
    private:
        using Callback = std::function<void(const HttpResponsePtr&)>&&;
        void get(const HttpRequestPtr& req, Callback callback) const;
    };
}  // namespace est_back::controller
