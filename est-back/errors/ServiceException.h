#pragma once

#include <stdexcept>

namespace est_back::errors {
    enum class ServiceError { NOT_FOUND, BAD_REQUEST };

    class ServiceException : public std::runtime_error {
    public:
        ServiceException(ServiceError errorType, const std::string& message);
        [[nodiscard]] ServiceError getErrorType() const;

    private:
        ServiceError errorType_;
    };
}  // namespace est_back::errors
