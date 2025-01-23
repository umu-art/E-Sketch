#include "ExeptionHandler.h"

using namespace drogon;

namespace est_back::errors {
    void customExceptionHandler(const std::exception& e, const HttpRequestPtr& req,
                                std::function<void(const HttpResponsePtr&)>&& callback) {
        std::string pathWithQuery = req->path();
        if (!req->query().empty())
            pathWithQuery += "?" + req->query();

        LOG_TRACE << "Custom handler: Handled exception in " << pathWithQuery << ", what(): " << e.what();

        if (const auto* serviceEx = dynamic_cast<const est_back::errors::ServiceException*>(&e)) {
            auto resp = controller::createErrorHttpResponse(*serviceEx);
            callback(resp);
            return;
        }

        // Unhandled exception
        auto response = HttpResponse::newHttpResponse();
        response->setStatusCode(k500InternalServerError);
        response->setBody(e.what());
        callback(response);
    }
}  // namespace est_back::errors
