/**
*
*  est_back_middleware_ExceptionHandlingMiddleware.h
*
*/

#pragma once

#include <drogon/HttpMiddleware.h>
#include "./errors/ServiceException.h" // 确保包含ServiceException的头文件
#include "./errors/HttpResponseFactory.h" // 确保包含HttpResponseFactory的头文件

using namespace drogon;

namespace est_back::middleware {

   class ExceptionHandlingMiddleware : public drogon::HttpMiddleware<ExceptionHandlingMiddleware> {
   public:
       ExceptionHandlingMiddleware() = default;

       void invoke(const HttpRequestPtr& req,
                   MiddlewareNextCallback &&nextCb,
                   MiddlewareCallback &&mcb) override;
   };

}  // namespace est_back::middleware
