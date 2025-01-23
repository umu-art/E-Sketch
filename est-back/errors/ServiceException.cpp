#include "ServiceException.h"

est_back::errors::ServiceException::ServiceException(ServiceError errorType, const std::string& message)
    : std::runtime_error(message), errorType_(errorType) {
}

est_back::errors::ServiceError est_back::errors::ServiceException::getErrorType() const {
    return errorType_;
}
