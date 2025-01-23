#pragma once

#include <drogon/HttpController.h>
#include "errors/ServiceException.h"

namespace est_back::controller {
    std::shared_ptr<drogon::HttpResponse> createErrorHttpResponse(const est_back::errors::ServiceException& exception);
}  // namespace est_back::controller