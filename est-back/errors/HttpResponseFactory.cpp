#include "HttpResponseFactory.h"

std::shared_ptr<drogon::HttpResponse> est_back::controller::createErrorHttpResponse(
    const est_back::errors::ServiceException& exception) {
    auto resp = drogon::HttpResponse::newHttpResponse();
    resp->setContentTypeCode(drogon::ContentType::CT_TEXT_PLAIN);
    switch (exception.getErrorType()) {
        case est_back::errors::ServiceError::BAD_REQUEST:
            resp->setStatusCode(drogon::HttpStatusCode::k400BadRequest);
            break;
        case est_back::errors::ServiceError::NOT_FOUND:
            resp->setStatusCode(drogon::HttpStatusCode::k404NotFound);
            break;
    }
    resp->setBody(exception.what());
    return resp;
}
