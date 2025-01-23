#include "ExceptionHandlingMiddleware.h"

using namespace drogon;

namespace est_back::middleware {

    void ExceptionHandlingMiddleware::invoke(const HttpRequestPtr& req,
                                             MiddlewareNextCallback &&nextCb,
                                             MiddlewareCallback &&mcb) {
        try {
            LOG_DEBUG << "Entering ExceptionHandlingMiddleware";
            // 在调用下一个中间件之前可以做一些操作
            nextCb([mcb = std::move(mcb)](const HttpResponsePtr &resp) {
                // 在下一个中间件返回后处理响应
                mcb(resp);
            });
        } catch (const errors::ServiceException& e) {
            LOG_ERROR << "Caught ServiceException: " << e.what();
            auto resp = controller::createErrorHttpResponse(e);
            mcb(resp); // 返回错误响应
        } catch (const std::exception& e) {
            LOG_ERROR << "Caught exception: " << e.what();
            auto resp = HttpResponse::newHttpResponse();
            resp->setStatusCode(HttpStatusCode::k500InternalServerError);
            resp->setContentTypeCode(ContentType::CT_TEXT_PLAIN);
            resp->setBody("Internal server error");
            mcb(resp); // 返回500错误响应
        }
    }
}
