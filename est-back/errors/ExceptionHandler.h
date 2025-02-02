#pragma once

#include <drogon/drogon.h>
#include <iostream>

#include "errors/ServiceException.h"
#include "HttpResponseFactory.h"

using namespace drogon;

namespace est_back::errors {
    void customExceptionHandler(const std::exception& e, const HttpRequestPtr& req,
                                std::function<void(const HttpResponsePtr&)>&& callback);
}
