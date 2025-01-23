#pragma once

#include <string>
#include <regex>

namespace est_back::utils {
    bool isValidUUID(const std::string& id);
    std::string toUpper(const std::string& s);
    std::string toLower(const std::string& s);
    std::string strVectorToString(const std::vector<std::string>& v);
}  // namespace est_back::utils